package backend

import (
	"context"
	"errors"
	"net/http"
	"rip/lib/auth"
	"time"

	"github.com/gin-gonic/gin"
)

type RegisterUserRequest struct {
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" example:"johndoe@example.com"`
	Login    string `json:"login" example:"johndoe123"`
	Password string `json:"password" example:"newsecurepassword456"`
}

// RegisterUser godoc
// @Summary Registers a new user
// @Description This endpoint registers a new user by accepting a JSON payload with the user's details (name, email, login, password).
// @Tags Users
// @Accept json
// @Produce json
// @Param request body RegisterUserRequest true "User registration data"
// @Success 201 {object} gin.H "User registered successfully"
// @Failure 400 {object} ErrorResponse "Invalid request format or user registration failed"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /user/register [post]
func (app *App) RegisterUser(c *gin.Context) {
	var req RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid request format"), err)
		return
	}

	userID, err := app.createUser(req.Login, req.Password, req.Name, req.Email)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] failed to save user in the database"), err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"userID":  userID,
	})
}

type UpdateUserProfileRequest struct {
	Password string `json:"password" example:"newsecurepassword456"`
	Email    string `json:"email" example:"test@test.com"`
	Name     string `json:"name" example:"Jane Doe 1"`
}

// UpdateUserProfile godoc
// @Summary Updates user profile
// @Description This endpoint allows users to update their profile details (name, email, password).
// @Tags Users
// @Accept json
// @Produce json
// @Param request body UpdateUserProfileRequest true "User profile update data"
// @Success 200 {object} gin.H "User profile updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request format or update failed"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /user/update [put]
func (app *App) UpdateUserProfile(c *gin.Context) {
	requestUserID, err := ExtractUserID(c)
	if err != nil {
		handleError(c, http.StatusUnauthorized, errors.New("[err] Unauthorized"), err)
		return
	}

	var req UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid request format"), err)
		return
	}

	user, err := app.getUserFirst(requestUserID)
	if err != nil {
		handleError(c, http.StatusNotFound, errors.New("[err] user not found"), err)
		return
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = &req.Email
	}
	if req.Password != "" {
		user.Password = []byte(req.Password)
	}

	if err := app.updateUser(&user); err != nil {
		handleError(c, http.StatusInternalServerError, errors.New("[err] failed to update user profile"), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User profile updated successfully",
		"user":    user,
	})
}

type UserLoginRequest struct {
	Login    string `json:"login" example:"user1"`
	Password string `json:"password" example:"userPass123"`
}

type UserLoginResponse struct {
	ExpiresIn   time.Duration `json:"expires_in" example:"86400"`
	AccessToken string        `json:"access_token" example:"JWT_TOKEN"`
	TokenType   string        `json:"token_type" example:"Bearer"`
	User        DbUser        `json:"user"` // Включаем всю структуру пользователя
}

// UserLogin godoc
// @Summary User login
// @Description Authenticates the user and returns a JWT token on successful login along with user details.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body UserLoginRequest true "User login credentials"
// @Success 200 {object} UserLoginResponse "User login successful and JWT token generated along with user details"
// @Failure 400 {object} ErrorResponse "Invalid login or password"
// @Failure 500 {object} ErrorResponse "Failed to generate or save session"
// @Router /user/login [post]
func (app *App) UserLogin(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid request"), err)
		return
	}

	isMatch, user, err := app.matchPassword(req.Login, req.Password)
	if err != nil || !isMatch {
		handleError(c, http.StatusBadRequest, errors.New("[err] invalid login or password"), err)
		return
	}

	token, err := auth.GenerateJWT(user.ID, string(user.Role), app.secret)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] failed to generate token"), err)
		return
	}

	ctx := context.Background()
	expiration := time.Hour * 24
	err = SaveSession(ctx, app.redisClient, user.ID, string(user.Role), expiration)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] failed to save session"), err)
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		ExpiresIn:   expiration,
		AccessToken: token,
		TokenType:   "Bearer",
		User:        user,
	})
}

// UserLogout godoc
// @Summary User logout
// @Description Logs the user out by deleting the session and clearing the authentication token from the cookie.
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} gin.H "User logged out successfully"
// @Failure 400 {object} ErrorResponse "Failed to clear session"
// @Failure 401 {object} ErrorResponse "Unauthorized: User not authenticated"
// @Router /user/logout [post]
func (app *App) UserLogout(c *gin.Context) {
	idAny, exists := c.Get("userID")
	if !exists {
		handleError(c, http.StatusUnauthorized, errors.New("[err] Unauthorized"))
		return
	}
	requestUserID, ok := idAny.(string)
	if !ok {
		handleError(c, http.StatusBadRequest, errors.New("[err] Unauthorized"), errors.New("userID is not of type *uuid.UUID"))
		return
	}

	ctx := context.Background()
	err := DeleteSession(ctx, app.redisClient, requestUserID)
	if err != nil {
		handleError(c, http.StatusBadRequest, errors.New("[err] failed to clear session"), err)
		return
	}

	// Удаление токена из cookie
	c.SetCookie("auth_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
		"status":  true,
	})
}
