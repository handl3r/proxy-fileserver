package services

import (
	"context"
	"google.golang.org/api/drive/v3"
)

type GoogleDriveService struct {
	service *drive.Service
	sharedFolder string
}

func NewGoogleDriveClientWithServiceAccount(ctx context.Context, sharedFolder string) (*GoogleDriveService, error){
	service, err := drive.NewService(ctx)
	if err != nil {
		return nil, err
	}
	return &GoogleDriveService{
		service: service,
		sharedFolder: sharedFolder,
	}, nil
}

func GetFile() {

}
