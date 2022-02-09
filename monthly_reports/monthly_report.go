package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/mkrs2404/eKYC/api/models"
	"github.com/mkrs2404/eKYC/database"
	"gopkg.in/alecthomas/kingpin.v2"
	"gorm.io/gorm/logger"
)

type IntResult struct {
	Date  time.Time
	Total int64
}

type FloatResult struct {
	Date  time.Time
	Total float64
}

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error fetching the environment values")
	}
	database.Connect(os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), logger.Error)
}

func main() {
	clientId := kingpin.Flag("client", "Client ID").Required().Int()
	kingpin.Parse()
	getMonthlyReportData(*clientId)
}

func getMonthlyReportData(clientId int) {
	format := "02-01-2006"
	yyyyFormat := "2006-01-02"
	var client models.Client
	var plans []models.Plan
	reportMap := make(map[time.Time]MonthlyData)
	var mapKeys []time.Time

	planMap := make(map[uint]models.Plan)

	database.DB.Find(&client, clientId)
	database.DB.Find(&plans)

	for _, plan := range plans {
		planMap[plan.ID] = plan
	}

	t := time.Now()
	firstDayTime := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDayTime.AddDate(0, 1, 0).Add(-time.Nanosecond).Format(yyyyFormat)
	firstDay := firstDayTime.Format(yyyyFormat)

	plan := planMap[client.Plan]

	//Getting total face matches by date
	var matchResult []IntResult
	database.DB.Table("api_calls").Select("date(created_at) as date, count(id) as total").Where("client_id = ? and type = ?", clientId, "face-match").Group("date(created_at)").
		Having("date(created_at) >= ? AND date(created_at) <= ?", firstDay, lastDay).Find(&matchResult)

	for _, record := range matchResult {
		mapKeys = append(mapKeys, record.Date)
		reportMap[record.Date] = MonthlyData{
			Date:           record.Date.Format(format),
			TotalFacematch: int(record.Total),
		}
	}

	//Getting total ocr by date
	var ocrResult []IntResult
	database.DB.Table("api_calls").Select("date(created_at) as date, count(id) as total").Where("client_id = ? and type = ?", clientId, "ocr").Group("date(created_at)").
		Having("date(created_at) >= ? AND date(created_at) <= ?", firstDay, lastDay).Find(&ocrResult)

	for _, record := range ocrResult {
		data, ok := reportMap[record.Date]
		if !ok {
			mapKeys = append(mapKeys, record.Date)
			reportMap[record.Date] = MonthlyData{
				Date:     record.Date.Format(format),
				TotalOcr: int(record.Total),
			}
		} else {
			data.TotalOcr = int(record.Total)
			reportMap[record.Date] = data
		}

	}

	//Getting total storage used by date
	var storageResult []FloatResult
	database.DB.Table("files").Select("date(created_at) as date, sum(file_size_kb) as total").Where("client_id = ?", clientId).Group("date(created_at)").
		Having("date(created_at) >= ? AND date(created_at) <= ?", firstDay, lastDay).Find(&storageResult)

	for _, record := range storageResult {
		data, ok := reportMap[record.Date]
		if !ok {
			mapKeys = append(mapKeys, record.Date)
			reportMap[record.Date] = MonthlyData{
				Date:                  record.Date.Format(format),
				TotalImageStorageInMb: record.Total / (1 << 10),
			}
		} else {
			data.TotalImageStorageInMb = record.Total / (1 << 10)
			reportMap[record.Date] = data
		}
	}

	sort.Slice(mapKeys, func(i, j int) bool { return mapKeys[i].Before(mapKeys[j]) })

	fileName := fmt.Sprintf("%s_%d", time.Now().Month(), clientId)
	csvFile, err := os.Create(fileName + ".csv")
	if err != nil {
		log.Fatalf("File creation failed: %s", err)
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)

	var csvData [][]string
	var totalMonthlyCost float64
	clientHeader := fmt.Sprintf("Report for client id: %d", clientId)
	csvData = append(csvData, []string{clientHeader})
	header := []string{"date", "total_facematch", "total_ocr", "total_image_storage_in_mb", "api_usage_cost_usd", "storage_cost_usd"}
	csvData = append(csvData, header)

	for _, record := range mapKeys {
		monthlyData := reportMap[record]
		monthlyData.ApiUsageCostUsd = plan.DailyBaseCost + float64(monthlyData.TotalFacematch+monthlyData.TotalOcr)*plan.ApiCost
		monthlyData.StorageCostUsd = monthlyData.TotalImageStorageInMb * plan.StorageCost

		totalImageStorageInMb := fmt.Sprintf("%f", monthlyData.TotalImageStorageInMb)
		apiUsageCostUsd := fmt.Sprintf("%f", monthlyData.ApiUsageCostUsd)
		storageCostUsd := fmt.Sprintf("%f", monthlyData.StorageCostUsd)

		totalMonthlyCost += monthlyData.ApiUsageCostUsd + monthlyData.StorageCostUsd
		row := []string{monthlyData.Date, strconv.Itoa(monthlyData.TotalFacematch), strconv.Itoa(monthlyData.TotalOcr), totalImageStorageInMb, apiUsageCostUsd, storageCostUsd}
		csvData = append(csvData, row)
	}

	totalCostFooter := fmt.Sprintf("Total Invoice Amount : USD %.2f", totalMonthlyCost)
	csvData = append(csvData, []string{totalCostFooter})
	err = csvWriter.WriteAll(csvData)
	if err != nil {
		log.Fatal(err)
	}

}
