package repositories

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"gym-core/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrorExerciseNotFound = errors.New("exercise not found")
)

type ExerciseRepository interface {
	GetExerciseById(ctx context.Context, id int) (*models.Exercise, error)
	GetAllExercises(ctx context.Context) ([]models.Exercise, error)
}

type exerciseRepository struct {
	db *pgxpool.Pool
}

func NewExerciseRepository(db *pgxpool.Pool) ExerciseRepository {
	return &exerciseRepository{db: db}
}

func (r *exerciseRepository) GetExerciseById(ctx context.Context, id int) (*models.Exercise, error) {
	// Шаг 1: Получаем данные о упражнении
	query := `
		SELECT id, name, description, type, difficulty, equipment, video_url, photo_urls
		FROM exercises
		WHERE id = $1
	`

	var exercise models.Exercise
	var description, equipment, videoUrl sql.NullString
	var photoUrls []string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&exercise.Id,
		&exercise.Name,
		&description,
		&exercise.Type,
		&exercise.Difficulty,
		&equipment,
		&videoUrl,
		&photoUrls,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrorExerciseNotFound
		}
		slog.Error("Failed to get exercise", "error", err)
		return nil, err
	}

	// Преобразуем null-значения в указатели
	if description.Valid {
		exercise.Description = &description.String
	}

	if equipment.Valid {
		exercise.Equipment = &equipment.String
	}

	if videoUrl.Valid {
		exercise.VideoUrl = &videoUrl.String
	}

	if len(photoUrls) > 0 {
		exercise.PhotoUrls = &photoUrls
	}

	// Шаг 2: Получаем связанные мышцы для этого упражнения
	musclesQuery := `
		SELECT m.id, m.name, m.description, m.photo, em.muscles_involved
		FROM muscles m
		JOIN exercise_muscle em ON m.id = em.muscle_id
		WHERE em.exercise_id = $1
	`

	// Используем pgx.Rows вместо sql.Rows
	rows, err := r.db.Query(ctx, musclesQuery, id)
	if err != nil {
		slog.Error("Failed to get muscles for exercise", "error", err, "exerciseId", id)
		return &exercise, err
	}
	defer rows.Close()

	exercise.Muscles = []models.Muscle{}

	// Создаем слайс для хранения всех мышц и используем pgx.CollectRows
	var muscles []models.Muscle
	muscles, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.Muscle, error) {
		var muscle models.Muscle
		var muscleDescription, photo sql.NullString

		err := row.Scan(
			&muscle.Id,
			&muscle.Name,
			&muscleDescription,
			&photo,
			&muscle.MusclesInvolved,
		)

		if err != nil {
			return models.Muscle{}, err
		}

		if muscleDescription.Valid {
			muscle.Description = &muscleDescription.String
		}

		if photo.Valid {
			muscle.Photo = &photo.String
		}

		return muscle, nil
	})

	if err != nil {
		slog.Error("Failed to collect muscles", "error", err)
		// Возвращаем упражнение без мышц в случае ошибки
		return &exercise, err
	}

	exercise.Muscles = muscles
	return &exercise, nil
}

// Добавляем новый метод для получения всех упражнений
func (r *exerciseRepository) GetAllExercises(ctx context.Context) ([]models.Exercise, error) {
	// Шаг 1: Получаем все упражнения
	query := `
		SELECT id, name, description, type, difficulty, equipment, video_url, photo_urls
		FROM exercises
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		slog.Error("Failed to get exercises", "error", err)
		return nil, err
	}
	defer rows.Close()

	var exercises []models.Exercise

	exercises, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.Exercise, error) {
		var exercise models.Exercise
		var description, equipment, videoUrl sql.NullString
		var photoUrls []string

		err := row.Scan(
			&exercise.Id,
			&exercise.Name,
			&description,
			&exercise.Type,
			&exercise.Difficulty,
			&equipment,
			&videoUrl,
			&photoUrls,
		)

		if err != nil {
			return models.Exercise{}, err
		}

		// Преобразуем null-значения в указатели
		if description.Valid {
			exercise.Description = &description.String
		}

		if equipment.Valid {
			exercise.Equipment = &equipment.String
		}

		if videoUrl.Valid {
			exercise.VideoUrl = &videoUrl.String
		}

		if len(photoUrls) > 0 {
			exercise.PhotoUrls = &photoUrls
		}

		return exercise, nil
	})

	if err != nil {
		slog.Error("Failed to collect exercises", "error", err)
		return nil, err
	}

	// Шаг 2: Получаем все мышцы для всех упражнений за один запрос
	if len(exercises) > 0 {
		// Собираем все ID упражнений
		exerciseIds := make([]int, len(exercises))
		exerciseMap := make(map[int]*models.Exercise, len(exercises))

		for i, ex := range exercises {
			exerciseIds[i] = ex.Id
			exerciseMap[ex.Id] = &exercises[i]
			exercises[i].Muscles = []models.Muscle{} // Инициализируем пустой слайс мышц
		}

		// Делаем один запрос для получения всех мышц для всех упражнений
		musclesQuery := `
			SELECT em.exercise_id, m.id, m.name, m.description, m.photo, em.muscles_involved
			FROM muscles m
			JOIN exercise_muscle em ON m.id = em.muscle_id
			WHERE em.exercise_id = ANY($1)
		`

		muscleRows, err := r.db.Query(ctx, musclesQuery, exerciseIds)
		if err != nil {
			slog.Error("Failed to get muscles for exercises", "error", err)
			return exercises, err
		}
		defer muscleRows.Close()

		for muscleRows.Next() {
			var exerciseId int
			var muscle models.Muscle
			var muscleDescription, photo sql.NullString

			err := muscleRows.Scan(
				&exerciseId,
				&muscle.Id,
				&muscle.Name,
				&muscleDescription,
				&photo,
				&muscle.MusclesInvolved,
			)

			if err != nil {
				slog.Error("Failed to scan muscle", "error", err)
				continue
			}

			if muscleDescription.Valid {
				muscle.Description = &muscleDescription.String
			}

			if photo.Valid {
				muscle.Photo = &photo.String
			}

			// Добавляем мышцу к соответствующему упражнению
			if ex, ok := exerciseMap[exerciseId]; ok {
				ex.Muscles = append(ex.Muscles, muscle)
			}
		}

		if err = muscleRows.Err(); err != nil {
			slog.Error("Error iterating muscle rows", "error", err)
		}
	}

	return exercises, nil
}
