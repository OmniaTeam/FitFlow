package repositories

import (
	"context"
	"database/sql"
	"errors"
	"gym-core/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrorUserNotFound = errors.New("user not found")
)

type UserRepository interface {
	GetUserById(ctx context.Context, id int) (*models.User, error)
	UpdateUser(ctx context.Context, id int, req models.UserUpdateRequest) (*models.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserById(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, last_name, first_name, gender, birthday, weight, height,
		       purpose, placement, level, training_count, food_prompt
		FROM users
		WHERE id = $1
	`

	var user models.User
	var gender sql.NullString
	var birthday sql.NullTime
	var weight, height sql.NullFloat64
	var purpose, placement, level, foodPrompt sql.NullString
	var trainingCount sql.NullInt32

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.Id,
		&user.LastName,
		&user.FirstName,
		&gender,
		&birthday,
		&weight,
		&height,
		&purpose,
		&placement,
		&level,
		&trainingCount,
		&foodPrompt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrorUserNotFound
		}
		return nil, err
	}

	if gender.Valid {
		user.Gender = &gender.String
	}

	if birthday.Valid {
		user.Birthday = &birthday.Time
	}

	if weight.Valid {
		user.Weight = &weight.Float64
	}

	if height.Valid {
		user.Height = &height.Float64
	}

	if purpose.Valid {
		user.Purpose = &purpose.String
	}

	if placement.Valid {
		user.Placement = &placement.String
	}

	if level.Valid {
		user.Level = &level.String
	}

	if trainingCount.Valid {
		trainingCountInt := int(trainingCount.Int32)
		user.TrainingCount = &trainingCountInt
	}

	if foodPrompt.Valid {
		user.FoodPrompt = &foodPrompt.String
	}

	return &user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, id int, req models.UserUpdateRequest) (*models.User, error) {
	query := `
		UPDATE users
		SET last_name = $1, 
		    first_name = $2, 
		    gender = $3, 
		    birthday = $4, 
		    weight = $5, 
		    height = $6,
		    purpose = $7,
		    placement = $8,
		    level = $9,
		    training_count = $10,
		    food_prompt = $11
		WHERE id = $12
	`

	result, err := r.db.Exec(ctx, query,
		req.LastName,
		req.FirstName,
		req.Gender,
		req.Birthday,
		req.Weight,
		req.Height,
		req.Purpose,
		req.Placement,
		req.Level,
		req.TrainingCount,
		req.FoodPrompt,
		id,
	)

	if err != nil {
		return nil, err
	}

	if result.RowsAffected() == 0 {
		return nil, ErrorUserNotFound
	}

	// Получаем обновленного пользователя
	return r.GetUserById(ctx, id)
}
