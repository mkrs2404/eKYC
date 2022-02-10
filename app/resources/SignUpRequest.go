package resources

import "github.com/go-playground/validator/v10"

//Structure for the Sign Up request
type SignUpRequest struct {
	Name  string `json:"name"`
	Email string `json:"email" binding:"required" validate:"email"`
	Plan  string `json:"plan" binding:"required" validate:"oneof=basic advanced enterprise"`
}

//Validate validates the Sign Up request parameters
func (req *SignUpRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(req)
	return err
}
