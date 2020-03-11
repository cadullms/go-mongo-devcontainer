package model

import (
	"time"
)

type Attempt struct {
	UserId string `json:"userId"`
    LessonId string `json:"lessonId"`
    ExerciseId string `json:"exerciseId"` 
    ActualAnswer string `json:"actualAnswer"`
    TimeToAnswerMs int32 `json:"timeToAnswerMs"`
    TimeStarted time.Time `json:"timeStarted"`
    State int32 `json:"state"`
}