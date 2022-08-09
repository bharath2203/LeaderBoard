package gamescore

import "TopKScores/api"

type Observer interface {
	Update(gameScore *api.GameScore)
	GetId() string
}
