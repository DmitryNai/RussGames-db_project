package main

import (
    "log"
    "os"

    "github.com/DmitryNai/RussGames-db_project/backend-go/internal/api"
    "github.com/DmitryNai/RussGames-db_project/backend-go/internal/config"
    "github.com/DmitryNai/RussGames-db_project/backend-go/internal/db"

    "github.com/gin-gonic/gin"
)

func main() {
    cfg := config.Load()

    gormDB, pgxPool, err := db.Init(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("db init: %v", err)
    }
    defer pgxPool.Close()

    r := gin.Default()
    api.RegisterRoutes(r, gormDB, pgxPool)

    addr := os.Getenv("HTTP_ADDR")
    if addr == "" {
        addr = ":8000"
    }
    log.Printf("listening on %s", addr)
    if err := r.Run(addr); err != nil {
        log.Fatalf("server failed: %v", err)
    }
}