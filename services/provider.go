package services

import (
	"proxy-fileserver/adapter"
	"proxy-fileserver/repository"
	"time"
)

type ServiceProvider interface {
	GeFileSystemService() *FileSystemService
	GetAuthService() *AuthService
}

type serviceProviderImpl struct {
	fileSystemService *FileSystemService
	authService       *AuthService
}

func NewServiceProvider(adapterProvider adapter.ProviderAdapter, repositoryProvider repository.ProviderRepository,
	privateKeyLocation, publicKeyLocation string, expiredTime time.Duration, sharedFolder string,
) ServiceProvider {
	return &serviceProviderImpl{
		fileSystemService: NewFileSystemService(
			adapterProvider.GetGoogleDriveFileSystem(),
			adapterProvider.GetLocalFileSystem(),
			repositoryProvider.GetFileInfoRepository(),
			sharedFolder,
		),
		authService: NewAuthService(privateKeyLocation, publicKeyLocation, expiredTime),
	}
}

func (s *serviceProviderImpl) GeFileSystemService() *FileSystemService {
	return s.fileSystemService
}
func (s *serviceProviderImpl) GetAuthService() *AuthService {
	return s.authService
}
