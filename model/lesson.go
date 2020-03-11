package model

type Lesson struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Exercises []Exercise `json:"exercises"`
}
