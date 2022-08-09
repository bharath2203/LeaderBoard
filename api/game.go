package api

import "encoding/json"

type GameInstance struct {
	GameId   string `json:"game_id"`
	UserName string `json:"user_name"`
}

func (gameInstance GameInstance) MarshalBinary() ([]byte, error) {
	return json.Marshal(gameInstance)
}

func (gameInstance *GameInstance) UnMarshalBinary(data []byte) error {
	return json.Unmarshal(data, &gameInstance)
}
