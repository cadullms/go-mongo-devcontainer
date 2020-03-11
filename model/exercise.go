package model

type Exercise struct {
	ExerciseId string `json:"exerciseId"`
    Question string `json:"question"`
    ExpectedAnswer string `json:"expectedAnswer"`
}