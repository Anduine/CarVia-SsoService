package http_handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sso-service/internal/domain"
	"sso-service/internal/lib/responseHTTP"
	"sso-service/internal/service"
	"sso-service/pkg/auth"
	"time"
)

type UsersHandler struct {
	service  *service.UsersService
	tokenTTL time.Duration
}

func NewUsersHandler(service *service.UsersService, tokenTTL time.Duration) *UsersHandler {
	return &UsersHandler{
		service:  service,
		tokenTTL: tokenTTL,
	}
}

func (h *UsersHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var regRequest domain.RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&regRequest)
	if err != nil {
		responseHTTP.JSONError(w, http.StatusBadRequest, "Неправильний запит")
		return
	}

	exists, err := h.service.ExistsByEmail(r.Context(), regRequest.Email)
	if err != nil {
		slog.Debug("Помилка при перевірці email", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusInternalServerError, "Помилка на сервері")
		return
	}
	if exists {
		responseHTTP.JSONError(w, http.StatusConflict, "Користувач з таким email вже існує")
		return
	}

	if err := h.service.CreateUser(r.Context(), &regRequest); err != nil {
		slog.Debug("Помилка при створені користувача", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusInternalServerError, "Помилка на сервері")
		return
	}

	token, err := auth.CreateToken(regRequest.Login, regRequest.UserID, h.tokenTTL)
	if err != nil {
		slog.Debug("Помилка при створені токена", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusInternalServerError, "Помилка на сервері")
		return
	}

	response := domain.TokenResponse{
		Token: token,
	}

	responseHTTP.JSONResp(w, http.StatusOK, response)
}

func (h *UsersHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq domain.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		slog.Debug("Помилка при декодуванні loginReq", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusBadRequest, "Неправильний запит")
		return
	}

	user, err := h.service.GetByEmail(r.Context(), loginReq.Email)
	if err != nil {
		slog.Debug("Користувача не знайдено", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusUnauthorized, "Користувача не знайдено")
		return
	}

	if err := auth.CheckPassword(user.HashPassword, loginReq.Password); err != nil {
		slog.Debug("Неправильний пароль", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusUnauthorized, "Неправильний пароль")
		return
	}

	token, err := auth.CreateToken(user.Login, user.UserID, h.tokenTTL)
	if err != nil {
		slog.Debug("Помилка при створені токена", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusInternalServerError, "Помилка на сервері")
		return
	}

	response := domain.TokenResponse{
		Token: token,
	}

	responseHTTP.JSONResp(w, http.StatusOK, response)
}

func (h *UsersHandler) UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value("username").(string)
	if !ok {
		slog.Debug("Помилка при отриманні username з context")
		responseHTTP.JSONError(w, http.StatusUnauthorized, "Не авторизовано")
		return
	}

	slog.Debug("Запит на отримання профілю")

	user, err := h.service.GetByUsername(r.Context(), username)
	if err != nil {
		slog.Debug("Помилка при отриманні користувача", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusInternalServerError, "Помилка на сервері")
		return
	}

	responseHTTP.JSONResp(w, http.StatusOK, user)
}

func (h *UsersHandler) UpdateUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		slog.Debug("Помилка при отриманні user_id з context")
		responseHTTP.JSONError(w, http.StatusUnauthorized, "Не авторизовано")
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		slog.Debug("Помилка парсингу форми", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusBadRequest, "Помилка парсингу форми")
		return
	}

	userData := domain.UserUpdateRequest{
		UserID:      userID,
		Login:       r.FormValue("Login"),
		Password:    r.FormValue("Password"),
		FirstName:   r.FormValue("FirstName"),
		LastName:    r.FormValue("LastName"),
		Email:       r.FormValue("Email"),
		Phonenumber: r.FormValue("Phonenumber"),
		Address:     r.FormValue("Address"),
	}

	if userData.Login == "" || userData.Password == "" || userData.FirstName == "" ||
		userData.LastName == "" || userData.Email == "" || userData.Phonenumber == "" || userData.Address == "" {
		slog.Debug("Не всі поля заповнені", "userData", userData)
		responseHTTP.JSONError(w, http.StatusBadRequest, "Усі поля повинні бути заповнені")
		return
	}

	file, header, err := r.FormFile("Avatar")
	if err == nil {
		defer file.Close()
		avatarPath, saveErr := h.service.SaveAvatar(header)
		if saveErr != nil {
			slog.Debug("Помилка збереження аватара", "err", saveErr.Error())
			responseHTTP.JSONError(w, http.StatusInternalServerError, "Помилка на сервері")
			return
		}
		userData.AvatarPath = avatarPath
	}

	err = h.service.UpdateUserProfile(r.Context(), userData)
	if err != nil {
		slog.Debug("Помилка оновлення профілю", "err", err.Error())
		responseHTTP.JSONError(w, http.StatusInternalServerError, "Помилка на сервері")
		return
	}

	responseHTTP.JSONResp(w, http.StatusOK, "Профіль оновлено")
}
