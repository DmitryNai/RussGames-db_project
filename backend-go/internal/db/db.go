package db

import (
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func Init(databaseURL string) (*gorm.DB, *pgxpool.Pool, error) {
    gormDB, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Warn),
    })
    if err != nil {
        return nil, nil, err
    }

    cfg, err := pgxpool.ParseConfig(databaseURL)
    if err != nil {
        return nil, nil, err
    }
    cfg.MaxConns = 10
    cfg.MinConns = 1
    cfg.MaxConnLifetime = time.Hour

    pool, err := pgxpool.NewWithConfig(nil, cfg)
    if err != nil {
        return nil, nil, err
    }

    return gormDB, pool, nil
}