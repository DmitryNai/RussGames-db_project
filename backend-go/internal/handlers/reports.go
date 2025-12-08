package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func ReportTopGames(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        type row struct {
            GameID    string  `json:"game_id"`
            Title     string  `json:"title"`
            Sales     int64   `json:"sales"`
            Revenue   float64 `json:"revenue"`
        }
        var rows []row
        if err := db.Raw(`SELECT g.id as game_id, g.title, g.sales_count as sales, SUM(p.price_paid) as revenue FROM games g LEFT JOIN purchases p ON p.game_id = g.id GROUP BY g.id ORDER BY sales_count DESC LIMIT 20`).Scan(&rows).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, rows)
    }
}

func ReportSalesByUser(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        type row struct {
            UserID      string  `json:"user_id"`
            Username    string  `json:"username"`
            TotalPurchases int   `json:"total_purchases"`
            TotalSpent   float64 `json:"total_spent"`
        }
        var rows []row
        if err := db.Raw(`SELECT u.id as user_id, u.username, COUNT(p.id) as total_purchases, COALESCE(SUM(p.price_paid),0) as total_spent FROM users u LEFT JOIN purchases p ON p.user_id = u.id GROUP BY u.id, u.username ORDER BY total_spent DESC LIMIT 50`).Scan(&rows).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, rows)
    }
}