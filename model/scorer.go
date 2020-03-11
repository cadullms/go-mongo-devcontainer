package model

func ScoreAttemptsForLesson(lesson *Lesson, userId string, attempts *[]Attempt) (LessonStatistics, []ExerciseStatistics) {
	var stats = LessonStatistics{}
	stats.LessonId = lesson.Id
	stats.UserId = userId
	var exerciseStats []ExerciseStatistics
	for _, exercise := range lesson.Exercises {
		subStat := ScoreAttemptsForExercise(lesson, userId, exercise.Id, attempts)
		if subStat.Completed {
			stats.CompletedExercises++
		}
		if subStat.SuccessfulAttempts > 0 || subStat.FailedAttempts > 0 {
			stats.StartedExercises++
		}
		exerciseStats = append(exerciseStats, subStat)
	}
	stats.TotalExercises = int32(len(lesson.Exercises))
	return stats, exerciseStats
}

func ScoreAttemptsForExercise(lesson *Lesson, userId string, exerciseId string, attempts *[]Attempt) ExerciseStatistics {
	var stats = ExerciseStatistics{}
	stats.LessonId = lesson.Id
	stats.UserId = userId
	stats.ExerciseId = exerciseId
	for _, attempt := range *attempts {
		if attempt.LessonId == lesson.Id && attempt.UserId == userId && attempt.ExerciseId == exerciseId {
			switch attempt.State {
			case AttemptStateCorrect:
				stats.SuccessfulAttempts++
			case AttemptStateIncorrect:
				stats.FailedAttempts++
			}
		}
	}
	stats.Completed = stats.SuccessfulAttempts > 3
	return stats
}
