package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "github.com/DmitryNai/RussGames-db_project/backend-go/internal/models"
    "github.com/google/uuid"
)

func ListGames(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var games []models.Game
        if err := db.Limit(100).Find(&games).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, games)
    }
}

func CreateGame(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input models.Game
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        if input.ID == uuid.Nil {
            input.ID = uuid.New()
        }
        if err := db.Create(&input).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusCreated, input)
    }
}

func GetGame(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")
        var game models.Game
        if err := db.First(&game, "id = ?", id).Error; err != nil {
            if err == gorm.ErrRecordNotFound {
                c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
                return
            }
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, game)
    }
}