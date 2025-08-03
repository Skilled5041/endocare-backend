package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/genai"

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

	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		log.Fatal("Missing required environment variable: GEMINI_API_KEY")
	}

	ctx2 := context.Background()
	client, err := genai.NewClient(ctx2, &genai.ClientConfig{
		APIKey: geminiAPIKey,
	})

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Use pgxpool instead of pgx.Connect
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database pool: %v", err)
	}
	defer pool.Close()

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

		queries := database.New(pool)
		res, err := queries.InsertSleep(c.Request.Context(), params)
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, expected RFC3339"})
			return
		}

		params := database.InsertDietParams{
			Meal:  pgtype.Text{String: req.Meal, Valid: true},
			Date:  pgtype.Date{Time: parsedTime, Valid: true},
			Items: req.Items,
			Notes: pgtype.Text{String: req.Notes, Valid: true},
		}

		queries := database.New(pool)
		res, err := queries.InsertDiet(c.Request.Context(), params)
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

		queries := database.New(pool)
		res, err := queries.InsertMenstrual(c.Request.Context(), params)
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

		queries := database.New(pool)
		res, err := queries.InsertSymptoms(c.Request.Context(), params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	r.GET("/get_all_sleep", func(c *gin.Context) {
		queries := database.New(pool)
		res, err := queries.GetAllSleep(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	r.GET("/get_all_diet", func(c *gin.Context) {
		queries := database.New(pool)
		res, err := queries.GetAllDiet(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	r.GET("/get_all_menstrual", func(c *gin.Context) {
		queries := database.New(pool)
		res, err := queries.GetAllMenstrual(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	r.GET("/get_all_symptoms", func(c *gin.Context) {
		queries := database.New(pool)
		res, err := queries.GetAllSymptoms(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	r.GET("/find_triggers", func(c *gin.Context) {
		queries := database.New(pool)

		sleepData, err := queries.GetAllSleep(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		dietData, err := queries.GetAllDiet(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		menstrualData, err := queries.GetAllMenstrual(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		symptomsData, err := queries.GetAllSymptoms(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		type triggerCounts struct {
			LowSleepHours  int
			MenstrualEvent map[string]int
			FlowLevel      map[string]int
			FoodItems      map[string]int
		}

		type TriggerDetail struct {
			Date            string  `json:"date"`
			TriggerSeverity float64 `json:"trigger_severity"`
		}

		triggers := triggerCounts{
			MenstrualEvent: make(map[string]int),
			FlowLevel:      make(map[string]int),
			FoodItems:      make(map[string]int),
		}

		// Track details per trigger for output
		var lowSleepDetails []TriggerDetail
		foodItemDetails := map[string][]TriggerDetail{}
		menstrualEventDetails := map[string][]TriggerDetail{}
		flowLevelDetails := map[string][]TriggerDetail{}

		// Map data by date
		sleepMap := map[string]database.Sleep{}
		for _, s := range sleepData {
			sleepMap[s.Date.Time.Format("2006-01-02")] = s
		}

		dietMap := map[string][]database.Diet{}
		for _, d := range dietData {
			date := d.Date.Time.Format("2006-01-02")
			dietMap[date] = append(dietMap[date], d)
		}

		menstrualMap := map[string]database.Menstrual{}
		for _, m := range menstrualData {
			menstrualMap[m.Date.Time.Format("2006-01-02")] = m
		}

		// Calculate mean and std dev of symptom severity
		var scores []float64
		for _, sym := range symptomsData {
			avg := float64(sym.Nausea.Int32+sym.Fatigue.Int32+sym.Pain.Int32) / 3.0
			scores = append(scores, avg)
		}
		if len(scores) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No symptom data found."})
			return
		}

		var sum float64
		for _, s := range scores {
			sum += s
		}
		mean := sum / float64(len(scores))

		var squaredDiffSum float64
		for _, s := range scores {
			diff := s - mean
			squaredDiffSum += diff * diff
		}
		stdDev := 0.0
		if len(scores) > 1 {
			stdDev = squaredDiffSum / float64(len(scores)-1)
			stdDev = math.Sqrt(stdDev)
		}

		// Calculate spike threshold based on symptom score differences
		type ScoredDay struct {
			Date  time.Time
			Score float64
		}
		var scoredDays []ScoredDay
		for _, sym := range symptomsData {
			score := float64(sym.Nausea.Int32+sym.Fatigue.Int32+sym.Pain.Int32) / 3.0
			scoredDays = append(scoredDays, ScoredDay{Date: sym.Date.Time, Score: score})
		}
		sort.Slice(scoredDays, func(i, j int) bool {
			return scoredDays[i].Date.Before(scoredDays[j].Date)
		})

		var diffs []float64
		for i := 1; i < len(scoredDays); i++ {
			diff := scoredDays[i].Score - scoredDays[i-1].Score
			diffs = append(diffs, diff)
		}
		var sumDiff float64
		for _, d := range diffs {
			sumDiff += d
		}
		meanDiff := sumDiff / float64(len(diffs))

		var sqSumDiff float64
		for _, d := range diffs {
			sqSumDiff += (d - meanDiff) * (d - meanDiff)
		}
		stdDiff := math.Sqrt(sqSumDiff / float64(len(diffs)))

		threshold := meanDiff + stdDiff

		// Find spike days based on diff threshold, keep symptom severity for spike day
		spikeDays := make(map[string]float64) // date => symptom severity
		for i := 1; i < len(scoredDays); i++ {
			diff := scoredDays[i].Score - scoredDays[i-1].Score
			if diff > threshold {
				dateStr := scoredDays[i].Date.Format("2006-01-02")
				spikeDays[dateStr] = scoredDays[i].Score
			}
		}

		// Check triggers on the day before spike days
		for spikeDateStr, severity := range spikeDays {
			spikeDate, _ := time.Parse("2006-01-02", spikeDateStr)
			dayBefore := spikeDate.AddDate(0, 0, -1).Format("2006-01-02")

			if sleep, ok := sleepMap[dayBefore]; ok {
				if sleep.Duration.Float64 < 6 {
					triggers.LowSleepHours++
					lowSleepDetails = append(lowSleepDetails, TriggerDetail{Date: dayBefore, TriggerSeverity: severity})
				}
			}

			if diets, ok := dietMap[dayBefore]; ok {
				for _, d := range diets {
					for _, item := range d.Items {
						triggers.FoodItems[item]++
						foodItemDetails[item] = append(foodItemDetails[item], TriggerDetail{Date: dayBefore, TriggerSeverity: severity})
					}
				}
			}

			if menstrual, ok := menstrualMap[dayBefore]; ok {
				triggers.MenstrualEvent[menstrual.PeriodEvent.String]++
				menstrualEventDetails[menstrual.PeriodEvent.String] = append(menstrualEventDetails[menstrual.PeriodEvent.String], TriggerDetail{Date: dayBefore, TriggerSeverity: severity})

				triggers.FlowLevel[menstrual.FlowLevel.String]++
				flowLevelDetails[menstrual.FlowLevel.String] = append(flowLevelDetails[menstrual.FlowLevel.String], TriggerDetail{Date: dayBefore, TriggerSeverity: severity})
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"symptom_spike_threshold": threshold,
			"symptom_average":         mean,
			"standard_deviation":      stdDev,

			"low_sleep_hours": map[string]interface{}{
				"count":   triggers.LowSleepHours,
				"details": lowSleepDetails,
			},
			"common_food_items": map[string]interface{}{
				"counts":  triggers.FoodItems,
				"details": foodItemDetails,
			},
			"menstrual_events": map[string]interface{}{
				"counts":  triggers.MenstrualEvent,
				"details": menstrualEventDetails,
			},
			"flow_levels": map[string]interface{}{
				"counts":  triggers.FlowLevel,
				"details": flowLevelDetails,
			},
		})
	})

	r.GET("/predict_flareups", func(c *gin.Context) {
		queries := database.New(pool)

		sleepData, err := queries.GetAllSleep(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		dietData, err := queries.GetAllDiet(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		menstrualData, err := queries.GetAllMenstrual(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		symptomsData, err := queries.GetAllSymptoms(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		type triggerCounts struct {
			LowSleepHours  int
			MenstrualEvent map[string]int
			FlowLevel      map[string]int
			FoodItems      map[string]int
		}

		type TriggerDetail struct {
			Date            string  `json:"date"`
			TriggerSeverity float64 `json:"trigger_severity"`
		}

		triggers := triggerCounts{
			MenstrualEvent: make(map[string]int),
			FlowLevel:      make(map[string]int),
			FoodItems:      make(map[string]int),
		}

		// Track details per trigger for output
		var lowSleepDetails []TriggerDetail
		foodItemDetails := map[string][]TriggerDetail{}
		menstrualEventDetails := map[string][]TriggerDetail{}
		flowLevelDetails := map[string][]TriggerDetail{}

		// Map data by date
		sleepMap := map[string]database.Sleep{}
		for _, s := range sleepData {
			sleepMap[s.Date.Time.Format("2006-01-02")] = s
		}

		dietMap := map[string][]database.Diet{}
		for _, d := range dietData {
			date := d.Date.Time.Format("2006-01-02")
			dietMap[date] = append(dietMap[date], d)
		}

		menstrualMap := map[string]database.Menstrual{}
		for _, m := range menstrualData {
			menstrualMap[m.Date.Time.Format("2006-01-02")] = m
		}

		// Calculate mean and std dev of symptom severity
		var scores []float64
		for _, sym := range symptomsData {
			avg := float64(sym.Nausea.Int32+sym.Fatigue.Int32+sym.Pain.Int32) / 3.0
			scores = append(scores, avg)
		}
		if len(scores) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No symptom data found."})
			return
		}

		var sum float64
		for _, s := range scores {
			sum += s
		}
		mean := sum / float64(len(scores))

		var squaredDiffSum float64
		for _, s := range scores {
			diff := s - mean
			squaredDiffSum += diff * diff
		}
		stdDev := 0.0
		if len(scores) > 1 {
			stdDev = squaredDiffSum / float64(len(scores)-1)
			stdDev = math.Sqrt(stdDev)
		}

		// Calculate spike threshold based on symptom score differences
		type ScoredDay struct {
			Date  time.Time
			Score float64
		}
		var scoredDays []ScoredDay
		for _, sym := range symptomsData {
			score := float64(sym.Nausea.Int32+sym.Fatigue.Int32+sym.Pain.Int32) / 3.0
			scoredDays = append(scoredDays, ScoredDay{Date: sym.Date.Time, Score: score})
		}
		sort.Slice(scoredDays, func(i, j int) bool {
			return scoredDays[i].Date.Before(scoredDays[j].Date)
		})

		var diffs []float64
		for i := 1; i < len(scoredDays); i++ {
			diff := scoredDays[i].Score - scoredDays[i-1].Score
			diffs = append(diffs, diff)
		}
		var sumDiff float64
		for _, d := range diffs {
			sumDiff += d
		}
		meanDiff := sumDiff / float64(len(diffs))

		var sqSumDiff float64
		for _, d := range diffs {
			sqSumDiff += (d - meanDiff) * (d - meanDiff)
		}
		stdDiff := math.Sqrt(sqSumDiff / float64(len(diffs)))

		threshold := meanDiff + stdDiff

		// Find spike days based on diff threshold, keep symptom severity for spike day
		spikeDays := make(map[string]float64) // date => symptom severity
		for i := 1; i < len(scoredDays); i++ {
			diff := scoredDays[i].Score - scoredDays[i-1].Score
			if diff > threshold {
				dateStr := scoredDays[i].Date.Format("2006-01-02")
				spikeDays[dateStr] = scoredDays[i].Score
			}
		}

		// Check triggers on the day before spike days
		for spikeDateStr, severity := range spikeDays {
			spikeDate, _ := time.Parse("2006-01-02", spikeDateStr)
			dayBefore := spikeDate.AddDate(0, 0, -1).Format("2006-01-02")

			if sleep, ok := sleepMap[dayBefore]; ok {
				if sleep.Duration.Float64 < 6 {
					triggers.LowSleepHours++
					lowSleepDetails = append(lowSleepDetails, TriggerDetail{Date: dayBefore, TriggerSeverity: severity})
				}
			}

			if diets, ok := dietMap[dayBefore]; ok {
				for _, d := range diets {
					for _, item := range d.Items {
						triggers.FoodItems[item]++
						foodItemDetails[item] = append(foodItemDetails[item], TriggerDetail{Date: dayBefore, TriggerSeverity: severity})
					}
				}
			}

			if menstrual, ok := menstrualMap[dayBefore]; ok {
				triggers.MenstrualEvent[menstrual.PeriodEvent.String]++
				menstrualEventDetails[menstrual.PeriodEvent.String] = append(menstrualEventDetails[menstrual.PeriodEvent.String], TriggerDetail{Date: dayBefore, TriggerSeverity: severity})

				triggers.FlowLevel[menstrual.FlowLevel.String]++
				flowLevelDetails[menstrual.FlowLevel.String] = append(flowLevelDetails[menstrual.FlowLevel.String], TriggerDetail{Date: dayBefore, TriggerSeverity: severity})
			}
		}

		// Check if any of these triggers have happened in the last 3 days of the data
		recentSleep := make(map[string]database.Sleep)
		for i := len(sleepData) - 3; i < len(sleepData); i++ {
			if i >= 0 {
				s := sleepData[i]
				recentSleep[s.Date.Time.Format("2006-01-02")] = s
			}
		}
		recentDiet := make(map[string][]database.Diet)
		for i := len(dietData) - 3; i < len(dietData); i++ {
			if i >= 0 {
				d := dietData[i]
				date := d.Date.Time.Format("2006-01-02")
				recentDiet[date] = append(recentDiet[date], d)
			}
		}
		recentMenstrual := make(map[string]database.Menstrual)
		for i := len(menstrualData) - 3; i < len(menstrualData); i++ {
			if i >= 0 {
				m := menstrualData[i]
				recentMenstrual[m.Date.Time.Format("2006-01-02")] = m
			}
		}
		recentSymptoms := make(map[string]database.Symptom)
		for i := len(symptomsData) - 3; i < len(symptomsData); i++ {
			if i >= 0 {
				s := symptomsData[i]
				recentSymptoms[s.Date.Time.Format("2006-01-02")] = s
			}
		}

		var recentFlareupPredictions []string
		for date := range recentSleep {
			if sleep, ok := recentSleep[date]; ok {
				if sleep.Duration.Float64 < 6 {
					recentFlareupPredictions = append(recentFlareupPredictions, fmt.Sprintf("Low sleep hours on %s", date))
				}
			}

			if diets, ok := recentDiet[date]; ok {
				for _, d := range diets {
					for _, item := range d.Items {
						recentFlareupPredictions = append(recentFlareupPredictions, fmt.Sprintf("%s consumed on %s", strings.Title(item), date))
					}
				}
			}

			if menstrual, ok := recentMenstrual[date]; ok {
				recentFlareupPredictions = append(recentFlareupPredictions, fmt.Sprintf("Menstrual event %s on %s", menstrual.PeriodEvent.String, date))
				recentFlareupPredictions = append(recentFlareupPredictions, fmt.Sprintf("Flow level %s on %s", menstrual.FlowLevel.String, date))
			}

			if sym, ok := recentSymptoms[date]; ok {
				avgSeverity := float64(sym.Nausea.Int32+sym.Fatigue.Int32+sym.Pain.Int32) / 3.0
				if avgSeverity > mean+stdDev { // Predict flareup if above average severity
					recentFlareupPredictions = append(recentFlareupPredictions, fmt.Sprintf("High symptom severity on %s: %.2f", date, avgSeverity))
				}
			}
		}

		if len(recentFlareupPredictions) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No recent flareup predictions found."})
			return
		}

		// Calculate probability of flareup based on recent data, and severity of triggers
		var totalTriggers int
		for _, count := range triggers.FoodItems {
			totalTriggers += count
		}
		totalTriggers += triggers.LowSleepHours
		for _, count := range triggers.MenstrualEvent {
			totalTriggers += count
		}
		for _, count := range triggers.FlowLevel {
			totalTriggers += count
		}
		if totalTriggers == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No triggers found in recent data."})
			return
		}
		probability := float64(totalTriggers) / float64(len(recentFlareupPredictions))
		probability = math.Min(probability, 1.0)        // Cap at 100%
		probability *= 100                              // Convert to percentage
		probability = math.Round(probability*100) / 100 // Round to 2 decimal places
		c.JSON(http.StatusOK, gin.H{
			"flareup_probability": probability,
			"flareup_predictions": recentFlareupPredictions,
		})
	})

	r.GET("recommendations", func(c *gin.Context) {
		queries := database.New(pool)

		sleepData, err := queries.GetAllSleep(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		dietData, err := queries.GetAllDiet(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		menstrualData, err := queries.GetAllMenstrual(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		symptomsData, err := queries.GetAllSymptoms(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		type triggerCounts struct {
			LowSleepHours  int
			MenstrualEvent map[string]int
			FlowLevel      map[string]int
			FoodItems      map[string]int
		}

		type TriggerDetail struct {
			Date            string  `json:"date"`
			TriggerSeverity float64 `json:"trigger_severity"`
		}

		triggers := triggerCounts{
			MenstrualEvent: make(map[string]int),
			FlowLevel:      make(map[string]int),
			FoodItems:      make(map[string]int),
		}

		// Track details per trigger for output
		var lowSleepDetails []TriggerDetail
		foodItemDetails := map[string][]TriggerDetail{}
		menstrualEventDetails := map[string][]TriggerDetail{}
		flowLevelDetails := map[string][]TriggerDetail{}

		// Map data by date
		sleepMap := map[string]database.Sleep{}
		for _, s := range sleepData {
			sleepMap[s.Date.Time.Format("2006-01-02")] = s
		}

		dietMap := map[string][]database.Diet{}
		for _, d := range dietData {
			date := d.Date.Time.Format("2006-01-02")
			dietMap[date] = append(dietMap[date], d)
		}

		menstrualMap := map[string]database.Menstrual{}
		for _, m := range menstrualData {
			menstrualMap[m.Date.Time.Format("2006-01-02")] = m
		}

		// Calculate mean and std dev of symptom severity
		var scores []float64
		for _, sym := range symptomsData {
			avg := float64(sym.Nausea.Int32+sym.Fatigue.Int32+sym.Pain.Int32) / 3.0
			scores = append(scores, avg)
		}
		if len(scores) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No symptom data found."})
			return
		}

		var sum float64
		for _, s := range scores {
			sum += s
		}
		mean := sum / float64(len(scores))

		var squaredDiffSum float64
		for _, s := range scores {
			diff := s - mean
			squaredDiffSum += diff * diff
		}
		stdDev := 0.0
		if len(scores) > 1 {
			stdDev = squaredDiffSum / float64(len(scores)-1)
			stdDev = math.Sqrt(stdDev)
		}

		// Calculate spike threshold based on symptom score differences
		type ScoredDay struct {
			Date  time.Time
			Score float64
		}
		var scoredDays []ScoredDay
		for _, sym := range symptomsData {
			score := float64(sym.Nausea.Int32+sym.Fatigue.Int32+sym.Pain.Int32) / 3.0
			scoredDays = append(scoredDays, ScoredDay{Date: sym.Date.Time, Score: score})
		}
		sort.Slice(scoredDays, func(i, j int) bool {
			return scoredDays[i].Date.Before(scoredDays[j].Date)
		})

		var diffs []float64
		for i := 1; i < len(scoredDays); i++ {
			diff := scoredDays[i].Score - scoredDays[i-1].Score
			diffs = append(diffs, diff)
		}
		var sumDiff float64
		for _, d := range diffs {
			sumDiff += d
		}
		meanDiff := sumDiff / float64(len(diffs))

		var sqSumDiff float64
		for _, d := range diffs {
			sqSumDiff += (d - meanDiff) * (d - meanDiff)
		}
		stdDiff := math.Sqrt(sqSumDiff / float64(len(diffs)))

		threshold := meanDiff + stdDiff

		// Find spike days based on diff threshold, keep symptom severity for spike day
		spikeDays := make(map[string]float64) // date => symptom severity
		for i := 1; i < len(scoredDays); i++ {
			diff := scoredDays[i].Score - scoredDays[i-1].Score
			if diff > threshold {
				dateStr := scoredDays[i].Date.Format("2006-01-02")
				spikeDays[dateStr] = scoredDays[i].Score
			}
		}

		// Check triggers on the day before spike days
		for spikeDateStr, severity := range spikeDays {
			spikeDate, _ := time.Parse("2006-01-02", spikeDateStr)
			dayBefore := spikeDate.AddDate(0, 0, -1).Format("2006-01-02")

			if sleep, ok := sleepMap[dayBefore]; ok {
				if sleep.Duration.Float64 < 6 {
					triggers.LowSleepHours++
					lowSleepDetails = append(lowSleepDetails, TriggerDetail{Date: dayBefore, TriggerSeverity: severity})
				}
			}

			if diets, ok := dietMap[dayBefore]; ok {
				for _, d := range diets {
					for _, item := range d.Items {
						triggers.FoodItems[item]++
						foodItemDetails[item] = append(foodItemDetails[item], TriggerDetail{Date: dayBefore, TriggerSeverity: severity})
					}
				}
			}

			if menstrual, ok := menstrualMap[dayBefore]; ok {
				triggers.MenstrualEvent[menstrual.PeriodEvent.String]++
				menstrualEventDetails[menstrual.PeriodEvent.String] = append(menstrualEventDetails[menstrual.PeriodEvent.String], TriggerDetail{Date: dayBefore, TriggerSeverity: severity})

				triggers.FlowLevel[menstrual.FlowLevel.String]++
				flowLevelDetails[menstrual.FlowLevel.String] = append(flowLevelDetails[menstrual.FlowLevel.String], TriggerDetail{Date: dayBefore, TriggerSeverity: severity})
			}
		}

		temp := float32(1)
		// Example output something like ["avoid inflammatory foods", "increase hydration", "improve sleep hygiene"], only 3
		result, err := client.Models.GenerateContent(ctx2, "gemini-2.5-flash-lite", genai.Text(`Be short and concise, and specific. Return an array of 3 recommendations to reduce flare-ups based on the following data:
			Sleep Data: `+fmt.Sprintf("%v", sleepData)+
			`Diet Data: `+fmt.Sprintf("%v", dietData)+
			`Menstrual Data: `+fmt.Sprintf("%v", menstrualData)+
			`Symptoms Data: `+fmt.Sprintf("%v", symptomsData)+
			`Triggers: `+fmt.Sprintf("%v", triggers)), &genai.GenerateContentConfig{
			SystemInstruction: &genai.Content{
				Role: "Output in the format of a JSON array with 3 items. Example: [\"recommendation1\", \"recommendation2\", \"recommendation3\"]. Output only the json array nothing more. Be very short and concise.",
			},
			Temperature:      &temp,
			MaxOutputTokens:  200,
			ResponseMIMEType: "application/json",
			ResponseSchema: &genai.Schema{
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeString,
				},
			},
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(result.Candidates) == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No recommendations generated"})
			return
		}

		recommendations := result.Text()
		c.String(http.StatusOK, recommendations)
	})

	fmt.Printf("Server is running on http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
