package resources

type FaceMatchRequest struct {
	Image1 string `json:"image1" binding:"required"`
	Image2 string `json:"image2" binding:"required"`
}
