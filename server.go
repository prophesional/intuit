package intuit

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type PlayerAPI struct {
	client *SQLClient
	mux    *mux.Router
}

func NewServer(config SQLConfig) (*PlayerAPI, error) {
	client, err := NewSQLClient(&config)
	if err != nil {
		return nil, err
	}

	return &PlayerAPI{
		client: client,
		mux:    mux.NewRouter(),
	}, nil

}

func (p *PlayerAPI) registerRoutes() *mux.Router {
	p.mux.HandleFunc("/api/players", p.Players)
	p.mux.HandleFunc("/api/players/{id}", p.GetPlayById)
	return p.mux

}

func (p *PlayerAPI) Start() {
	p.registerRoutes()
	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      p.mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

}

func (p *PlayerAPI) GetPlayById(res http.ResponseWriter, req *http.Request) {
	values := mux.Vars(req)
	var player *Player
	var err error
	i, ok := values["id"]
	if ok {
		player, err = p.client.GetPlayerById(i)
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
		if player == nil {
			http.Error(res, fmt.Sprintf("player %v not found", i), 404)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(res)

		enc.Encode(player)
		return

	} else {
		http.Error(res, "unable to parse request", 401)
		return
	}

}

func (p *PlayerAPI) GetAllPlayers(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	players, err := p.client.GetAllPlayers()
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	enc := json.NewEncoder(res)

	enc.Encode(players)
}

func (p *PlayerAPI) Players(res http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		p.GetAllPlayers(res, req)
	}
	if req.Method == "POST" {
		p.upload(res, req)
	}
}
func (p *PlayerAPI) InsertPlayers(res http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()
	var players []*Player

	err := decoder.Decode(&players)
	if err != nil {
		fmt.Println("json decoding err", err)
		http.Error(res, err.Error(), 400)
		return
	}
	err = p.client.InsertPlayers(players)
	if err != nil {
		fmt.Println("ingestion error", err)
		http.Error(res, err.Error(), 500)

	}
}

func (p *PlayerAPI) upload(res http.ResponseWriter, req *http.Request) {
	//this function returns the filename(to save in database) of the saved file or an error if it occurs

	req.ParseMultipartForm(32 << 20)

	//ParseMultipartForm parses a request body as multipart/form-data

	file, handler, err := req.FormFile("People.csv") //retrieve the file from form data
	if file == nil {
		http.Error(res, "file People.csv not found", 404)
		return
	}
	//defer file.Close() //close the file when we finish

	if err != nil {
		fmt.Println("error reading file", err)
		http.Error(res, err.Error(), 400)
		return
	}

	//this is path which  we want to store the file

	f, err := os.OpenFile("/tmp/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("error saving file", err)
		http.Error(res, err.Error(), 400)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	players, err := ConvertToPlayer("/tmp/" + handler.Filename)
	if err != nil {
		fmt.Println("error converting csv file", err)
		http.Error(res, err.Error(), 400)
		return
	}

	err = p.client.InsertPlayers(players)
	if err != nil {
		fmt.Println("ingestion error", err)
		http.Error(res, err.Error(), 500)

	}
	res.WriteHeader(201)
}
