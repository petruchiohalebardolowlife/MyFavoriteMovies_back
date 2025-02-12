package repository

import (
	"errors"
	"myfavouritemovies/database"
	"myfavouritemovies/models"
)

func AddUser(user *models.User) error {
  if len(user.NickName) == 0 || len(user.UserName) == 0 || len(user.Password) == 0 {
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

func UpdateUser(user *models.User, nickName *string, password *string) error {
  if len(*nickName) == 0 && len(*password) == 0 {
    return errors.New("at least one field must be provided")
  }

  if len(*nickName) > 0 {
    user.NickName = *nickName
  }
  if len(*password) > 0 {
    user.Password = *password
  }

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