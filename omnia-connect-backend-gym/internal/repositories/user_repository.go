package repositories

import (
	"api-gateway/internal/domain"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

var (
	ErrorUserNotFound     = errors.New("user not found")
	ErrorUserAlreadyExist = errors.New("user already exist")
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
	//GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
	GetUserByGoogleID(ctx context.Context, id string) (*domain.User, error)
	GetUserByVkID(ctx context.Context, id string) (*domain.User, error)
	//UpdateUser(ctx context.Context, user *domain.User) error
	CreateUser(ctx context.Context, user *domain.User) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User
	query := "SELECT id, last_name, first_name, email, roles FROM users WHERE id = $1"
	err := ur.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.LastName, &user.FirstName, &user.Email, &user.Roles)
	if err != nil {
		slog.ErrorContext(ctx, "failed get user by id", slog.String("error", err.Error()))
		return nil, err
	}
	slog.DebugContext(ctx, "succeed to get user by id", slog.Any("user", user))
	return &user, nil
}

//func (ur *userRepository) GetUserByLogin(ctx context.Context, login string) (*domain.User, error) {
//	var user domain.User
//	query := "SELECT id, fio, login, password, roles FROM users WHERE login = $1 and google_id is null and vk_id is null"
//	err := ur.db.QueryRow(ctx, query, login).Scan(&user.ID, &user.FIO, &user.Email, &user.Password, &user.Roles)
//	if err != nil {
//		if errors.Is(err, pgx.ErrNoRows) {
//			return nil, ErrorUserNotFound
//		}
//		slog.ErrorContext(ctx, "failed get user by login", slog.String("error", err.Error()))
//		return nil, err
//	}
//	slog.DebugContext(ctx, "succeed to get user by login", slog.Any("user", user))
//	return &user, nil
//}

func (ur *userRepository) GetUserByGoogleID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	query := "SELECT id, last_name, first_name, email, roles FROM users WHERE google_id = $1"
	err := ur.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.LastName, &user.FirstName, &user.Email, &user.Roles)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrorUserNotFound
		}
		slog.ErrorContext(ctx, "failed get user by login", slog.String("error", err.Error()))
		return nil, err
	}
	slog.DebugContext(ctx, "succeed to get user by login", slog.Any("user", user))
	return &user, nil
}

func (ur *userRepository) GetUserByVkID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	query := "SELECT id, last_name, first_name, email, roles FROM users WHERE vk_id = $1"
	err := ur.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.LastName, &user.FirstName, &user.Email, &user.Roles)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrorUserNotFound
		}
		slog.ErrorContext(ctx, "failed get user by login", slog.String("error", err.Error()))
		return nil, err
	}
	slog.DebugContext(ctx, "succeed to get user by login", slog.Any("user", user))
	return &user, nil
}

//func (ur *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
//	query := "UPDATE users SET fio = &1, login = $2, password = $3, roles = $3 WHERE id = $5"
//	_, err := ur.db.Exec(ctx, query, user.FIO, user.Email, user.Password, user.Roles, user.ID)
//	if err != nil {
//		slog.ErrorContext(ctx, "failed update user", slog.String("error", err.Error()))
//		return err
//	}
//	return nil
//}

//func (ur *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
//
//	query := "INSERT INTO users (login, password, roles) VALUES ($1, $2, $3) returning id"
//	_, err := ur.db.Exec(ctx, query, user.Email, user.Password, user.Roles)
//	if err != nil {
//		slog.ErrorContext(ctx, "failed update user", slog.String("error", err.Error()))
//		return err
//	}
//	return nil
//}

func (ur *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := "INSERT INTO users (last_name, first_name, email, roles, google_id, vk_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

	var id int
	err := ur.db.QueryRow(ctx, query, user.LastName, user.FirstName, user.Email, user.Roles, user.GoogleID, user.VkID).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "users_login_key" {
				return ErrorUserAlreadyExist
			}
		}
		slog.ErrorContext(ctx, "failed to create user", slog.String("error", err.Error()))
		return err
	}

	user.ID = id

	return nil
}
