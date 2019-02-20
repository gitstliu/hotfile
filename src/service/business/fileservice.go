package business

import (
	"datafile"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gitstliu/log4go"
)

type FileService struct {
}

func (fileService *FileService) SaveFile(path string, serviceId int64, userId string, group string, source []byte) (string, error) {

	currPath := strings.TrimRight(path, "/")

	fileName, fileId, fileNameErr := datafile.GetFileName(userId, group)
	if fileNameErr != nil {
		return "", fileNameErr
	}

	currDate := time.Now().Format("2006-01-02")
	fullFilePath := currPath + "/" + userId + "/" + group + "/" + currDate
	fullFileName := fullFilePath + "/" + fileName
	url := strings.Replace(fullFileName, "/", "|", -1)
	mkdirErr := os.MkdirAll(fullFilePath, 0777)

	if mkdirErr != nil {
		return url, mkdirErr
	}

	currFile, openFileErr := os.OpenFile(fullFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)

	if openFileErr != nil {
		return url, openFileErr
	}

	defer currFile.Close()

	length, writeErr := currFile.Write(source)
	image := datafile.ImageInfo{Actions: "Add", ServiceId: serviceId, FileId: fileId, Path: path, UserId: userId, Group: group, Url: url}
	datafile.SaveRecord(image.ToRecord())

	if writeErr != nil {
		return url, writeErr
	}

	if length != len(source) {
		return url, errors.New("Write file size unequal with read!!")
	}

	return url, nil
}

func (fileService *FileService) ReadFile(fileName string) ([]byte, error) {

	fullPath := strings.Replace(fileName, "|", "/", -1)

	file, openErr := os.OpenFile(fullPath, os.O_RDONLY, 0666)

	defer file.Close()

	if openErr != nil {
		return nil, openErr
	}

	buffer, readFileErr := ioutil.ReadAll(file)

	if readFileErr != nil {
		return nil, readFileErr
	}

	return buffer, nil
}

func (fileService *FileService) ReplaceFile(path string, serviceId int64, userId string, group string, url string, source []byte) (string, error) {
	log4go.Debug("url = %v", url)
	fullFileName := strings.Replace(url, "|", "/", -1)
	pathSplitIndex := strings.LastIndex(fullFileName, "/")
	//	fullFilePath := fullFileName[0:pathSplitIndex]
	log4go.Debug("fullFileName = %v", fullFileName)
	log4go.Debug("pathSplitIndex = %v", pathSplitIndex)
	fileIdMeta := fullFileName[pathSplitIndex+1:]
	log4go.Debug("fileIdMeta = %v", fileIdMeta)
	fileId, fileIdParseErr := strconv.ParseInt(fileIdMeta, 10, 64)

	if fileIdParseErr != nil {
		return url, fileIdParseErr
	}

	tempFileName := fullFileName + ".tmp"
	os.Rename(fullFileName, tempFileName)
	currFile, openFileErr := os.OpenFile(fullFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)

	if openFileErr != nil {
		os.Rename(tempFileName, fullFileName)
		return url, openFileErr
	}

	defer currFile.Close()

	length, writeErr := currFile.Write(source)
	image := datafile.ImageInfo{Actions: "Update", ServiceId: serviceId, FileId: fileId, Path: path, UserId: userId, Group: group, Url: url}
	datafile.SaveRecord(image.ToRecord())

	if writeErr != nil {
		os.Rename(tempFileName, fullFileName)
		return url, writeErr
	}

	if length != len(source) {
		os.Rename(tempFileName, fullFileName)
		return url, errors.New("Write file size unequal with read!!")
	}

	os.Remove(tempFileName)

	return url, nil
}

func (fileService *FileService) DeleteFile(path string, serviceId int64, userId string, group string, url string) error {
	fullFileName := strings.Replace(url, "|", "/", -1)
	removeErr := os.Remove(fullFileName)

	if removeErr != nil {
		return removeErr
	}

	pathSplitIndex := strings.LastIndex(fullFileName, "/")
	//	fullFilePath := fullFileName[0:pathSplitIndex]
	log4go.Debug("fullFileName = %v", fullFileName)
	log4go.Debug("pathSplitIndex = %v", pathSplitIndex)
	fileIdMeta := fullFileName[pathSplitIndex+1:]

	fileId, fileIdParseErr := strconv.ParseInt(fileIdMeta, 10, 64)

	if fileIdParseErr != nil {
		return fileIdParseErr
	}

	image := datafile.ImageInfo{Actions: "Delete", ServiceId: serviceId, FileId: fileId, Path: path, UserId: userId, Group: group, Url: url}
	datafile.SaveRecord(image.ToRecord())
	return nil
}
