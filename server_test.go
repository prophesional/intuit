package intuit

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
)

func (p *PlayerAPITestSuite) TestPlayerAPI_GetAllPlayers() {
	expectedPlayer, err := p.client.GetPlayerById("playerWithDeathDate")
	p.Require().NotNil(expectedPlayer)
	response, err := expectedPlayer.MarshalJSON()

	//fmt.Println(string(response))
	p.Require().NoError(err)

	p.Require().NoError(err)
	server := &PlayerAPI{
		client: p.client,
		mux:    mux.NewRouter(),
	}
	router := server.registerRoutes()

	req, err := http.NewRequest("GET", "/api/players/playerWithDeathDate", nil)
	values := mux.Vars(req)

	fmt.Println(values)
	p.Require().NoError(err)
	rr := httptest.NewRecorder()
	http.Handle("/", router)
	router.ServeHTTP(rr, req)
	p.Assert().Equal(rr.Code, 200)
	p.Assert().Equal(rr.Header().Get("Content-Type"), "application/json")
	p.Assert().JSONEq(rr.Body.String(), string(response))

}
