package adapter

import (
	"context"
	"proxy-fileserver/configs"
)

type ProviderAdapter interface {
	GetGoogleDriveFileSystem() *GoogleDriveFileSystem
	GetLocalFileSystem() *LocalFileSystem
}

type providerAdapterImpl struct {
	googleDriveFileSystem *GoogleDriveFileSystem
	localFileSystem       *LocalFileSystem
}

func NewProviderAdapter(ctx context.Context, config *configs.Config) (ProviderAdapter, error) {
	googleDriveFileSystem, err := NewGoogleDriveFileSystem(ctx, config.SharedRootFolder, config.SharedRootFolderID)
	if err != nil {
		return nil, err
	}
	localFileSystem := NewLocalFileSystem(config.SharedRootFolderLocal)
	return &providerAdapterImpl{
		googleDriveFileSystem: googleDriveFileSystem,
		localFileSystem:       localFileSystem,
	}, nil
}

func (p *providerAdapterImpl) GetGoogleDriveFileSystem() *GoogleDriveFileSystem {
	return p.googleDriveFileSystem
}

func (p *providerAdapterImpl) GetLocalFileSystem() *LocalFileSystem {
	return p.localFileSystem
}
