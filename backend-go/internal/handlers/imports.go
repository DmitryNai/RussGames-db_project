package handlers

import (
    "context"
    "encoding/csv"
    "io"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

// BatchImportUsers handles CSV upload and COPY into users table
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
        header, err := r.Read()
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid csv header"})
            return
        }
        _ = header

        ctx := context.Background()
        conn, err := pgxPool.Acquire(ctx)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        defer conn.Release()

        // We will stream lines to COPY
        pgConn := conn.Conn()
        copyCmd := "COPY users(id,username,email,password_hash,created_at,status,country,profile) FROM STDIN WITH (FORMAT csv)"
        w, err := pgConn.PgConn().CopyFrom(ctx, copyCmd)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        count := 0
        for {
            rec, err := r.Read()
            if err == io.EOF {
                break
            }
            if err != nil {
                // optionally log to batch_errors
                continue
            }
            // build CSV line (assumes commas are safe — for production escape properly)
            line := "\"" + rec[0] + "\",\"" + rec[1] + "\",\"" + rec[2] + "\",\"" + rec[3] + "\",\"" + rec[4] + "\",\"" + rec[5] + "\",\"" + rec[6] + "\",\"" + rec[7] + "\"
"
            if _, err := w.Write([]byte(line)); err != nil {
                w.Close()
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
            count++
        }

        if err := w.Close(); err != nil && err != pgx.ErrNoRows {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"imported": count})
    }
}