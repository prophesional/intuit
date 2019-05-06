package intuit

import (
	"encoding/json"
	"time"
)

const (
	layout = "2006-01-02"
)

type Player struct {
	PlayerID     string `json:"playerId"`
	BirthDate    time.Time
	BirthCountry string
	BirthState   string
	BirthCity    string
	DeathDate    time.Time
	DeathCountry string
	DeathState   string
	DeathCity    string
	NameFirst    string
	NameLast     string
	NameGiven    string
	Weight       int
	Height       int
	Bats         string
	Throws       string
	Debut        time.Time
	FinalGame    time.Time
	RetroID      string
	BbrefID      string
}

func (p Player) MarshalJSON() ([]byte, error) {
	deathY, deathM, deathD := p.getDeathDates()
	birthY, birthM, birthD := p.getBirthDates()
	return json.Marshal(
		&struct {
			PlayerID     string `json:"playerID"`
			BirthDay     *int
			BirthYear    *int
			BirthMonth   *int
			BirthCountry string
			BirthState   string
			BirthCity    string
			DeathDate    *int
			DeathYear    *int
			DeathMonth   *int
			DeathCountry *string
			DeathState   *string
			DeathCity    *string
			NameFirst    string
			NameLast     string
			NameGiven    string
			Weight       int
			Height       int
			Bats         string
			Throws       string
			Debut        string
			FinalGame    string
			RetroID      string
			BbrefID      string
		}{
			PlayerID:     p.PlayerID,
			BirthDay:     birthD,
			BirthYear:    birthY,
			BirthMonth:   birthM,
			BirthCountry: p.BirthCountry,
			BirthState:   p.BirthState,
			BirthCity:    p.BirthCity,
			DeathDate:    deathD,
			DeathMonth:   deathM,
			DeathYear:    deathY,
			DeathCountry: getString(p.DeathCountry),
			DeathCity:    getString(p.DeathCity),
			DeathState:   getString(p.DeathState),
			NameGiven:    p.NameGiven,
			NameLast:     p.NameLast,
			NameFirst:    p.NameFirst,
			Weight:       p.Weight,
			Height:       p.Height,
			Bats:         p.Bats,
			Throws:       p.Throws,
			Debut:        p.Debut.Format(layout),
			FinalGame:    p.FinalGame.Format(layout),
			RetroID:      p.RetroID,
			BbrefID:      p.BbrefID,
		})

}

func (p *Player) getDeathDates() (*int, *int, *int) {

	if p.DeathDate.IsZero() {
		return nil, nil, nil
	}
	year := p.DeathDate.Year()
	day := p.DeathDate.Day()
	month := p.DeathDate.Month()
	monthConverted := int(month)
	return &year, &monthConverted, &day

}

func (p *Player) getBirthDates() (*int, *int, *int) {
	year := p.BirthDate.Year()
	day := p.BirthDate.Day()
	month := p.BirthDate.Month()
	monthConverted := int(month)
	return &year, &monthConverted, &day
}

func getString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
