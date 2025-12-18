package handlers

import (
	"net/http"

	"github.com/DmitryNai/RussGames-db_project/backend-go/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListAuditLogs(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var logs []models.AuditLog
		if err := db.Limit(100).Find(&logs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, logs)
	}
}

func GetAuditLog(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var log models.AuditLog
		if err := db.First(&log, "id = ?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, log)
	}
}


func CreateAuditLog(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.AuditLog
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, input)
	}
}
