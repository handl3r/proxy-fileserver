package services

import (
	"proxy-fileserver/adapter"
)

type ServiceProvider interface {
	GeFileSystemService() *FileSystemService
}

type serviceProviderImpl struct {
	fileSystemService *FileSystemService
}

func NewServiceProvider(adapterProvider adapter.ProviderAdapter) ServiceProvider {
	return &serviceProviderImpl{
		fileSystemService: NewFileSystemService(
			adapterProvider.GetGoogleDriveFileSystem(),
			adapterProvider.GetLocalFileSystem(),
		),
	}
}

func (s *serviceProviderImpl) GeFileSystemService() *FileSystemService {
	return s.fileSystemService
}
