package character

import "encoding/json"

func UnmarshalMetadata(data []byte) (*Character, error) {
	r := new(Character)
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Character) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Character struct {
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	Personality   string      `json:"personality"`
	FirstMes      string      `json:"first_mes"`
	Avatar        string      `json:"avatar"`
	Chat          string      `json:"chat"`
	MesExample    string      `json:"mes_example"`
	Scenario      string      `json:"scenario"`
	CreateDate    string      `json:"create_date"`
	Talkativeness json.Number `json:"talkativeness"`
}
