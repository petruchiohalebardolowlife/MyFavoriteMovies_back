package repository

import (
	"errors"
	"myfavouritemovies/database"
	"myfavouritemovies/models"
	"myfavouritemovies/security"
	"strings"
)

func AddUser(user *models.User) error {
  if len(strings.ReplaceAll(user.NickName, " ", "")) == 0 || len(strings.ReplaceAll(user.UserName, " ", "")) == 0 || len(strings.ReplaceAll(user.PasswordHash, " ", "")) == 0 {
    return errors.New("some of fields are empty")
  }

  if err := database.DB.Where("user_name = ?", user.UserName).First(&models.User{}).Error; err == nil {
    return errors.New("user with this username already exists")
  }

  if err := database.DB.Create(user).Error; err != nil {
    return err
  }

  return nil
}

func UpdateNickName(user *models.User, nickname string) error {
  if len(strings.ReplaceAll(nickname, " ", "")) == 0 {
    return errors.New("nickname cannot be empty")
  }
  user.NickName = nickname
  if err := database.DB.Save(user).Error; err != nil {
    return err
  }
  return nil
}

func UpdatePassWord(user *models.User, password string) error {
  if len(strings.ReplaceAll(password, " ", "")) == 0 {
    return errors.New("password cannot be empty")
  }
  hash, err := security.GenerateHashPassword(password)
  if err != nil {
    return err
  }
  user.PasswordHash = hash
  if err := database.DB.Save(user).Error; err != nil {
    return err
  }
  return nil
}

func DeleteUser(userID uint) error {
  if err := database.DB.Where("id = ?", userID).Delete(&models.User{}).Error; err != nil {
    return errors.New("user not found")
  }

  return nil
}

func GetUserByID(userID uint) (*models.User, error) {
  var user *models.User
  if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
    return nil, err
  }

  return user, nil
}

func GetUserByUserName(userName string) (*models.User, error) {
  var user *models.User
  if err := database.DB.Where("user_name = ?", userName).First(&user).Error; err != nil {
    return nil, err
  }

  return user, nil
}
