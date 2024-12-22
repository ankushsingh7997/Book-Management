package models

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
	Phone    string `json:"phone"`
}
type Claims struct {
	UserId uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

func init() {
	db.AutoMigrate(&User{})
}

func (u *User) HashPassword() error {
	if len(u.Password) == 0 {
		return errors.New("password connot be empty")
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}
	u.Password = string(hashPassword)
	return nil
}

func (u *User) CreateUser() (*User, error) {

	if err := u.HashPassword(); err != nil {
		return nil, err
	}
	if err := db.Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func GetUserByID(id uint) (*User, error) {
	var user User
	if err := db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) UpdateUser() error {
	if u.ID == 0 {
		return errors.New("user ID is required")
	}
	existingUser, err := GetUserByID(u.ID)
	if err != nil {
		return err
	}
	if u.Password != "" && u.Password != existingUser.Password {
		if err := u.HashPassword(); err != nil {
			return err
		} else {
			u.Password = existingUser.Password
		}

	}
	return db.Save(u).Error
}

func DeleteUser(id uint) error {
	return db.Delete(&User{}, id).Error
}

// GetAllUsers retrieves all users from the database
func GetAllUsers() ([]User, error) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) Login(password string) (string, error) {
	if err := u.ComparePassword(password); err != nil {
		return "", errors.New("invallid credentials")
	}

	claims := &Claims{
		UserId: u.ID,
		Email:  u.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "THis-Is-mY-SecreT"
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "THis-Is-mY-SecreT"
	}
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
