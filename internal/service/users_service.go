package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"sso-service/internal/domain"
	"sso-service/pkg/auth"
	"time"

	"github.com/google/uuid"
)

type UsersService struct {
	repo              domain.UserRepository
	storageServiceURL string
	httpClient        http.Client
}

func NewUsersService(repo domain.UserRepository, storageServiceURL string) *UsersService {
	return &UsersService{
		repo:              repo,
		storageServiceURL: storageServiceURL,
		httpClient:        http.Client{Timeout: 10 * time.Second},
	}
}

func (s *UsersService) StorageRequest(requestURL string, requestBody *bytes.Buffer, contentType string) error {
	req, err := http.NewRequest("POST", requestURL, requestBody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("помилка запиту до storage: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("storage повернув помилку: %d, %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func (s *UsersService) SaveAvatar(fileHeader *multipart.FileHeader) (string, error) {
	var newReqBody bytes.Buffer
	writer := multipart.NewWriter(&newReqBody)

	file, err := fileHeader.Open()
	if err != nil {
		slog.Debug("Помилка відкриття файлу", "err", err.Error())
		return "", err
	}
	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)
	generatedName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	part, err := writer.CreateFormFile("file", generatedName)
	if err != nil {
		slog.Debug("Помилка створення форми", "err", err.Error())
		return "", err
	}

	if _, err := io.Copy(part, file); err != nil {
		slog.Debug("Помилка копіювання файлу", "err", err.Error())
		return "", err
	}

	err = writer.Close()
	if err != nil {
		slog.Debug("Помилка закриття writer", "err", err.Error())
		return "", err
	}

	requestURL := fmt.Sprintf("%s/api/storage/upload_avatar", s.storageServiceURL)
	if err := s.StorageRequest(requestURL, &newReqBody, writer.FormDataContentType()); err != nil {
		slog.Debug("Помилка при збереженні аватара на сервісі storage", "err", err.Error())
		return "", err
	}

	return generatedName, nil
}

func (s *UsersService) DeleteAvatar(filename string) error {
	payload, err := json.Marshal(map[string]string{
		"filename": filename,
	})
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf("%s/api/storage/delete_avatar", s.storageServiceURL)
	s.StorageRequest(requestURL, bytes.NewBuffer(payload), "application/json")

	return nil
}

func (s *UsersService) CreateUser(ctx context.Context, req *domain.RegisterRequest) error {
	exists, err := s.repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("user already exists")
	}

	hashedPwd, err := auth.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Login:        req.Login,
		Email:        req.Email,
		Role:         "user",
		HashPassword: hashedPwd,
		Address:      req.Address,
		Phonenumber:  req.Phonenumber,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
	}

	userID, err := s.repo.CreateUser(ctx, user)

	req.UserID = userID
	return err
}

func (s *UsersService) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *UsersService) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return s.repo.ExistsByEmail(ctx, email)
}

func (s *UsersService) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	return s.repo.GetByUsername(ctx, username)
}

func (s *UsersService) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	return s.repo.ExistsByUsername(ctx, username)
}

func (s *UsersService) UpdateUserProfile(ctx context.Context, userData domain.UserUpdateRequest) error {
	hashPassword, err := auth.HashPassword(userData.Password)
	if err != nil {
		slog.Debug("Ошибка при хешировании пароля:", "err", err.Error())
		return err
	}

	userData.HashPassword = hashPassword

	return s.repo.UpdateUserProfile(ctx, userData)
}
