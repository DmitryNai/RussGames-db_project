package handlers

import (
	"context"
	"encoding/csv"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgconn"
)

func BatchImportUsers(pgxPool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, _, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
			return
		}
		defer file.Close()

		r := csv.NewReader(file)
		// Read header
		_, err = r.Read()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid csv header"})
			return
		}

		ctx := context.Background()
		conn, err := pgxPool.Acquire(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer conn.Release()

		copyCmd := "COPY users(id,username,email,password_hash,created_at,status,country,profile) FROM STDIN WITH (FORMAT csv)"

		pr, pw := io.Pipe()
		go func() {
			defer pw.Close()
			for {
				rec, err := r.Read()
				if err == io.EOF {
					return
				}
				if err != nil {
					continue
				}
				for i := range rec {
					rec[i] = `"` + escapeQuotes(rec[i]) + `"`
				}
				line := ""
				for i, f := range rec {
					if i > 0 {
						line += ","
					}
					line += f
				}
				line += "\n"
				if _, err := pw.Write([]byte(line)); err != nil {
					return
				}
			}
		}()

		pgConn := conn.Conn().PgConn()
		_, err = pgConn.CopyFrom(ctx, pr, copyCmd)
		if err != nil {
			pr.Close()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"imported": "ok"})
	}
}

func escapeQuotes(s string) string {
	if s == "" {
		return ""
	}
	b := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '"' {
			b = append(b, '"', '"')
		} else {
			b = append(b, s[i])
		}
	}
	return string(b)
}
