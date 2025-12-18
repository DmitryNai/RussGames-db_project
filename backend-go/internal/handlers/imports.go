package handlers

import (
	"context"
	"encoding/csv"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

var identRe = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

func safeIdent(s string) (string, bool) {
	if !identRe.MatchString(s) {
		return "", false
	}
	return `"` + s + `"`, true
}

func BatchImportHandler(pgxPool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		tableRaw := c.Query("table")
		schemaRaw := c.DefaultQuery("schema", "public")

		table, ok := safeIdent(tableRaw)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table name"})
			return
		}

		schema, ok := safeIdent(schemaRaw)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schema name"})
			return
		}

		file, _, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
			return
		}
		defer file.Close()

		r := csv.NewReader(file)

		// header
		header, err := r.Read()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid csv header"})
			return
		}

		cols := make([]string, 0, len(header))
		for _, h := range header {
			h = strings.TrimSpace(h)
			col, ok := safeIdent(h)
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "invalid column name: " + h,
				})
				return
			}
			cols = append(cols, col)
		}

		copyCmd := "COPY " + schema + "." + table +
			" (" + strings.Join(cols, ",") + ")" +
			" FROM STDIN WITH (FORMAT csv, NULL '\\N')"

		ctx := context.Background()
		conn, err := pgxPool.Acquire(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer conn.Release()

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
					if strings.TrimSpace(rec[i]) == "" {
						rec[i] = `\N`
					} else {
						rec[i] = escapeCSV(rec[i])
					}
				}

				line := strings.Join(rec, ",") + "\n"
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

func escapeCSV(s string) string {
	s = strings.ReplaceAll(s, `"`, `""`)
	return `"` + s + `"`
}
