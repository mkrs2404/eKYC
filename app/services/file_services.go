package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/mkrs2404/eKYC/app/database"
	"github.com/mkrs2404/eKYC/app/messages"
	"github.com/mkrs2404/eKYC/app/models"
)

var allowedFileSize float64 = 5

//SaveFile creates the file object and saves it to the database, returning the file's uuid and err, if any
func SaveFile(bucketName string, filePath string, fileSize int64, fileType string, clientId uint) (uuid.UUID, error) {

	var file models.File
	file.ClientID = clientId
	file.FileName = filepath.Base(filePath)
	file.FileSizeKB = float64(fileSize) / (1 << 10)
	file.FileExtension = filepath.Ext(file.FileName)
	file.FileType = fileType
	file.FileStoragePath = fmt.Sprintf("%s/%s", bucketName, filePath)
	fileNameWithoutExt := strings.Split(file.FileName, ".")[0]
	file.ID, _ = uuid.Parse(fileNameWithoutExt)
	err := database.DB.Create(&file).Error
	return file.ID, err
}

//ValidateFile validates the uploaded file's size and extension
func ValidateFile(file *multipart.FileHeader) error {

	var err error
	fileExt := filepath.Ext(file.Filename)
	fileSizeMb := float64(file.Size) / (1 << 20)

	if !(fileExt == ".png" || fileExt == ".jpeg") {
		err = errors.New(messages.WRONG_FILE_EXTENSION)
	}
	if fileSizeMb > allowedFileSize {
		err = errors.New(messages.WRONG_FILE_SIZE)
	}
	return err
}

func GetFileForClient(imageId string, clientId uint) (models.File, error) {

	var file models.File
	err := database.DB.Where("id = ? AND client_id = ?", imageId, clientId).First(&file).Error
	return file, err
}

func CreateEmptyFile(fileName string) (string, error) {

	os.MkdirAll("./downloads", os.ModePerm)
	file, err := os.Create(fmt.Sprintf("./downloads/%s", fileName))
	filePath := file.Name()
	return filePath, err
}

//deleteLocalFile deletes the file that was saved locally by saveFileToDisk method
func DeleteLocalFile(filePath string) {
	os.Remove(filePath)
}
