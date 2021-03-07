package services

import (
	"proxy-fileserver/adapter"
	"proxy-fileserver/repository"
)

type ServiceProvider interface {
	GeFileSystemService() *FileSystemService
}

type serviceProviderImpl struct {
	fileSystemService *FileSystemService
}

func NewServiceProvider(adapterProvider adapter.ProviderAdapter, repositoryProvider repository.ProviderRepository) ServiceProvider {
	return &serviceProviderImpl{
		fileSystemService: NewFileSystemService(
			adapterProvider.GetGoogleDriveFileSystem(),
			adapterProvider.GetLocalFileSystem(),
			repositoryProvider.GetFileInfoRepository(),
		),
	}
}

func (s *serviceProviderImpl) GeFileSystemService() *FileSystemService {
	return s.fileSystemService
}
