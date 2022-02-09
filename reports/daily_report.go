package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/mkrs2404/eKYC/api/models"
	"github.com/mkrs2404/eKYC/database"
	"gorm.io/gorm/logger"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error fetching the environment values")
	}
	database.Connect(os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), logger.Error)
}

func main() {
	getDailyReportData()
}

func getDailyReportData() {

	format := "02-01-2006"
	today := time.Now().Format(format)
	var clients []models.Client
	var dailyDataList []DailyData
	var plans []models.Plan

	planMap := make(map[uint]models.Plan)

	database.DB.Find(&clients)
	database.DB.Find(&plans)

	for _, plan := range plans {
		planMap[plan.ID] = plan
	}

	for _, client := range clients {
		var dailyData DailyData
		var totalFaceMatch int64
		var totalOcr int64
		var totalStorageKb float64
		plan := planMap[client.Plan]

		database.DB.Model(&models.Api_Calls{}).Where("client_id = ? AND type = ?", client.ID, "face-match").Count(&totalFaceMatch)
		database.DB.Model(&models.Api_Calls{}).Where("client_id = ? AND type = ?", client.ID, "ocr").Count(&totalOcr)
		database.DB.Model(&models.File{}).Select("sum(file_size_kb)").Where("client_id = ?", client.ID).Find(&totalStorageKb)

		dailyData.ClientId = strconv.Itoa(int(client.ID))
		dailyData.Name = client.Name
		dailyData.Date = today
		dailyData.Plan = plan.PlanName
		dailyData.TotalFacematch = int(totalFaceMatch)
		dailyData.TotalOcr = int(totalOcr)
		dailyData.TotalImageStorageInMb = totalStorageKb / (1 << 10)
		dailyData.ApiUsageCostUsd = plan.DailyBaseCost + float64(totalFaceMatch+totalOcr)*plan.ApiCost
		dailyData.StorageCostUsd = dailyData.StorageCostUsd * plan.StorageCost

		dailyDataList = append(dailyDataList, dailyData)
	}

	csvFile, err := os.Create(today + ".csv")
	if err != nil {
		log.Fatalf("File creation failed: %s", err)
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)

	var data [][]string
	header := []string{"client_id", "name", "plan", "date", "total_facematch", "total_ocr", "total_image_storage_in_mb", "api_usage_cost_usd", "storage_cost_usd"}
	data = append(data, header)
	for _, record := range dailyDataList {
		totalImageStorageInMb := fmt.Sprintf("%f", record.TotalImageStorageInMb)
		apiUsageCostUsd := fmt.Sprintf("%f", record.ApiUsageCostUsd)
		storageCostUsd := fmt.Sprintf("%f", record.StorageCostUsd)

		row := []string{record.ClientId, record.Name, record.Plan, record.Date, strconv.Itoa(record.TotalFacematch), strconv.Itoa(record.TotalOcr), totalImageStorageInMb, apiUsageCostUsd, storageCostUsd}
		data = append(data, row)
	}

	err = csvWriter.WriteAll(data)
	if err != nil {
		log.Fatal(err)
	}
}
