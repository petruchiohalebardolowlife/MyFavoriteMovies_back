package security

import (
	"errors"
	"myfavouritemovies/database"
	"myfavouritemovies/models"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHashPassword(password string) (string, error) {
  passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  if err != nil {
    return "", err
  }
  return string(passwordHash), nil
}

func CheckPassword(passwordHash string, password string) error {
  err:=bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
  if err!=nil {
    return err
  }
  return nil
}

func SignIn (userName string, password string) (*models.User, error) {
  var user models.User
  if err := database.DB.Where("user_name = ?", userName).First(&user).Error; err != nil {
    return nil, errors.New("incorrect username or password")
}
  if err := CheckPassword(user.PasswordHash, password); err!= nil {
    return nil, errors.New("incorrect username or password")
  }

  return &user, nil
}