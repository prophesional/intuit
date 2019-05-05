package intuit

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3" // init mysql driver support
	"go.uber.org/zap"
)

type logger struct {
	*zap.Logger
}

func (l *logger) logError(err error, msg string) {
	l.Error(msg, zap.Error(err))
}
func (l *logger) logInfo(msg string) {
	l.Info(msg)
}

// SQLClient connects to the DB to perform operations
type SQLClient struct {
	*sql.DB
	*logger
	sqlType string
}

func NewLocalSqlClient(dbpath string, schema []string) (*SQLClient, error) {

	/*	path := filepath.Base(dbpath)
		fmt.Println("Path is ", path)
		_, noPathErr := os.Stat(path)

		fmt.Println("Path Error is ", noPathErr)

		if noPathErr != nil {
			mkdirErr := os.MkdirAll(path, 0755)
			if mkdirErr != nil {
				fmt.Println("Mkdir error ", mkdirErr)
				return nil, mkdirErr
			}
		}*/
	if _, err := os.Create(dbpath); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return nil, err
	}

	l, err := zap.NewProduction()
	if err != nil {
		return nil, err

	}
	client := &SQLClient{
		DB:      db,
		logger:  &logger{Logger: l},
		sqlType: "sqlite3",
	}

	for _, statement := range schema {
		if statement == "" {
			continue
		}
		_, err = db.Exec(statement)
		if err != nil {
			return nil, err
		}
	}

	return client, nil

}
func NewLocalSqlClientWithFile(dbpath string, schemaFile string) (*SQLClient, error) {
	raw, err := ioutil.ReadFile(schemaFile)
	if err != nil {
		return nil, err
	}
	stmt := string(raw)

	var schema []string
	// separate the individual CREATE statements
	for idx := strings.Index(stmt, "CREATE"); idx >= 0; idx = strings.Index(stmt, "CREATE") {
		next := strings.Index(stmt[idx+1:], "CREATE")
		if next < 0 {
			next = len(stmt) - 1
		} else {
			next += idx + 1
		}
		schema = append(schema, strings.TrimSpace(stmt[idx:next]))
		stmt = stmt[next:]
	}

	return NewLocalSqlClient(dbpath, schema)
}

// NewSQLClinet constructs a new SQL db client and connects to the database
func NewSQLClinet(config *SQLConfig) (*SQLClient, error) {
	logger := logger{}
	l, err := zap.NewProduction()
	if err != nil {
		return nil, err

	}
	logger.Logger = l
	c := SQLClient{}
	c.logger = &logger

	connStr, err := config.ConnectionString()
	if err != nil {
		return nil, err
	}
	db, err := sql.Open(config.Type, connStr)
	if err != nil {
		l.Error("cannot open sql connection", zap.Error(err))
		return nil, err
	}
	c.DB = db
	c.sqlType = config.Type
	return &c, c.Ping()
}

// GetPlayerById returns a player object from the db based on the playerId
func (c *SQLClient) GetPlayerById(id string) (*Player, error) {
	row, err := c.Query("SELECT * from Players where playerId  = ?", id)
	if row == nil {
		return nil, nil
	}
	if err != nil {
		c.logError(err, fmt.Sprintf("Error executing select statement for user id : %v", id))
		return nil, err
	}
	var player *Player
	for row.Next() {
		player, err = c.convertRow(row)
	}
	return player, err
}

// GetAllPlayers returns all players currently in the database
func (c *SQLClient) GetAllPlayers() ([]*Player, error) {
	var players []*Player
	rows, err := c.Query("SELECT * from Players")

	if err != nil {
		c.logError(err, "Error executing select statement for user id")
		return nil, err
	}

	for rows.Next() {
		if player, err := c.convertRow(rows); err == nil {
			players = append(players, player)
		} else {
			c.logError(err, "error ocurred scanning row ")
			continue
		}
	}
	return players, nil
}

// InsertPlayers inserts a batch of players into the Db
func (c *SQLClient) InsertPlayers(players []*Player) error {
	args := []interface{}{}
	sqlString := `INSERT into Players (PlayerID,BirthDate,BirthCountry,BirthState,BirthCity,DeathDate,DeathCountry,DeathState,DeathCity,NameFirst, 
				 NameLast,NameGiven,Weight,Height,Bats,Throws,Debut,FinalGame,RetroID,BbrefID ) VALUES `

	for _, player := range players {
		deathDate := getSQLDate(player.DeathDate)
		debutDate := getSQLDate(player.Debut)
		finalGame := getSQLDate(player.FinalGame)
		values := fmt.Sprintf("(%v)", strings.Repeat("?,", 20))
		sqlString += values[0:len(values)-2] + "),"
		if player.DeathDate.IsZero() {
			fmt.Println(player.DeathDate)
			deathDate = nil
		} else {
			deathDate = player.DeathDate
		}
		args = append(args, player.PlayerID, player.BirthDate, player.BirthCountry, player.BirthState, player.BirthCity, deathDate, player.DeathCountry,
			player.DeathState, player.DeathCity, player.NameFirst, player.NameLast, player.NameGiven, player.Weight, player.Height, player.Bats,
			player.Throws, debutDate, finalGame, player.RetroID, player.BbrefID)
	}
	updatedSQL := sqlString[0 : len(sqlString)-1]
	var upsertSQL string
	if c.sqlType != "sqlite3" {
		upsertSQL = updatedSQL + " ON DUPLICATE KEY UPDATE playerId=playerId"
	} else {
		upsertSQL = updatedSQL
	}
	stmt, err := c.Prepare(upsertSQL)

	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(args...)
	return err

}

