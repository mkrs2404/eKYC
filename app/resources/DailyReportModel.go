package resources

type DailyData struct {
	ClientId              string
	Name                  string
	Plan                  string
	Date                  string
	TotalFacematch        int
	TotalOcr              int
	TotalImageStorageInMb float64
	ApiUsageCostUsd       float64
	StorageCostUsd        float64
}
