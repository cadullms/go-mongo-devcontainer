package model

import (
	"time"
)

type Attempt struct {
    Id string `json:"_id"`
	UserId string `json:"userId"`
    LessonId string `json:"lessonId"`
    ExerciseId string `json:"exerciseId"` 
    ActualAnswer string `json:"actualAnswer"`
    TimeToAnswerMs int32 `json:"timeToAnswerMs"`
    TimeStarted time.Time `json:"timeStarted"`
    State AttemptState `json:"state"`
}