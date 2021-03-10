package services

import (
	"io"
	"proxy-fileserver/adapter"
	"proxy-fileserver/common/log"
	"proxy-fileserver/enums"
	"proxy-fileserver/models"
	"proxy-fileserver/repository"
	"strings"
	"time"
)

type FileSystemService struct {
	GoogleDriveFileSystem *adapter.GoogleDriveFileSystem
	LocalFileSystem       *adapter.LocalFileSystem

	FileInForRepo *repository.FileInfoRepository
}

func NewFileSystemService(googleDrive *adapter.GoogleDriveFileSystem, localStorage *adapter.LocalFileSystem, fileInFor *repository.FileInfoRepository) *FileSystemService {
	return &FileSystemService{
		GoogleDriveFileSystem: googleDrive,
		LocalFileSystem:       localStorage,
		FileInForRepo:         fileInFor,
	}
}

// Use for gin
func (s *FileSystemService) GetSourceStream(filePath string) (io.Reader, enums.Response) {
	isProcessing, err := s.LocalFileSystem.IsExisted(filePath + enums.SuffixProcessing)
	if err != nil {
		log.Errorf("Failure when checking if file exist from local file system %s with error: %v",
			filePath+enums.SuffixProcessing, err)
		return nil, enums.ErrorSystem
	}
	existed, err := s.LocalFileSystem.IsExisted(filePath)
	if err != nil {
		log.Errorf("Failure when checking if file exist from local file system %s with error: %v", filePath, err)
		return nil, enums.ErrorSystem
	}
	if existed {
		srcStream, err := s.GetSourceStreamFromLocalFileSystem(filePath)
		if err != nil {
			log.Errorf("Failure getting source stream file on local file system with filepath: %s, err: %s", filePath, err)
			return nil, enums.ErrorSystem
		}
		go func() {
			now := time.Now()
			err := s.FileInForRepo.Update(models.FileInfo{
				FilePath:       filePath,
				LastDownloadAt: now,
			})
			if err != nil {
				log.Errorf("Can not update last_download_at for file %s with error: %v", filePath, err)
			}
			log.Infof("Updated last_download_at for file %s, last_download_at: %v", filePath, now)
		}()
		return srcStream, nil
	}
	if isProcessing {
		id, srcStream, err := s.GoogleDriveFileSystem.GetStreamSourceByFilePath(filePath)
		if err == enums.ErrFileNotExist {
			return nil, enums.ErrorNoContent
		}
		if err != nil {
			log.Errorf("Failure getting source stream from drive with path %s, id %s , error: %v", filePath, id, err)
			return nil, enums.ErrorSystem
		}
		return srcStream, nil
	}
	id, srcStream, err := s.GoogleDriveFileSystem.GetStreamSourceByFilePath(filePath)
	if err == enums.ErrFileNotExist {
		return nil, enums.ErrorNoContent
	}
	if err != nil {
		log.Errorf("Failure getting source stream from drive with path %s, id %s , error: %v", filePath, id, err)
		return nil, enums.ErrorSystem
	}
	go func() {
		log.Infof("Start streaming file from drive to local file system with filepath: %s, id: %s", filePath, id)
		err = s.StreamFromDriveToLocalFileSystem(id, filePath)
		if err != nil {
			log.Errorf("Failure when streaming file from drive to local file system with filepath: %s, id: %s, err: %v", filePath, id, err)
			return
		}
		log.Infof("Finished streaming file from drive to local file system with filepath: %s, id: %s", filePath, id)
	}()
	go func() {
		now := time.Now()
		err := s.FileInForRepo.Create(models.FileInfo{
			FilePath:       filePath,
			LastDownloadAt: now,
		})
		if err != nil {
			log.Errorf("Can not create last_download_at for file %s with error: %v", filePath, err)
			return
		}
		log.Infof("Created last_download_at for file %s, last_download_at: %v", filePath, now)
	}()
	return srcStream, nil
}

// Use only for HTTP Basic
// StreamFile Public method control all process to stream file to client and sync file from drive to local server
func (s *FileSystemService) StreamFile(outStreamHttp io.Writer, filePath string) enums.Response {
	existed, err := s.LocalFileSystem.IsExisted(filePath)
	if err != nil {
		log.Errorf("Can not check if file exist from file %s with error: %v", filePath, err)
		return enums.ErrorSystem
	}
	if existed {
		err := s.StreamFromLocalFileSystem(outStreamHttp, filePath)
		if err != nil {
			log.Errorf("Failure when streaming file from local server to client with filepath: %s", filePath)
			return enums.ErrorSystem
		}
		log.Infof("Finish stream file %s from  local file system to client", filePath)
		return nil
	}

	id, srcStream, err := s.GoogleDriveFileSystem.GetStreamSourceByFilePath(filePath)
	if err == enums.ErrFileNotExist {
		return enums.ErrorNoContent
	}
	if err != nil {
		log.Errorf("Can not get source stream from drive with file path %s, id %s , error: %v", filePath, id, err)
		return enums.ErrorSystem
	}
	// TODO make new function to stream from reader to multi writer with condition: when 1 writer is failure, another still continue
	go func() {
		err = s.StreamFromDriveToLocalFileSystem(id, filePath)
		if err != nil {
			log.Errorf("Failure when streaming file from drive to local file system with filepath: %s, id: %s, err: %v", filePath, id, err)
			return
		}
		log.Infof("Finished getting new file from drive to local server with filepath: %s, id: %s", filePath, id)
	}()
	_, err = io.Copy(outStreamHttp, srcStream)
	if err != nil {
		log.Errorf("Error when streaming file from driver to client with filePath: %s, id: %s, error: %s", filePath, id, err)
		return enums.ErrorSystem
	}
	log.Infof("Finished stream file from drive to client with filepath: %s, id: %s", filePath, id)
	return nil
}

// Use for gin
func (s *FileSystemService) GetSourceStreamFromLocalFileSystem(filePath string) (io.Reader, error) {
	srcStream, err := s.LocalFileSystem.GetStreamSourceByFilePath(filePath)
	if err != nil {
		return nil, err
	}
	return srcStream, nil
}

// Use for http basic
func (s *FileSystemService) StreamFromLocalFileSystem(outStream io.Writer, filePath string) error {
	srcStream, err := s.GetSourceStreamFromLocalFileSystem(filePath)
	if err != nil {
		return err
	}
	_, err = io.Copy(outStream, srcStream)
	if err != nil {
		return err
	}
	return nil
}

func (s *FileSystemService) StreamFromDriveToLocalFileSystem(id string, filePath string) error {
	processingFile := filePath + enums.SuffixProcessing
	newFileStream, err := s.LocalFileSystem.NewFile(processingFile)
	if err != nil {
		return err
	}
	defer func() {
		err := s.LocalFileSystem.RenameFile(processingFile, filePath)
		if err != nil {
			log.Errorf("Failure when rename file %s to %s after stream file to local server with error: %v",
				processingFile, filePath, err)
		}
	}()
	srcStreamDrive, err := s.GoogleDriveFileSystem.GetStreamBySourceByID(id)
	if err != nil {
		return err
	}
	_, err = io.Copy(newFileStream, srcStreamDrive)
	if err != nil {
		return err
	}
	return nil
}

// FilePath format: x/y/z/t.a
func (s *FileSystemService) validateFilePath(filePath string) bool {
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
