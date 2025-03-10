package repository

import (
	"errors"
	"myfavouritemovies/database"
	"myfavouritemovies/models"
	"time"

	"gorm.io/gorm"
)

func AddSession(session *models.Session) (*models.Session, error) {
  var currentSessions []models.Session
  if err := database.DB.Where("user_id = ?", session.UserID).Order("created_at ASC").Find(&currentSessions).Error; err != nil {
    return nil, err
  }
  if len(currentSessions) >= 3 {
    if err := database.DB.Delete(&currentSessions[0]).Error; err != nil {
      return nil, err
    }
  }
  if err := database.DB.Create(&session).Error; err != nil {
    return nil, err
  }

  return session, nil
}

func DeleteSession(ID string) error {
  if err := database.DB.Where("id = ?", ID).Delete(&models.Session{}).Error; err != nil {
    return err
  }

  return nil
}

func CleanExpiredTokens(duration time.Duration) {
  ticker := time.NewTicker(duration)
  go func() {
    for range ticker.C {
      database.DB.Unscoped().Where("expires_at < ?", time.Now()).Delete(&models.Session{})
      database.DB.Unscoped().Where("expires_at < ?", time.Now()).Delete(&models.BlackListToken{})
    }
  }()
}

func CheckTokenInBlackList(tokenUUID string) error {
  err := database.DB.Where("id = ?", tokenUUID).First(&models.BlackListToken{}).Error
  if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
      return nil
    }
  }

  return errors.New("your refresh token in blacklist")
}

func AddTokenInBlackList(claims *models.TokenClaims) error {
  if err := database.DB.Create(&models.BlackListToken{ID: claims.ID, UserID: claims.UserID, ExpiresAt: claims.ExpiresAt.Time}).Error; err != nil {
    return err
  }

  return nil
}
