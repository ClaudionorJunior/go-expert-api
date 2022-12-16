package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ClaudionorJunior/go-expert-api/internal/dto"
	"github.com/ClaudionorJunior/go-expert-api/internal/entity"
	"github.com/ClaudionorJunior/go-expert-api/internal/infra/database"
	"github.com/go-chi/jwtauth"
)

type UserHandler struct {
	UserDB database.UserInterface
}

func NewUserHandler(db database.UserInterface) *UserHandler {
	return &UserHandler{
		UserDB: db,
	}
}

// GetJWT godoc
// @Summary Get a token JWT
// @Description Get a token JWT
// @Tags users
// @Accept  json
// @Produce  json
// @Param request body dto.GetJWTInput true "user credencials"
// @Success 200 {object} dto.GetJWTOutput
// @Failure 400 {object} dto.ErrorMessageResponse
// @Failure 401 {object} dto.ErrorMessageResponse
// @Router /users/generate_token [post]
func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("jwtExpiresIn").(int)
	var authForm dto.GetJWTInput
	err := json.NewDecoder(r.Body).Decode(&authForm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorMessageResponse{
			Message: "Invalid data",
		})
		return
	}

	u, err := h.UserDB.FindByEmail(authForm.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(dto.ErrorMessageResponse{
			Message: "Invalid user credencials",
		})
		return
	}

	if !u.ValidatePassword(authForm.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(dto.ErrorMessageResponse{
			Message: "Invalid user credencials",
		})
		return
	}

	_, tokenString, _ := jwt.Encode(map[string]interface{}{
		"sub": u.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	})

	accessToken := dto.GetJWTOutput{AccessToken: tokenString}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

// Create user godoc
// @Summary Create user
// @Description Create user
// @Tags users
// @Accept  json
// @Produce  json
// @Param request body dto.CreateUserInput true "user request"
// @Success 201
// @Failure 400 {object} dto.ErrorMessageResponse
// @Failure 500 {object} dto.ErrorMessageResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorMessageResponse{
			Message: "Invalid user data",
		})
		return
	}

	err = h.UserDB.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorMessageResponse{
			Message: "Error creating user",
		})
		return
	}
	w.WriteHeader(http.StatusCreated)
}
