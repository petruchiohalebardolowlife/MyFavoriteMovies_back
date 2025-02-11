package repository

import (
	"errors"
	"myfavouritemovies/database"
	"myfavouritemovies/structs"
)

func AddUser(user *structs.User) error {
  if user.NickName == "" || user.UserName == "" || user.Password == "" {
    return errors.New("Some of filed are empty")
  }
  if err:=database.DB.Where("user_name = ?", user.UserName).First(&structs.User{}).Error; err == nil {
    return errors.New("User with this username already exist")
  }

  if err := database.DB.Create(user).Error; err != nil {
    return err
  }

  return nil
}

func UpdateUser(user *structs.User, nickName *string, password *string) error {
  if nickName == nil && password == nil {
      return errors.New("all fields are empty")
  }

  if nickName != nil {
    if *nickName == "" {
      return errors.New("nickname cannot be empty")
    }
    user.NickName = *nickName
  }

  if password != nil {
    if *password == "" {
        return errors.New("password cannot be empty")
    }
    user.Password = *password
  }

  if err := database.DB.Save(user).Error; err != nil {
    return err
  }

  return nil
}

func DeleteUser(userID uint) error {
  if err := database.DB.Where("id = ?", userID).Delete(&structs.User{}).Error; err != nil {
      return errors.New("User not found")
  }

  return nil
}