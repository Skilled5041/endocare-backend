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
			Date        string  `json:"date"`
			Duration    float64 `json:"duration"`
			Quality     int32   `json:"quality"`
			Disruptions string  `json:"disruptions"`
			Notes       string  `json:"notes"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		parsedDate, err := time.Parse(time.RFC3339, req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, expected RFC3339"})
			return
		}

		params := database.InsertSleepParams{
			Date:        pgtype.Date{Time: parsedDate, Valid: true},
			Duration:    pgtype.Float8{Float64: req.Duration, Valid: true},
			Quality:     pgtype.Int4{Int32: req.Quality, Valid: true},
			Disruptions: pgtype.Text{String: req.Disruptions, Valid: true},
			Notes:       pgtype.Text{String: req.Notes, Valid: true},
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
			Date  string   `json:"date"`
			Items []string `json:"items"`
			Notes string   `json:"notes"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		parsedTime, err := time.Parse(time.RFC3339, req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid time format, expected HH:MM:SS"})
			return
		}

		params := database.InsertDietParams{
			Meal:  pgtype.Text{String: req.Meal, Valid: true},
			Date:  pgtype.Date{Time: parsedTime, Valid: true},
			Items: req.Items,
			Notes: pgtype.Text{String: dbURL, Valid: true},
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
			PeriodEvent string `json:"period_event"`
			Date        string `json:"date"`
			FlowLevel   string `json:"flow_level"`
			Notes       string `json:"notes"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		parsedDate, err := time.Parse(time.RFC3339, req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, expected RFC3339"})
			return
		}

		params := database.InsertMenstrualParams{
			PeriodEvent: pgtype.Text{String: req.PeriodEvent, Valid: true},
			Date:        pgtype.Date{Time: parsedDate, Valid: true},
			FlowLevel:   pgtype.Text{String: req.FlowLevel, Valid: true},
			Notes:       pgtype.Text{String: req.Notes, Valid: true},
		}

		res, err := queries.InsertMenstrual(ctx, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	r.POST("/insert_symptoms", func(c *gin.Context) {
		var req struct {
			Date    string `json:"date"`
			Nausea  int32  `json:"nausea"`
			Fatigue int32  `json:"fatigue"`
			Pain    int32  `json:"pain"`
			Notes   string `json:"notes"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		parsedDate, err := time.Parse(time.RFC3339, req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, expected RFC3339"})
			return
		}
		params := database.InsertSymptomsParams{
			Date:    pgtype.Date{Time: parsedDate, Valid: true},
			Nausea:  pgtype.Int4{Int32: req.Nausea, Valid: true},
			Fatigue: pgtype.Int4{Int32: req.Fatigue, Valid: true},
			Pain:    pgtype.Int4{Int32: req.Pain, Valid: true},
			Notes:   pgtype.Text{String: req.Notes, Valid: true},
		}

		res, err := queries.InsertSymptoms(ctx, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	r.GET("/get_all_sleep", func(c *gin.Context) {
		res, err := queries.GetAllSleep(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	r.GET("/get_all_diet", func(c *gin.Context) {
		res, err := queries.GetAllDiet(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	r.GET("/get_all_menstrual", func(c *gin.Context) {
		res, err := queries.GetAllMenstrual(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	r.GET("/get_all_symptoms", func(c *gin.Context) {
		res, err := queries.GetAllSymptoms(ctx)
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
