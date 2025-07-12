package service

import (
	"context"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"sso/internal/domain"
	"sso/pkg/auth"

	"github.com/google/uuid"
)

type UsersService struct {
	repo domain.UserRepository
}

func NewUsersService(repo domain.UserRepository) *UsersService {
	return &UsersService{repo: repo}
}

func (s *UsersService) SaveAvatarFile(userID int, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	dir := "internal/storage/avatars"
	os.MkdirAll(dir, os.ModePerm)

	filename := uuid.New().String() + filepath.Ext(fileHeader.Filename)
	path := filepath.Join(dir, filename)

	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return filename, nil
}

func (s *UsersService) CreateUser(ctx context.Context, req *domain.RegisterRequest) error {
	exists, err := s.repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user already exists")
	}

	hashedPwd, err := auth.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Login:         req.Login,
		Email:         req.Email,
		Role:          "user",
		HashPassword:  hashedPwd,
		Address:       req.Address, 
		Phone:   req.Phonenumber,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
	}

  userID, err := s.repo.CreateUser(ctx, user)

  req.UserID = userID
	return err
}

func (s *UsersService) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	return s.repo.GetByUsername(ctx, username)
}

func (s *UsersService) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return s.repo.ExistsByEmail(ctx, email)
}

func (s *UsersService) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	return s.repo.ExistsByUsername(ctx, username)
}

func (s *UsersService) UpdateUserProfile(ctx context.Context, userData domain.UserUpdateRequest) error {
	hashPassword, err := auth.HashPassword(userData.Password)
	if err != nil {
		log.Println("Ошибка при хешировании пароля:", err)
		return err
	}

	userData.HashPassword = hashPassword

	return s.repo.UpdateUserProfile(ctx, userData)
}
