package intuit

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
)

var log *zap.Logger

func init() {
	log, _ = zap.NewProduction()

}

type AWSConfig struct {
	Region          string `env:"AWS_REGION" cli:"aws-region"`
	AccessKeyID     string `env:"AWS_ACCESS_KEY" cli:"aws-access-key"`
	SecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY" cli:"aws-secret-access-key"`
}

func (awsConfig *AWSConfig) GetAWSSession() (*session.Session, error) {
	if awsConfig.Region == "" {
		reg := os.Getenv("AWS_REGION")
		if reg == "" {
			return nil, fmt.Errorf("AWS Region must be supplied to authenticate")
		}
		awsConfig.Region = reg
	}

	if awsConfig.SecretAccessKey == "" || awsConfig.AccessKeyID == "" {
		return session.NewSession(&aws.Config{Region: aws.String(awsConfig.Region)})
	}
	return session.NewSession(awsConfig.convertConfigs())
}

func (awsConfig *AWSConfig) convertConfigs() *aws.Config {
	creds := credentials.Value{
		AccessKeyID:     awsConfig.AccessKeyID,
		SecretAccessKey: awsConfig.SecretAccessKey,
	}
	c := aws.Config{
		Region:      aws.String(awsConfig.Region),
		Credentials: credentials.NewStaticCredentialsFromCreds(creds),
	}
	return &c
}

type SQLConfig struct {
	ServerName        string `env:"DB_SERVER_NAME" cli:"db-server-name"`
	DatabaseName      string `env:"DB_DATABASE_NAME" cli:"db-database-name"`
	UserName          string `env:"DB_DATABASE_USERNAME" cli:"db-database-username"`
	Password          string `env:"DB_DATABASE_PASSWORD" cli:"db-database-password"`
	Type              string `env:"DB_DATABASE_TYPE" cli:"db-database-type"`
	AWSConfig         AWSConfig
	DatabaseSecretKey string `env:"DB_DATABASE_SECRET_KEY" cli:"db-database-secret-key"`
}

type dbOptions struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	DatabaseName string `json:"databaseName"`
	Host         string `json:"host"`
}

func (c *SQLConfig) ConnectionString() (string, error) {
	//Use Secrets Manager for all settings
	if c.DatabaseSecretKey != "" {
		if dbOpts, err := c.getDBSecrets(); err == nil {

			return c.getDBConnectionString(dbOpts), nil
		} else {
			return "", err
		}
	}
	//use IAM with RDS
	if c.UserName == "" || c.ServerName == "" || c.DatabaseName == "" {
		return "", fmt.Errorf("Username, ServerName and Database name are required for RDS access with IAM.  username: %v servername: %v databaseName: %v", c.UserName, c.ServerName, c.DatabaseName)
	}
	var passwd string
	var err error
	if c.Password == "" {
		passwd, err = c.generateAuthToken()
		if err != nil {
			return "", err
		}
	} else {
		passwd = c.Password
	}

	dbOpts := &dbOptions{
		Password:     passwd,
		Username:     c.UserName,
		Host:         c.ServerName,
		DatabaseName: c.DatabaseName,
	}
	return c.getDBConnectionString(dbOpts), nil

}

func (c *SQLConfig) generateAuthToken() (string, error) {
	session, err := c.AWSConfig.GetAWSSession()
	if err != nil {
		return "", err
	}
	if session == nil {
		return "", errors.New("nil AWS session detected")
	}
	return rdsutils.BuildAuthToken(c.ServerName, c.AWSConfig.Region, c.UserName, session.Config.Credentials)
}

func (c *SQLConfig) getDBSecrets() (*dbOptions, error) {

	s, err := c.AWSConfig.GetAWSSession()
	if err != nil {
		log.Error("error occurred getting session: %s ", zap.Error(err))
		return nil, err
	}

	secretsClient := secretsmanager.New(s)
	secretsInput := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(c.DatabaseSecretKey),
	}

	result, err := secretsClient.GetSecretValue(secretsInput)

	if err != nil {
		log.Error(fmt.Sprintf("error occurred getting secret value %s", c.DatabaseSecretKey), zap.Error(err))
		return nil, err
	}

	var dbOpts dbOptions
	err = json.Unmarshal([]byte(*result.SecretString), &dbOpts)
	if err != nil {
		log.Error("error occurred while unmarshalling secret string ", zap.Error(err))
		return nil, err
	}
	return &dbOpts, nil
}

func (c *SQLConfig) getDBConnectionString(options *dbOptions) string {
	return fmt.Sprintf("%v:%v@tcp(%v)/%v?parseTime=true&allowCleartextPasswords=true&tls=false",
		options.Username,
		options.Password,
		options.Host,
		options.DatabaseName)
}
