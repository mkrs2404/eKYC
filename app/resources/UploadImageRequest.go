package resources

import (
	"mime/multipart"

	"github.com/go-playground/validator/v10"
)

type UploadImageRequest struct {
	ImageType string                `form:"type" binding:"required" validate:"oneof=face id_card"`
	Image     *multipart.FileHeader `form:"file" binding:"required" swaggerignore:"true" swaggertype:"file"`
}

//Validate validates the Sign Up request parameters
func (req *UploadImageRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(req)
	return err
}
