package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"

	"terrahack2025-backend/database"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("Missing required environment variable: DATABASE_URL")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer func() {
		if err := conn.Close(ctx); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	queries := database.New(conn)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.POST("/insert_sleep", func(c *gin.Context) {
		var req struct {
			Duration      float64 `json:"duration"`
			Efficiency    int32   `json:"efficiency"`
			DeepPct       int32   `json:"deep_pct"`
			Latency       int32   `json:"latency"`
			NumAwakenings int32   `json:"num_awakenings"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		params := database.InsertSleepParams{
			Duration:      pgtype.Float8{Float64: req.Duration, Valid: true},
			Efficiency:    pgtype.Int4{Int32: req.Efficiency, Valid: true},
			DeepPct:       pgtype.Int4{Int32: req.DeepPct, Valid: true},
			Latency:       pgtype.Int4{Int32: req.Latency, Valid: true},
			NumAwakenings: pgtype.Int4{Int32: req.NumAwakenings, Valid: true},
		}

		res, err := queries.InsertSleep(ctx, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	r.POST("/insert_diet", func(c *gin.Context) {
		var req struct {
			Meal  string   `json:"meal"`
			Time  string   `json:"time"`
			Items []string `json:"items"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		parsedTime, err := time.Parse(time.RFC3339, req.Time)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid time format, expected HH:MM:SS"})
			return
		}

		params := database.InsertDietParams{
			Meal:  pgtype.Text{String: req.Meal, Valid: true},
			Time:  pgtype.Timestamp{Time: parsedTime, Valid: true},
			Items: req.Items,
		}

		res, err := queries.InsertDiet(ctx, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	r.POST("/insert_menstrual", func(c *gin.Context) {
		var req struct {
			CycleDay    int32    `json:"cycle_day"`
			PainRating  int32    `json:"pain_rating"`
			StressLevel int32    `json:"stress_level"`
			Medication  []string `json:"medication"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		params := database.InsertMenstrualParams{
			CycleDay:    pgtype.Int4{Int32: req.CycleDay, Valid: true},
			PainRating:  pgtype.Int4{Int32: req.PainRating, Valid: true},
			StressLevel: pgtype.Int4{Int32: req.StressLevel, Valid: true},
			Medication:  req.Medication,
		}

		res, err := queries.InsertMenstrual(ctx, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	fmt.Printf("Server is running on http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
