package intuit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type PlayerAPI struct {
	client *SQLClient
	mux    *mux.Router
}

func NewServer(config SQLConfig) (*PlayerAPI, error) {
	client, err := NewSQLClinet(&config)
	if err != nil {
		return nil, err
	}

	return &PlayerAPI{
		client: client,
		mux:    mux.NewRouter(),
	}, nil

}

func (p *PlayerAPI) registerRoutes() *mux.Router {
	p.mux.HandleFunc("/api/players", p.GetAllPlayers)
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
