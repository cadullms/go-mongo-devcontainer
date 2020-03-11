package model

type Exercise struct {
	Id string `json:"exerciseId"`
    Question string `json:"question"`
    ExpectedAnswer string `json:"expectedAnswer"`
}