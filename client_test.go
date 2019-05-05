package intuit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type PlayerAPITestSuite struct {
	suite.Suite
	client   *SQLClient
	dataPath string
	players  []*Player
}

func TestRun(t *testing.T) {
	suite.Run(t, &PlayerAPITestSuite{})
}
func (p *PlayerAPITestSuite) SetupTest() {
	p.dataPath = "./data/test.csv"

	client, err := NewLocalSqlClientWithFile("./test.sql", "./data/sqllite-schema.sql")
	p.Require().NoError(err)
	p.client = client
	players, err := ConvertToPlayer(p.dataPath)
	p.players = players
	p.Require().NoError(err)

	err = p.client.InsertPlayers(players)
	p.Require().NoError(err)

}

func (p *PlayerAPITestSuite) TearDownTest() {

	//os.Remove("./test.sql")
	//p.players = make([]*Player, 0)
}

func (p *PlayerAPITestSuite) TestConvertPlayers() {
	p.Assert().Len(p.players, 3)
	p.Assert().Equal("playerWithDeathDate", p.players[0].PlayerID)
	p.Assert().Equal(time.Date(2040, 01, 02, 0, 0, 0, 0, time.UTC), p.players[0].DeathDate)

}

func (p *PlayerAPITestSuite) TestGetAllPlayers() {

}

func (p *PlayerAPITestSuite) TestGetPlayerById() {

	player, err := p.client.GetPlayerById("playerWithDeathDate")
	p.Require().NoError(err)
	p.Assert().NotNil(player)
	p.Assert().Equal(player.DeathDate, p.players[0].DeathDate)

}

func (p *PlayerAPITestSuite) TestIngestPlayer() {
	playersFromDb, err := p.client.GetAllPlayers()
	p.Require().NoError(err)

	p.Assert().NotNil(playersFromDb)
	p.Assert().Len(playersFromDb, 3)

}
