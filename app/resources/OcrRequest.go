package resources

type OcrRequest struct {
	Image string `json:"image" binding:"required"`
}
