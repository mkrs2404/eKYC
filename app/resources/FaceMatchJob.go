package resources

type FaceMatchJob struct {
	Image1 string
	Image2 string
	JobId  uint
}

func CreateFaceMatchJob(img1 string, img2 string, id uint) *FaceMatchJob {
	var job FaceMatchJob
	job.Image1 = img1
	job.Image2 = img2
	job.JobId = id
	return &job
}
