package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"sso-service/internal/domain"
)

type PostgresUserRepo struct {
	db *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) CreateUser(ctx context.Context, user domain.User) (int, error) {
	query := `INSERT INTO users (login, hash_password, role, email, address, phonenumber, first_name, last_name ) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING user_id`

	var userID int

	err := r.db.QueryRowContext(ctx, query, user.Login, user.HashPassword, user.Role, user.Email, user.Address, user.Phonenumber, user.FirstName, user.LastName).Scan(&userID)
	return userID, err
}

func (r *PostgresUserRepo) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	query := `SELECT user_id, login, hash_password, role, email, address, phonenumber, first_name, last_name, avatar_path
	FROM users WHERE login = $1`

	var user domain.User
	var avatar sql.NullString

	err := r.db.QueryRowContext(ctx, query, username).Scan(&user.UserID, &user.Login, &user.HashPassword, &user.Role,
		&user.Email, &user.Address, &user.Phonenumber, &user.FirstName, &user.LastName, &avatar)

	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user not found")
		}
		slog.Debug("Помилка при отриманні користувача з БД", "err", err.Error())
		return user, err
	}

	if avatar.Valid {
		user.AvatarPath = avatar.String
	} else {
		user.AvatarPath = ""
	}

	return user, nil
}

func (r *PostgresUserRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `SELECT user_id, login, hash_password, role, email, address, phonenumber, first_name, last_name, avatar_path
	FROM users WHERE email = $1`

	var user domain.User
	var avatar sql.NullString

	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.UserID, &user.Login, &user.HashPassword, &user.Role,
		&user.Email, &user.Address, &user.Phonenumber, &user.FirstName, &user.LastName, &avatar)

	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user not found")
		}
		slog.Debug("Помилка при отриманні користувача з БД", "err", err.Error())
		return user, err
	}

	if avatar.Valid {
		user.AvatarPath = avatar.String
	} else {
		user.AvatarPath = ""
	}

	return user, nil
}

func (r *PostgresUserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *PostgresUserRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE login = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, username).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *PostgresUserRepo) UpdateUserProfile(ctx context.Context, userData domain.UserUpdateRequest) error {
	baseQuery := `
	UPDATE users SET 
		login = $2,
		hash_password = $3,
		first_name = $4,
		last_name = $5,
		email = $6,
		phonenumber = $7,
		address = $8
		%s
	WHERE user_id = $1;
	`

	args := []any{
		userData.UserID,
		userData.Login,
		userData.HashPassword,
		userData.FirstName,
		userData.LastName,
		userData.Email,
		userData.Phonenumber,
		userData.Address,
	}

	avatarQuery := ""
	if userData.AvatarPath != "" {
		avatarQuery = ", avatar_path = $9"
		args = append(args, userData.AvatarPath)
	}

	finalQuery := fmt.Sprintf(baseQuery, avatarQuery)

	_, err := s.db.Exec(finalQuery, args...)
	if err != nil {
		slog.Debug("Помилка при оновленні профіля користувача", "err", err.Error())
		return err
	}

	return nil
}
