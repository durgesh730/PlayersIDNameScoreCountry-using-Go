package models

type Player struct {
	Id      string `json:"id,omitempty" bson:"id,omitempty"`
	Name    string `json:"name" bson:"name"`
	Country string `json:"country" bson:"country"`
	Score   int    `json:"score" bson:"score"`
}
