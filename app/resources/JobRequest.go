package resources

type JobRequest struct {
	JobId int `json:"job_id" binding:"required"`
}
