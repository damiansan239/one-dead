package datastore

import (
	"time"
)

type Player struct {
	Name string
}

type Game struct {
	Tries    int
	Duration int
	Won      bool
	Id       string
	Players  []*Player
	Date     time.Time
}
