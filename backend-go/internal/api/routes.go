package api

import (
    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v5/pgxpool"
    "gorm.io/gorm"

    "github.com/DmitryNai/RussGames-db_project/backend-go/internal/handlers"
)

func RegisterRoutes(r *gin.Engine, gormDB *gorm.DB, pgxPool *pgxpool.Pool) {
    v1 := r.Group("/api/v1")
    {
        v1.GET("/users", handlers.ListUsers(gormDB))
        v1.POST("/users", handlers.CreateUser(gormDB))
        v1.GET("/users/:id", handlers.GetUser(gormDB))
        v1.PUT("/users/:id", handlers.UpdateUser(gormDB))
        v1.DELETE("/users/:id", handlers.DeleteUser(gormDB))

        v1.GET("/games", handlers.ListGames(gormDB))
        v1.POST("/games", handlers.CreateGame(gormDB))
        v1.GET("/games/:id", handlers.GetGame(gormDB))

        v1.POST("/purchases", handlers.CreatePurchase(gormDB))

        v1.POST("/batch-import/users", handlers.BatchImportUsers(pgxPool))

        v1.GET("/reports/top-games", handlers.ReportTopGames(gormDB))
        v1.GET("/reports/sales-by-user", handlers.ReportSalesByUser(gormDB))
    }
}