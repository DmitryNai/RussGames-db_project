package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/DmitryNai/RussGames-db_project/backend-go/internal/models"
	"github.com/google/uuid"
)

func CreatePurchase(db *gorm.DB) gin.HandlerFunc {
	type req struct {
		UserID    uuid.UUID  `json:"user_id" binding:"required"`
		GameID    uuid.UUID  `json:"game_id" binding:"required"`
		PricePaid float64    `json:"price_paid" binding:"required"`
		Method    string     `json:"method" binding:"required"`
		LicenseID *uuid.UUID `json:"license_id"`
	}

	return func(c *gin.Context) {
		var r req
		if err := c.ShouldBindJSON(&r); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec("SET LOCAL app.current_user = ?", r.UserID.String()).Error; err != nil {
				return err
			}

			tr := models.Transaction{
				ID:        uuid.New(),
				UserID:    r.UserID,
				Amount:    r.PricePaid,
				Currency:  "RUB",
				CreatedAt: time.Now(),
				Provider:  nil,
				Status:    "completed",
				Details:   nil,
			}
			if err := tx.Create(&tr).Error; err != nil {
				return err
			}

			p := models.Purchase{
				ID:            uuid.New(),
				UserID:        r.UserID,
				GameID:        r.GameID,
				TransactionID: &tr.ID,
				PricePaid:     r.PricePaid,
				PurchasedAt:   time.Now(),
				Method:        r.Method,
				LicenseID:     r.LicenseID,
			}
			if err := tx.Create(&p).Error; err != nil {
				return err
			}

			if r.LicenseID != nil {
				if err := tx.Model(&models.GameLicense{}).Where("id = ?", r.LicenseID).Updates(map[string]interface{}{
					"assigned_to_user": r.UserID,
					"assigned_at":      time.Now(),
					"state":            "assigned",
				}).Error; err != nil {
					return err
				}
			}

			lib := models.Library{
				ID:        uuid.New(),
				UserID:    r.UserID,
				GameID:    r.GameID,
				LicenseID: r.LicenseID,
				AddedAt:   time.Now(),
				Active:    true,
			}
			if err := tx.Create(&lib).Error; err != nil {
				return err
			}

			return nil
		}, &sql.TxOptions{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"ok": true})
	}
}
