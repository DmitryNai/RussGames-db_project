package config

import "os"

type Config struct {
    DatabaseURL string
}

func Load() Config {
    return Config{
        DatabaseURL: getenv("DATABASE_URL", "postgres://postgres:example@db:5432/russgames?sslmode=disable"),
    }
}

func getenv(k, def string) string {
    v := os.Getenv(k)
    if v == "" {
        return def
    }
    return v
}