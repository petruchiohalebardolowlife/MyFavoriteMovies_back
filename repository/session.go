package repository

import (
	"myfavouritemovies/database"
	"myfavouritemovies/models"
)

func AddSession(session *models.Session) (*models.Session, error) {
	if err := database.DB.Create(&session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func GetSession (ID string) (*models.Session, error) {
  var session models.Session
  if err := database.DB.Where("id = ?", ID).First(&session).Error; err != nil {
    return nil, err
  }

  return &session, nil
}

func RevokeSession (ID string) error {
  if err := database.DB.Model(&models.Session{}).Where("id = ?", ID).Update("revoke", true).Error; err != nil {
    return err
  }

  return nil
}

func DeleteSession(ID uint) error {
  if err := database.DB.Where("id = ?", ID).Delete(&models.Session{}).Error; err != nil {
      return err
  }

  return nil
}
