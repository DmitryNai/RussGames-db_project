package main

import (
    "encoding/csv"
    "flag"
    "fmt"
    "math/rand"
    "os"
    "time"

    "github.com/brianvoe/gofakeit/v6"
    "github.com/google/uuid"
)

func main() {
    var users = flag.Int("users", 1000, "number of users to generate")
    var games = flag.Int("games", 600, "number of games to generate")
    var purchases = flag.Int("purchases", 5000, "number of purchases to generate")
    flag.Parse()

    rand.Seed(time.Now().UnixNano())
    gofakeit.Seed(time.Now().UnixNano())

    outdir := "./data"
    if err := os.MkdirAll(outdir, 0o755); err != nil {
        panic(err)
    }

    fmt.Printf("Generating users=%d games=%d purchases=%d into %s\n", *users, *games, *purchases, outdir)

    // Users
    f, err := os.Create(outdir + "/users.csv")
    if err != nil {
        panic(err)
    }
    w := csv.NewWriter(f)
    w.Write([]string{"id", "username", "email", "password_hash", "created_at", "status", "country", "profile"})
    for i := 0; i < *users; i++ {
        id := uuid.New().String()
        username := fmt.Sprintf("user%04d", i+1)
        email := username + "@example.com"
        password := uuid.New().String()
        created := gofakeit.DateRange(time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC), time.Now()).Format("2006-01-02 15:04:05")
        status := "active"
        country := gofakeit.CountryAbr()
        w.Write([]string{id, username, email, password, created, status, country, "{}"})
    }
    w.Flush()
    f.Close()

    // Games (simplified)
    fg, err := os.Create(outdir + "/games.csv")
    if err != nil {
        panic(err)
    }
    wg := csv.NewWriter(fg)
    wg.Write([]string{"id", "developer_id", "title", "description", "genre", "price", "release_date", "avg_rating", "sales_count", "created_at", "metadata"})
    for i := 0; i < *games; i++ {
        id := uuid.New().String()
        dev := uuid.New().String()
        title := gofakeit.HipsterSentence(3)
        desc := gofakeit.Sentence(10)
        genre := gofakeit.RandomString([]string{"Action", "RPG", "Strategy", "Indie"})
        price := fmt.Sprintf("%.2f", gofakeit.Price(0, 59.99))
        rdate := gofakeit.DateRange(time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC), time.Now()).Format("2006-01-02")
        avg := fmt.Sprintf("%.2f", gofakeit.Float32Range(0, 10))
        sales := fmt.Sprintf("%d", gofakeit.Number(0, 2000))
        created := gofakeit.DateRange(time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC), time.Now()).Format("2006-01-02 15:04:05")
        wg.Write([]string{id, dev, title, desc, genre, price, rdate, avg, sales, created, "{}"})
    }
    wg.Flush()
    fg.Close()

    // Purchases (simplified)
    fp, err := os.Create(outdir + "/purchases.csv")
    if err != nil {
        panic(err)
    }
    wp := csv.NewWriter(fp)
    wp.Write([]string{"id", "user_id", "game_id", "transaction_id", "price_paid", "purchased_at", "method", "license_id"})
    for i := 0; i < *purchases; i++ {
        id := uuid.New().String()
        user := uuid.New().String()
        game := uuid.New().String()
        tr := uuid.New().String()
        price := fmt.Sprintf("%.2f", gofakeit.Price(0, 59.99))
        purchased := gofakeit.DateRange(time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC), time.Now()).Format("2006-01-02 15:04:05")
        method := gofakeit.RandomString([]string{"card", "wallet", "promo"})
        wp.Write([]string{id, user, game, tr, price, purchased, method, ""})
    }
    wp.Flush()
    fp.Close()

    fmt.Println("done")
}
