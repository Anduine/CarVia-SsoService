package http_handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sso/internal/domain"
	"sso/internal/service"
	"sso/pkg/auth"
	"strings"
	"time"
)

type UsersHandler struct {
	log *slog.Logger
	service *service.UsersService
	tokenTTL time.Duration
}

func NewUsersHandler(log *slog.Logger, service *service.UsersService, tokenTTL time.Duration) *UsersHandler {
	return &UsersHandler{
		log: log,
		service: service, 
		tokenTTL: tokenTTL,
	}
}

func (h *UsersHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var regRequest domain.RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&regRequest)
	if err != nil {
		http.Error(w, "Неправильний запит", http.StatusBadRequest)
		return
	}

	exists, err := h.service.ExistsByEmail(r.Context(), regRequest.Email)
	if err != nil {
		h.log.Error("Ошибка в проверке: ", slog.Any("Error:", err))
		http.Error(w, "Помилка на сервері", http.StatusInternalServerError)
		return
	}
	if exists {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "Користувач з таким email вже існує", http.StatusConflict)
		return
	}

	if err := h.service.CreateUser(r.Context(), &regRequest); err != nil {
		h.log.Error("Ошибка при создании пользователя: ", slog.Any("Error:", err)) 
		http.Error(w, "Помилка на сервері", http.StatusInternalServerError)
		return
	}

	token, err := auth.CreateToken(regRequest.Login, regRequest.UserID, h.tokenTTL)
	if err != nil {
		h.log.Error("Ошибка создании токена: ", slog.Any("Error:", err)) 
		http.Error(w, "Помилка на сервері", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *UsersHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq domain.LoginRequest

	//log.Println("Запрос логина: ", r.Header, r.Body)

	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		h.log.Info("Ошибка при логине: ", slog.Any("Error:", err))
		http.Error(w, "Неправильний запит", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetByUsername(r.Context(), loginReq.Login)
	if err != nil {
		h.log.Info("Пользователь не найден: ", slog.Any("Error:", err))
		http.Error(w, "Користувач не знайден", http.StatusUnauthorized)
		return
	}

	if err := auth.CheckPassword(user.HashPassword, loginReq.Password); err != nil {
		h.log.Info("Неверный пароль", slog.Any("Error:", err))
		http.Error(w, "Неправильний пароль", http.StatusUnauthorized)
		return
	}

	token, err := auth.CreateToken(user.Login, user.UserID, h.tokenTTL)
	if err != nil {
		h.log.Info("Ошибка при создании токена: ", slog.Any("Error:", err))
		http.Error(w, "Помилка на сервері", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *UsersHandler) UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		h.log.Info("Ошибка при проверке токена: ", slog.Any("Auth Header:", authHeader))
		http.Error(w, "Помилка на сервері", http.StatusInternalServerError)
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenStr == "" {
		h.log.Info("Ошибка при проверке Bearer: ", slog.Any("Token:", tokenStr))
		http.Error(w, "Помилка на сервері", http.StatusInternalServerError)
		return
	}

	claims, err := auth.ParseToken(tokenStr)
	if err != nil {
		h.log.Info("Ошибка при парсинге токена: ", slog.Any("Error:", err))
		http.Error(w, "Не авторизовано", http.StatusUnauthorized)
		return
	}

	user, err := h.service.GetByUsername(r.Context(), claims.Username)
	if err != nil {
		h.log.Error("Ошибка при получении пользователя:", slog.Any("Error:", err))
		http.Error(w, "Помилка на сервері", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UsersHandler) UpdateUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		h.log.Error("Помилка парсингу форми: ", slog.Any("Error:", err))
		http.Error(w, "Помилка парсингу форми", http.StatusBadRequest)
		return
	}

	userData := domain.UserUpdateRequest{
		UserID:     userID,
		Login:      r.FormValue("login"),
		Password:   r.FormValue("password"),
		FirstName:  r.FormValue("first_name"),
		LastName:   r.FormValue("last_name"),
		Email:      r.FormValue("email"),
		Phone:      r.FormValue("phonenumber"),
		Address:    r.FormValue("address"),
	}

	if userData.Login == "" || userData.Password == "" || userData.FirstName == "" ||
	userData.LastName == "" || userData.Email == "" || userData.Phone == "" || userData.Address == "" {
		http.Error(w, "Усі поля повинні бути заповнені", http.StatusBadRequest)
		return
	}

	// Обробка аватара (необов'язкове)
	file, header, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()
		avatarPath, saveErr := h.service.SaveAvatarFile(userID, file, header)
		if saveErr != nil {
			h.log.Error("Помилка збереження аватара: ", slog.Any("Error:", saveErr))
			http.Error(w, "Помилка на сервері", http.StatusInternalServerError)
			return
		}
		userData.AvatarPath = &avatarPath
	}

	if err := h.service.UpdateUserProfile(r.Context(), userData); err != nil {
		h.log.Error("Помилка оновлення профілю: ", slog.Any("Error:", err))
		http.Error(w, "Помилка на сервері", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Профіль оновлено успішно",
	})
}
