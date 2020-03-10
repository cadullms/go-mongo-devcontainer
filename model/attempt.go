package model

import (
	"time"
)

type Attempt struct {
	userId string 
    lessonId string
    exerciseId string 
    actualAnswer string
    timeToAnswerMs int32
    timeStarted time.Time
    state string
}