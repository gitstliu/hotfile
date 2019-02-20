package datafile

import (
	"os"
	"strconv"
	"strings"
	//	"strings"
	"idadapter"
	"time"
)

type ImageInfo struct {
	Actions   string
	ServiceId int64
	FileId    int64
	Path      string
	UserId    string
	Group     string
	Url       string
}

func (this *ImageInfo) ToRecord() string {
	return "<ImageInfo>" + this.Actions + ":" + strconv.FormatInt(this.ServiceId, 10) + ":" + strconv.FormatInt(this.FileId, 10) + ":" + this.Path + ":" + this.UserId + ":" + this.Group + ":" + this.Url
}

func GetFileName(userId string, group string) (string, int64, error) {
	newId, newIdErr := idadapter.GetNewId(strings.Join([]string{userId, group}, ":"))

	if newIdErr != nil {
		return "", newId, newIdErr
	}

	fileName := strconv.FormatInt(newId, 10)
	return fileName, newId, nil
}

func SaveRecord(logValue string) error {
	year, month, day := time.Now().Date()
	dateValue := strconv.FormatInt(int64(year), 10) + strconv.FormatInt(int64(month), 10) + strconv.FormatInt(int64(day), 10)

	recordFile, openRecordFileErr := os.OpenFile("records/logicalnode-"+dateValue+".oc", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	defer recordFile.Close()

	if openRecordFileErr != nil {
		return openRecordFileErr
	}

	//	var buffer bytes.Buffer

	//	for _, nodeToRecord := range nodeToRecords {
	//		buffer.WriteString(nodeToRecord.ToRecord())
	//		buffer.WriteString("\n")
	//	}

	_, writeError := recordFile.WriteString(logValue)

	if writeError != nil {
		return writeError
	}

	return nil
}
