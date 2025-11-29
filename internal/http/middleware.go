package http

import (
	"concert/internal/concert"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SetCookie(w http.ResponseWriter, userID uint) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_concert",
		Value:    fmt.Sprintf("user_%d", userID),
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600,
	})
}

// used when logout
func DeleteCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session_concert",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // delere the cookie
	})
}

func Authenticated(r *http.Request) bool {
	cookie, err := r.Cookie("session_concert")
	if err != nil {
		return false
	}
	return cookie != nil && cookie.Value != "" && strings.HasPrefix(cookie.Value, "user_")
}

func GetCurrentUserID(r *http.Request) (uint, error) {
	cookie, err := r.Cookie("session_concert")
	if err != nil {
		return 0, err
	}

	if !strings.HasPrefix(cookie.Value, "user_") {
		return 0, fmt.Errorf("invalid session format")
	}

	userIDString := strings.TrimPrefix(cookie.Value, "user_")
	userID, err := strconv.ParseUint(userIDString, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(userID), nil
}

func GetUserFromCookie(db *gorm.DB, r *http.Request) (*concert.User, error) {
	userID, err := GetCurrentUserID(r)
	if err != nil {
		return nil, err
	}

	var user concert.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func NeedsAuth(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !Authenticated(r) {
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func NeedsRole(db *gorm.DB, role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := GetUserFromCookie(db, r)
			if err != nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			if user.Role != role {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// checks if a user with the given email or username already exists
func UserExists(db *gorm.DB, email, username string) (bool, error) {
	var existingUser concert.User
	err := db.Where("email = ? OR username = ?", email, username).First(&existingUser).Error
	if err == nil {
		return true, nil
	}
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	return false, err
}

func FindUserByEmailOrUsername(db *gorm.DB, emailOrUsername string) (*concert.User, error) {
	var user concert.User
	err := db.Where("email = ? OR username = ?", emailOrUsername, emailOrUsername).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func InsertUser(db *gorm.DB, email, username, password, firstName, lastName, role string) (*concert.User, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	if role == "" {
		role = "user"
	}
	user := concert.User{
		Email:        email,
		Username:     username,
		PasswordHash: hashedPassword,
		FirstName:    firstName,
		LastName:     lastName,
		Role:         role,
	}

	if err := db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// Authenticate a user by email/username and password
func AuthenticateUser(db *gorm.DB, emailOrUsername, password string) (*concert.User, error) {
	user, err := FindUserByEmailOrUsername(db, emailOrUsername)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid email/username or password")
		} else {
			return nil, fmt.Errorf("error finding user: %w", err)
		}
	}

	// check if password is correct
	if !VerifyPassword(password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid email/username or password")
	}

	return user, nil
}
