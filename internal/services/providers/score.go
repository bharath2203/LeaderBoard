package providers

import (
	"TopKScores/api"
)

// TopKScoreProvider is an interface that exposes two functionalities.
// 1. Add/Update a GameScore to the stream.
// 2. Return top K scores objects sorted in score, descending order.
type TopKScoreProvider interface {
	AddScore(gameScore *api.GameScore) error
	GetTopKScores(k int64) ([]*api.GameScore, error)
}
