package model

type Lesson struct {
	Id string `json:"_id"`
	Name string `json:"name"`
	Exercises []Exercise `json:"exercises"`
}
