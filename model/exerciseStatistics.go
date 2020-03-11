package model

type ExerciseStatistics struct {
	LessonId           string `json:""`
	UserId             string `json:""`
	ExerciseId         string `json:""`
	SuccessfulAttempts int32 `json:""`
	FailedAttempts     int32 `json:""`
	Completed          bool  `json:""`
}
