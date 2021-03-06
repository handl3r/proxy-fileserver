package adapter

import (
	"context"
	"fmt"
	"google.golang.org/api/drive/v3"
	"io"
	"proxy-fileserver/enums"
	"strings"
)

type GoogleDriveFileSystem struct {
	service            *drive.Service
	sharedRootFolder   string
	sharedRootFolderID string
}

type TreeNode struct {
	Name string
	ID   string
}

func NewGoogleDriveFileSystem(ctx context.Context, sharedRootFolder, sharedRootFolderID string) (*GoogleDriveFileSystem, error) {
	service, err := drive.NewService(ctx)
	if err != nil {
		return nil, err
	}
	return &GoogleDriveFileSystem{
		service:            service,
		sharedRootFolder:   sharedRootFolder,
		sharedRootFolderID: sharedRootFolderID,
	}, nil
}

// path must be: {shared-folder/*}
func (s *GoogleDriveFileSystem) validateFilePath(filePath string) bool {
	files := strings.Split(filePath, "/")
	if len(files) < 2 {
		return false
	}
	for _, file := range files {
		if len(strings.TrimSpace(file)) == 0 {
			return false
		}
	}
	return true
}

func (s *GoogleDriveFileSystem) buildQuerySearchFile(filePath string) string {
	subQueries := make([]string, 0)
	files := strings.Split(filePath, "/")
	for _, file := range files[0:(len(files) - 1)] {
		subQuery := fmt.Sprintf("(name = '%s' and mimeType = 'application/vnd.google-apps.folder')", file)
		subQueries = append(subQueries, subQuery)
	}
	lastQuery := fmt.Sprintf("(name = '%s' and mimeType =! 'application/vnd.google-apps.folder')", files[len(files)-1])
	subQueries = append(subQueries, lastQuery)
	return strings.Join(subQueries, " or ")
}

// TODO remove search for sharedRootFolder
// Return fileID, isExisted, error
func (s *GoogleDriveFileSystem) GetFileIDByPath(filePath string) (string, bool, error) {
	listFileInPath := strings.Split(filePath, "/")
	numPathLevel := len(listFileInPath)
	query := s.buildQuerySearchFile(filePath)
	fileList, err := s.service.Files.List().Fields("files(id, name, parents)").Q(query).Do()
	if err != nil {
		return "", false, err
	}
	if len(fileList.Files) < numPathLevel {
		return "", false, nil
	}
	//listGroupFileWithName := make([][]TreeNode, numPathLevel-1)
	lastFileID := ""
	isExistedPath := true
	preNodeID := s.sharedRootFolderID
	for _, fileInPath := range listFileInPath[1:] {
		existedNode := false
		for _, file := range fileList.Files {
			if file.Name == fileInPath && file.Id == preNodeID {
				preNodeID = file.Id
				lastFileID = file.Id
				existedNode = true
				break
			}
		}
		if existedNode == false {
			isExistedPath = false
			break
		}
	}
	if isExistedPath {
		return lastFileID, true, nil
	}
	return "", false, nil
}

// Return id, srcStream, error
func (s *GoogleDriveFileSystem) GetStreamSourceByFilePath(filePath string) (string, io.Reader, error) {
	id, existed, err := s.GetFileIDByPath(filePath)
	if err != nil {
		return "", nil, err
	}
	if !existed {
		return "", nil, enums.ErrFileNotExist
	}
	stream, err := s.GetStreamBySourceByID(id)
	return id, stream, err
}

func (s *GoogleDriveFileSystem) IsExistedByID(id string) (bool, error) {
	files, err := s.service.Files.List().Do()
	if err != nil {
		return false, err
	}
	if len(files.Files) == 0 {
		return false, err
	}
	for _, file := range files.Files {

		if file.Id == id {
			return true, nil
		}
	}
	return false, nil
}

func (s *GoogleDriveFileSystem) GetStreamBySourceByID(id string) (io.Reader, error) {
	resp, err := s.service.Files.Get(id).Download()
	if err != nil {
		return nil, err
	}
	return resp.Body, nil

}