func getSQLDate(t time.Time) interface{} {

	if t.IsZero() {
		return nil
	}
	return t
}

func (c *SQLClient) convertRow(rows *sql.Rows) (player *Player, err error) {
	player = &Player{}
	var nullableDeathDate, nullableFinalGame mysql.NullTime
	var nullableDeathCity, nullableDeathState, nullableDeathCountry sql.NullString

	err = rows.Scan(&player.PlayerID,
		&player.BirthDate, &player.BirthCountry, &player.BirthState, &player.BirthCity,
		&nullableDeathDate, &nullableDeathCity, &nullableDeathState, &nullableDeathCountry,
		&player.NameFirst, &player.NameLast, &player.NameGiven,
		&player.Weight, &player.Height, &player.Bats, &player.Throws,
		&player.Debut, &nullableFinalGame,
		&player.RetroID, &player.BbrefID,
	)
	if err != nil {
		return
	}

	if nullableDeathDate.Valid {
		player.DeathDate = nullableDeathDate.Time
	}

	if nullableFinalGame.Valid {
		player.FinalGame = nullableFinalGame.Time
	}

	if nullableDeathCity.Valid {
		player.DeathCity = nullableDeathCity.String
	}

	if nullableDeathState.Valid {
		player.DeathState = nullableDeathState.String
	}

	if nullableDeathCountry.Valid {
		player.DeathCountry = nullableDeathCountry.String
	}

	return
}

// Close closes the sql client
func (c *SQLClient) Close() {
	c.Close()
}

// ingestion

func ConvertToPlayer(path string) ([]*Player, error) {

	var players []*Player
	f, err := os.Open(path) //parse if it starts with S3
	if err != nil {
		return nil, err
	}
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}
	for i := 1; i < len(lines); i++ {

		line := lines[i]

		weight, err := strconv.Atoi(line[16])
		if err != nil {
			fmt.Println(err)
			continue
		}
		height, err := strconv.Atoi(line[17])
		if err != nil {
			fmt.Println(err)
			continue
		}

		deathDate, err := getDate(line[7], line[8], line[9])
		if err != nil {
			fmt.Println(err)
			continue
		}
		birthDate, err := getDate(line[1], line[2], line[3])
		if err != nil {
			fmt.Println(err)
			continue
		}

		debutDate, err := parseDate(line[20])
		if err != nil {
			fmt.Println(err)
			continue
		}
		finalGameDate, err := parseDate(line[21])
		if err != nil {
			fmt.Println(err)
			continue
		}

		player := &Player{
			PlayerID:     line[0],
			BirthCountry: line[4],
			BirthState:   line[5],
			BirthCity:    line[6],
			BirthDate:    birthDate,
			DeathDate:    deathDate,
			DeathCountry: line[10],
			DeathState:   line[11],
			DeathCity:    line[12],
			NameFirst:    line[13],
			NameLast:     line[14],
			NameGiven:    line[15],
			Weight:       weight,
			Height:       height,
			Bats:         line[18],
			Throws:       line[19],
			Debut:        debutDate,
			FinalGame:    finalGameDate,
			RetroID:      line[22],
			BbrefID:      line[23],
		}

		players = append(players, player)

	}

	return players, nil

}

func getDate(year, month, day string) (time.Time, error) {
	if year == "" || month == "" || day == "" {
		return time.Time{}, nil
	}

	dYear, err := strconv.Atoi(year)
	if err != nil {
		return time.Time{}, err
	}

	dMonth, err := strconv.Atoi(month)
	if err != nil {
		return time.Time{}, err
	}
	dDay, err := strconv.Atoi(day)
	if err != nil {
		return time.Time{}, err
	}

	if dYear == 0 || dMonth == 0 || dDay == 0 {
		return time.Time{}, nil
	}
	return time.Date(dYear, time.Month(dMonth), dDay, 0, 0, 0, 0, time.UTC), nil
}

func parseDate(d string) (time.Time, error) {

	date, err := time.Parse(layout, d)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}
