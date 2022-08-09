package api

import "encoding/json"

type GameScore struct {
	GameId    string  `json:"game_id" validate:"required"`
	UserName  string  `json:"user_name" validate:"required"`
	UserScore float64 `json:"user_score" validate:"required"`
}

func (score GameScore) MarshalBinary() ([]byte, error) {
	return json.Marshal(score)
}

func (score *GameScore) UnMarshalBinary(data []byte) error {
	return json.Unmarshal(data, &score)
}
