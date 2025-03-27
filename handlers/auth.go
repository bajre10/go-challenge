package handlers

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"net/http"
	"os"
	"shortcut-challenge/api"
	"shortcut-challenge/api/appErrors"
	"shortcut-challenge/database"
	"shortcut-challenge/models"
	"shortcut-challenge/utils"
	"strconv"
	"time"
)

const DEFAULT_JWT_EXP = 3600

type UserResponseDTO struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type AuthResponseDTO struct {
	Token string `json:"token"`
	User  UserResponseDTO
}

func createToken(user *models.User) (string, error) {
	exp, err := strconv.Atoi(utils.GetEnv("JWT_EXP", string(rune(DEFAULT_JWT_EXP))))
	if err != nil {
		exp = DEFAULT_JWT_EXP
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    user.ID,
		"role":   user.Role.Name,
		"roleID": user.Role.ID,
		"exp":    time.Now().Add(time.Duration(exp) * time.Second).Unix(),
		"iat":    time.Now().Unix(),
	})

	secret := []byte(os.Getenv("JWT_SECRET"))
	token, err := claims.SignedString(secret)
	if err != nil {
		return "", err
	}
	return token, nil
}

func createSendToken(w http.ResponseWriter, user *models.User) {
	token, err := createToken(user)
	if err != nil {
		api.ThrowInternalError(w)
		return
	}

	exp, err := strconv.Atoi(utils.GetEnv("JWT_EXP", string(rune(DEFAULT_JWT_EXP))))
	if err != nil {
		exp = DEFAULT_JWT_EXP
	}

	cookieExp := time.Now().Add(time.Duration(exp) * time.Second)
	cookie := http.Cookie{Name: "jwt", Value: token, Expires: cookieExp, Path: "/"}
	http.SetCookie(w, &cookie)

	// Send success auth response
	api.SendJson(w, http.StatusOK, AuthResponseDTO{
		User: UserResponseDTO{
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role.Name,
		},
		Token: token,
	})
}

func GetClaims(r *http.Request) *jwt.MapClaims {
	// Get JWT claims instance from request context
	if db, ok := r.Context().Value("JWT_CLAIMS").(*jwt.MapClaims); ok {
		return db
	}
	return nil
}

type RegisterUserDTO struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

func registerUser(w http.ResponseWriter, r *http.Request, role *models.Role) {
	db := database.GetDBInstance(r)

	var userDTO RegisterUserDTO

	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		api.ThrowRequestError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate email and password
	if userDTO.Email == "" || userDTO.Password == "" {
		api.ThrowRequestError(w, "Email and Password are required", http.StatusBadRequest)
		return
	}

	// Check if passwords match
	if userDTO.Password != userDTO.ConfirmPassword {
		api.ThrowRequestError(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	user := models.User{
		Name:     userDTO.Name,
		Email:    userDTO.Email,
		Password: userDTO.Password,
		Role:     *role,
	}

	// Insert the user into the database
	if err := db.Create(&user).Error; err != nil {
		var dbError *appErrors.DBError
		if errors.As(err, &dbError) {
			api.ThrowRequestError(w, dbError.Message, http.StatusBadRequest)
			return
		} else if errors.Is(err, gorm.ErrDuplicatedKey) {
			api.ThrowRequestError(w, "User already exists", http.StatusConflict)
			return
		} else {
			api.ThrowRequestError(w, "Error registering user", http.StatusInternalServerError)
			return
		}
	}

	createSendToken(w, &user)
}

// Register godoc
// @Summary Register a new user
// @Description Registers a new user with the given email, password, and name. Validates the input, checks password match, assigns a default role, and creates the user in the database.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body RegisterUserDTO true "User registration details"
// @Success 201 {object} AuthResponseDTO "User successfully registered"
// @Failure 400 {object} api.RequestError "Invalid input or user already exists"
// @Failure 500 {object} api.RequestError "Internal server error"
// @Router /auth/register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	db := database.GetDBInstance(r)

	var userRole models.Role
	if err := db.Where("name = ?", string(models.USER)).First(&userRole).Error; err != nil {
		api.ThrowInternalError(w)
		return
	}

	registerUser(w, r, &userRole)
}

// RegisterAdmin godoc
// @Security BearerAuth
// @Summary Register a new admin
// @Description Registers a new admin with the given email, password, and name. Validates the input, checks password match and creates the user in the database.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param user body RegisterUserDTO true "User registration details"
// @Success 201 {object} AuthResponseDTO "User successfully registered"
// @Failure 400 {object} api.RequestError "Invalid input or user already exists"
// @Failure 500 {object} api.RequestError "Internal server error"
// @Router /auth/admin/register [post]
func RegisterAdmin(w http.ResponseWriter, r *http.Request) {
	db := database.GetDBInstance(r)

	var userRole models.Role
	if err := db.Where("name = ?", string(models.ADMIN)).First(&userRole).Error; err != nil {
		api.ThrowInternalError(w)
		return
	}

	registerUser(w, r, &userRole)
}

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login godoc
// @Summary Login user
// @Description Authenticates a user by validating the email and password, and generates an authentication token if successful.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param login body LoginDTO true "User login details"
// @Success 200 {object} AuthResponseDTO "Authentication successful, token returned"
// @Failure 400 {object} api.RequestError "Invalid input"
// @Failure 401 {object} api.RequestError "Incorrect email or password"
// @Failure 500 {object} api.RequestError "Internal server error"
// @Router /auth/login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	db := database.GetDBInstance(r)

	var loginDTO LoginDTO
	if err := json.NewDecoder(r.Body).Decode(&loginDTO); err != nil {
		api.ThrowRequestError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate email and password
	if loginDTO.Email == "" || loginDTO.Password == "" {
		api.ThrowRequestError(w, "Email and Password are required", http.StatusBadRequest)
		return
	}

	// Get user from DB
	fetchedUser := models.User{}
	if err := db.Joins("Role").Where("email = ?", loginDTO.Email).First(&fetchedUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.ThrowRequestError(w, "Incorrect email or password", http.StatusUnauthorized)
		} else {
			api.ThrowInternalError(w)
		}
		return
	}

	// Compare passwords
	if err := fetchedUser.CheckPassword(loginDTO.Password); err != nil {
		api.ThrowRequestError(w, "Incorrect email or password", http.StatusUnauthorized)
		return
	}

	createSendToken(w, &fetchedUser)
}
