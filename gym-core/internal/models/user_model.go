package models

import "time"

type User struct {
	Id            int        `json:"id"`
	LastName      string     `json:"last_name"`
	FirstName     string     `json:"first_name"`
	Gender        *string    `json:"gender"`
	Birthday      *time.Time `json:"birthday"`
	Weight        *float64   `json:"weight"`
	Height        *float64   `json:"height"`
	Purpose       *string    `json:"purpose"`
	Placement     *string    `json:"placement"`
	Level         *string    `json:"level"`
	TrainingCount *int       `json:"training_count"`
	FoodPrompt    *string    `json:"food_prompt"`
}

type UserUpdateRequest struct {
	LastName      string    `json:"last_name" binding:"required"`
	FirstName     string    `json:"first_name" binding:"required"`
	Gender        string    `json:"gender" binding:"required"`
	Birthday      time.Time `json:"birthday" binding:"required"`
	Weight        float64   `json:"weight" binding:"required"`
	Height        float64   `json:"height" binding:"required"`
	Purpose       *string   `json:"purpose"`
	Placement     *string   `json:"placement"`
	Level         *string   `json:"level"`
	TrainingCount *int      `json:"training_count"`
	FoodPrompt    *string   `json:"food_prompt"`
}
