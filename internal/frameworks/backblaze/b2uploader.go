package backblaze

import (
	"errors"
	"fmt"
	"io"

	"gopkg.in/kothar/go-backblaze.v0"
)

var accountID = "635d33259f1e"
var applicationKEY = "005b52216b4aaa94e054e59b99e37fffdc87ee0fde"
var bucketName = "sma-bucket"

func connectBackBlazeBucket() (*backblaze.Bucket, error) {
	b2, _ := backblaze.NewB2(backblaze.Credentials{
		ApplicationKey: applicationKEY,
		KeyID:          accountID,
	})

	response, err := b2.Bucket(bucketName)

	if err != nil {
		return &backblaze.Bucket{}, errors.New("erro ao listar bucket")
	}

	return response, nil
}

func UploadFileService(fileName string, fileContent io.Reader) (string, error) {

	bucket, err := connectBackBlazeBucket()

	if err != nil {
		return "", err
	}
	metadata := make(map[string]string)

	response, err := bucket.UploadFile(fileName, metadata, fileContent)

	if err != nil {
		return "", fmt.Errorf("falha ao fazer upload: %w", err)
	}

	fileNameURL, err := bucket.FileURL(response.Name)

	if err != nil {
		return "", fmt.Errorf("falha ao fazer download: %w", err)
	}

	return fileNameURL, nil
}
