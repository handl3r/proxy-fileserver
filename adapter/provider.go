package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"io/ioutil"
	"net/http"
	"os"
	"proxy-fileserver/common/log"
	"proxy-fileserver/configs"
	"proxy-fileserver/enums"
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
	var service *drive.Service
	var err error
	if config.GoogleDriveOAuthConfig.Enable {
		credentials, err := ioutil.ReadFile(config.GoogleDriveOAuthConfig.CredentialFile)
		if err != nil {
			return nil, err
		}
		gConfig, err := google.ConfigFromJSON(credentials, drive.DriveReadonlyScope, drive.DriveMetadataScope)
		if err != nil {
			return nil, err
		}
		client := getDriveClient(gConfig, config.GoogleDriveOAuthConfig.TokenFile)
		service, err = drive.New(client)
		if err != nil {
			log.Errorf("Can not init service google drive with client, error: %v", err)
			return nil, err
		}
	} else {
		service, err = drive.NewService(ctx)
		if err != nil {
			log.Errorf("Can not init service google drive application credential, error: %v", err)
			return nil, err
		}
	}
	googleDriveFileSystem := NewGoogleDriveFileSystem(ctx, service, config.SharedRootFolder, config.SharedRootFolderID)
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

func getDriveClient(config *oauth2.Config, tokenLocation string) *http.Client {
	var token *oauth2.Token
	var err error
	token, err = getTokenFromFile(tokenLocation)
	if err != nil {
		log.Errorf("Can not get G Oauth2 Token from file: %s, error: %v", tokenLocation, err)
		token, err = getTokenFromCallback(config)
		if err != nil {
			log.Errorf("Can not get G OAuth2 Token from Callback with error: %v", err)
			return nil
		}
		err = saveToken(tokenLocation, token)
		if err != nil {
			log.Errorf("Can not save G Oauth2 Token to file: %s", tokenLocation)
		}
	}
	return config.Client(context.Background(), token)

}

func getTokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func getTokenFromCallback(config *oauth2.Config) (*oauth2.Token, error) {
	log.Infof("Get G OAuth2 Token from callback")
	authURL := config.AuthCodeURL(enums.StateToken, oauth2.AccessTypeOffline)
	log.Infof("Access following link[%s], grant permission then type authorization here: ", authURL)
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, err
	}
	tok, err := config.Exchange(context.TODO(), authCode)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) error {
	log.Infof("Save new G Oauth2 Token to %s", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}
