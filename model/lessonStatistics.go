package model

type LessonStatistics struct {
	LessonId string `json:""`
	UserId string `json:""`
	TotalExercises int32  `json:""`
	StartedExercises int32  `json:""`
	CompletedExercises int32  `json:""`
}
