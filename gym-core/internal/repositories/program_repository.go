package repositories

import (
	"context"
	"database/sql"
	"errors"
	"gym-core/internal/models"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrorProgramNotFound = errors.New("program not found")
)

type ProgramRepository interface {
	GetProgramByUserId(ctx context.Context, userId int) (*models.Program, error)
}

type programRepository struct {
	db *pgxpool.Pool
}

func NewProgramRepository(db *pgxpool.Pool) ProgramRepository {
	return &programRepository{db: db}
}

func (r *programRepository) GetProgramByUserId(ctx context.Context, userId int) (*models.Program, error) {
	// Шаг 1: Получаем программу пользователя
	query := `
		SELECT id, user_id, name, description
		FROM programs
		WHERE user_id = $1
	`

	var program models.Program
	var description sql.NullString

	err := r.db.QueryRow(ctx, query, userId).Scan(
		&program.Id,
		&program.UserId,
		&program.Name,
		&description,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrorProgramNotFound
		}
		slog.Error("Failed to get program", "error", err)
		return nil, err
	}

	// Преобразуем null-значения в указатели
	if description.Valid {
		program.Description = &description.String
	}

	// Шаг 2: Получаем все тренировки для этой программы
	workoutsQuery := `
		SELECT id, program_id, name, description, date_time, status, notes
		FROM workouts
		WHERE program_id = $1
		ORDER BY date_time
	`

	rows, err := r.db.Query(ctx, workoutsQuery, program.Id)
	if err != nil {
		slog.Error("Failed to get workouts for program", "error", err, "programId", program.Id)
		return &program, err
	}
	defer rows.Close()

	program.Workouts = []models.Workout{}
	workoutIds := make([]int, 0)

	// Собираем тренировки
	for rows.Next() {
		var workout models.Workout
		var workoutDescription, notes sql.NullString
		var dateTime sql.NullTime
		var programId sql.NullInt32

		err := rows.Scan(
			&workout.Id,
			&programId,
			&workout.Name,
			&workoutDescription,
			&dateTime,
			&workout.Status,
			&notes,
		)

		if err != nil {
			slog.Error("Failed to scan workout", "error", err)
			continue
		}

		if programId.Valid {
			id := int(programId.Int32)
			workout.ProgramId = &id
		}

		if workoutDescription.Valid {
			workout.Description = &workoutDescription.String
		}

		if dateTime.Valid {
			workout.DateTime = &dateTime.Time
		}

		if notes.Valid {
			workout.Notes = &notes.String
		}

		workout.Exercises = []models.WorkoutExercise{}
		program.Workouts = append(program.Workouts, workout)
		workoutIds = append(workoutIds, workout.Id)
	}

	if err = rows.Err(); err != nil {
		slog.Error("Error iterating workout rows", "error", err)
		return &program, err
	}

	// Если нет тренировок, возвращаем программу
	if len(workoutIds) == 0 {
		return &program, nil
	}

	// Шаг 3: Получаем все упражнения для всех тренировок
	// Создаем карту для быстрого доступа к тренировкам
	workoutMap := make(map[int]*models.Workout)
	for i := range program.Workouts {
		workoutMap[program.Workouts[i].Id] = &program.Workouts[i]
	}

	workoutExercisesQuery := `
		SELECT we.id, we.workout_id, we.exercise_id, we.order_number,
			   e.id as e_id, e.name as e_name, e.description as e_description, 
			   e.type as e_type, e.difficulty as e_difficulty, 
			   e.equipment as e_equipment, e.video_url as e_video_url, 
			   e.photo_urls as e_photo_urls
		FROM workout_exercise we
		JOIN exercises e ON we.exercise_id = e.id
		WHERE we.workout_id = ANY($1)
		ORDER BY we.workout_id, we.order_number
	`

	weRows, err := r.db.Query(ctx, workoutExercisesQuery, workoutIds)
	if err != nil {
		slog.Error("Failed to get workout exercises", "error", err)
		return &program, err
	}
	defer weRows.Close()

	// Карта ID упражнений в тренировке
	workoutExerciseIds := make([]int, 0)
	workoutExerciseMap := make(map[int]*models.WorkoutExercise)

	// Карта ID упражнений и упражнений для загрузки мышц
	exerciseIds := make([]int, 0)
	exerciseMap := make(map[int]*models.Exercise)

	// Собираем упражнения для тренировок
	for weRows.Next() {
		var we models.WorkoutExercise
		var e models.Exercise
		var eDescription, eEquipment, eVideoUrl sql.NullString
		var photoUrls []string

		err := weRows.Scan(
			&we.Id,
			&we.WorkoutId,
			&we.ExerciseId,
			&we.OrderNumber,
			&e.Id,
			&e.Name,
			&eDescription,
			&e.Type,
			&e.Difficulty,
			&eEquipment,
			&eVideoUrl,
			&photoUrls,
		)

		if err != nil {
			slog.Error("Failed to scan workout exercise", "error", err)
			continue
		}

		// Заполняем данные упражнения
		if eDescription.Valid {
			e.Description = &eDescription.String
		}

		if eEquipment.Valid {
			e.Equipment = &eEquipment.String
		}

		if eVideoUrl.Valid {
			e.VideoUrl = &eVideoUrl.String
		}

		if len(photoUrls) > 0 {
			e.PhotoUrls = &photoUrls
		}

		// Инициализируем пустой слайс мышц
		e.Muscles = []models.Muscle{}

		// Присваиваем упражнение
		we.Exercise = e
		we.Sets = []models.ExerciseSet{}

		// Добавляем в соответствующую тренировку
		if workout, ok := workoutMap[we.WorkoutId]; ok {
			workout.Exercises = append(workout.Exercises, we)
			workoutExerciseIds = append(workoutExerciseIds, we.Id)
			workoutExerciseMap[we.Id] = &workout.Exercises[len(workout.Exercises)-1]

			// Добавляем ID упражнения для загрузки мышц
			exerciseIds = append(exerciseIds, e.Id)
			exerciseMap[e.Id] = &workout.Exercises[len(workout.Exercises)-1].Exercise
		}
	}

	if err = weRows.Err(); err != nil {
		slog.Error("Error iterating workout exercise rows", "error", err)
		return &program, err
	}

	// Если нет упражнений, возвращаем программу
	if len(workoutExerciseIds) == 0 {
		return &program, nil
	}

	// Шаг 3.5: Получаем все мышцы для всех упражнений
	if len(exerciseIds) > 0 {
		musclesQuery := `
			SELECT em.exercise_id, m.id, m.name, m.description, m.photo, em.muscles_involved
			FROM muscles m
			JOIN exercise_muscle em ON m.id = em.muscle_id
			WHERE em.exercise_id = ANY($1)
			ORDER BY em.exercise_id, m.name
		`

		muscleRows, err := r.db.Query(ctx, musclesQuery, exerciseIds)
		if err != nil {
			slog.Error("Failed to get muscles for exercises", "error", err)
			// Продолжаем выполнение даже если не удалось получить мышцы
		} else {
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
				if exercise, ok := exerciseMap[exerciseId]; ok {
					exercise.Muscles = append(exercise.Muscles, muscle)
				}
			}

			if err = muscleRows.Err(); err != nil {
				slog.Error("Error iterating muscle rows", "error", err)
			}
		}
	}

	// Шаг 4: Получаем все подходы для всех упражнений
	setsQuery := `
		SELECT id, workout_exercise_id, set_number, 
		       planned_sets, planned_reps, planned_weight,
		       completed_sets, completed_reps, completed_weight, is_completed
		FROM exercise_set
		WHERE workout_exercise_id = ANY($1)
		ORDER BY workout_exercise_id, set_number
	`

	setRows, err := r.db.Query(ctx, setsQuery, workoutExerciseIds)
	if err != nil {
		slog.Error("Failed to get exercise sets", "error", err)
		return &program, err
	}
	defer setRows.Close()

	// Собираем подходы для упражнений
	for setRows.Next() {
		var set models.ExerciseSet
		var plannedSets, plannedReps, completedSets, completedReps sql.NullInt32
		var plannedWeight, completedWeight sql.NullFloat64

		err := setRows.Scan(
			&set.Id,
			&set.WorkoutExerciseId,
			&set.SetNumber,
			&plannedSets,
			&plannedReps,
			&plannedWeight,
			&completedSets,
			&completedReps,
			&completedWeight,
			&set.IsCompleted,
		)

		if err != nil {
			slog.Error("Failed to scan exercise set", "error", err)
			continue
		}

		if plannedSets.Valid {
			val := int(plannedSets.Int32)
			set.PlannedSets = &val
		}

		if plannedReps.Valid {
			val := int(plannedReps.Int32)
			set.PlannedReps = &val
		}

		if plannedWeight.Valid {
			set.PlannedWeight = &plannedWeight.Float64
		}

		if completedSets.Valid {
			val := int(completedSets.Int32)
			set.CompletedSets = &val
		}

		if completedReps.Valid {
			val := int(completedReps.Int32)
			set.CompletedReps = &val
		}

		if completedWeight.Valid {
			set.CompletedWeight = &completedWeight.Float64
		}

		// Добавляем в соответствующее упражнение
		if we, ok := workoutExerciseMap[set.WorkoutExerciseId]; ok {
			we.Sets = append(we.Sets, set)
		}
	}

	if err = setRows.Err(); err != nil {
		slog.Error("Error iterating exercise set rows", "error", err)
	}

	return &program, nil
}
