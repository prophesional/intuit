package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gravitational/configure"
	"github.com/prophesional/intuit"
)

func main() {
	h := Handler{}
	lambda.Start(h.Handle)
}

type Handler struct {
	sqlClient intuit.SQLClient
}

func (h *Handler) Handle(ctx context.Context, s3Event *events.S3Event) {
	path := "/tmp/tempfile"
	var eplayers []*intuit.Player
	var config intuit.SQLConfig

	err := configure.ParseEnv(&config)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println("Debug:  Config is: ", config)
	for _, r := range s3Event.Records {

		sess := session.Must(session.NewSession())
		// Create a downloader with the session and default options
		downloader := s3manager.NewDownloader(sess)
		// Create a file to write the S3 Object contents to.
		f, err := os.Create(path)

		if err != nil {
			fmt.Println(err)
		}
		// Write the contents of S3 Object to the file
		_, err = downloader.Download(f, &s3.GetObjectInput{
			Bucket: aws.String(r.S3.Bucket.Name),
			Key:    aws.String(r.S3.Object.Key),
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// print the content as 'bytes'

		players, err := intuit.ConvertToPlayer(path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		eplayers = players

		/*if err = h.sqlClient.InsertPlayers(players); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		*/
	}

	url := "intuit-demo.prophesionalizm.net"
	b, err := json.Marshal(eplayers)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

}
