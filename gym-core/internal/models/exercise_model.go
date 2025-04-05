package models

// Модель упражнения
type Exercise struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Type        string    `json:"type"`
	Difficulty  string    `json:"difficulty"`
	Equipment   *string   `json:"equipment"`
	VideoUrl    *string   `json:"video_url"`
	PhotoUrls   *[]string `json:"photo_urls"`
	Muscles     []Muscle  `json:"muscles"` // Связанные мышцы
}

// Модель мышцы
type Muscle struct {
	Id              int     `json:"id"`
	Name            string  `json:"name"`
	Description     *string `json:"description"`
	Photo           *string `json:"photo"`
	MusclesInvolved float64 `json:"muscles_involved"` // Степень вовлеченности мышцы (NOT NULL в таблице exercise_muscle)
}
