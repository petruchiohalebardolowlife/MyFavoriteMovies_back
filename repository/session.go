package repository

import (
	"myfavouritemovies/database"
	"myfavouritemovies/models"
	"time"
)

func AddSession(session *models.Session) (*models.Session, error) {
  var currentSessions []models.Session
  if err := database.DB.Where("user_id = ?", session.UserID).Order("created_at ASC").Find(&currentSessions).Error; err != nil {
    return nil, err
  }
  if len(currentSessions)>=3 {
    if err := database.DB.Delete(&currentSessions[0]).Error; err != nil {
      return nil, err
    }
  }
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

func DeleteSession(ID uint) error {
  if err := database.DB.Where("id = ?", ID).Delete(&models.Session{}).Error; err != nil {
      return err
  }

  return nil
}

func CleanExpiredSessions (duration time.Duration) {
  ticker := time.NewTicker(duration)
  go func() {
    for range ticker.C {
      database.DB.Where("expires_at < ?", time.Now()).Delete(&models.Session{})
    }
  }()
} 
