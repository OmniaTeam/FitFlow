package models

import "time"

// Program представляет программу тренировок для пользователя
type Program struct {
	Id          int       `json:"id"`
	UserId      int       `json:"user_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Workouts    []Workout `json:"workouts"` // Связанные тренировки
}

// Workout представляет отдельную тренировку в программе
type Workout struct {
	Id          int               `json:"id"`
	ProgramId   *int              `json:"program_id"`
	Name        string            `json:"name"`
	Description *string           `json:"description"`
	DateTime    *time.Time        `json:"date_time"`
	Status      string            `json:"status"`
	Notes       *string           `json:"notes"`
	Exercises   []WorkoutExercise `json:"exercises"` // Связанные упражнения
}

// WorkoutExercise представляет упражнение в рамках конкретной тренировки
type WorkoutExercise struct {
	Id          int           `json:"id"`
	WorkoutId   int           `json:"workout_id"`
	ExerciseId  int           `json:"exercise_id"`
	OrderNumber int           `json:"order_number"`
	Exercise    Exercise      `json:"exercise"` // Данные об упражнении
	Sets        []ExerciseSet `json:"sets"`     // Подходы для этого упражнения
}

// ExerciseSet представляет подход для упражнения
type ExerciseSet struct {
	Id                int      `json:"id"`
	WorkoutExerciseId int      `json:"workout_exercise_id"`
	SetNumber         int      `json:"set_number"`
	PlannedSets       *int     `json:"planned_sets"`
	PlannedReps       *int     `json:"planned_reps"`
	PlannedWeight     *float64 `json:"planned_weight"`
	CompletedSets     *int     `json:"completed_sets"`
	CompletedReps     *int     `json:"completed_reps"`
	CompletedWeight   *float64 `json:"completed_weight"`
	IsCompleted       bool     `json:"is_completed"`
}
