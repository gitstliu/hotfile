package syscommon

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gitstliu/log4go"
)

type TimeSpan struct {
	startNS int64
	endNS   int64
}

type NodeToRecord interface {
	ToRecord() string
}

func StringToInt32(meta string, defaultValue int) int {
	result, err := strconv.Atoi(meta)
	if err == nil {
		return result
	} else {
		return defaultValue
	}
}

func Int32ToString(meta int) string {
	return strconv.Itoa(meta)
}

func PanicHandler() {
	if r := recover(); r != nil {
		log4go.Error("Run time Error %v", r)
		//		fmt.Println(r)
		//		fmt.Printf("%T", r)
		panic(r)
	}
}

func InterfacesToStrings(src []interface{}) []string {
	result := []string{}
	for _, value := range src {
		result = append(result, value.(string))
	}
	return result
}

func (this *TimeSpan) Start() {
	this.startNS = time.Now().UnixNano()
}

func (this *TimeSpan) End() {
	this.endNS = time.Now().UnixNano()
}

func (this *TimeSpan) GetTimeSpanMS() float64 {

	return float64(this.endNS-this.startNS) / 1000000
}

func GetFilesWithFolder(folderName string) ([]string, error) {

	result := make([]string, 0, 1000)

	infos, readDirError := ioutil.ReadDir(folderName)

	if readDirError != nil {
		return nil, readDirError
	}

	for _, info := range infos {

		if !info.IsDir() {
			result = append(result, folderName+"/"+info.Name())
		}
	}

	return result, nil
}

func ObjectToJson(value interface{}) (string, error) {
	meta, err := json.Marshal(value)
	return string(meta), err
}

func ObjectsToJson(values []interface{}) ([]interface{}, error) {
	result := [](interface{}){}

	for _, currValue := range values {
		meta, err := ObjectToJson(currValue)

		if err != nil {
			return nil, err
		} else {
			result = append(result, meta)
		}
	}

	return result, nil
}

func DecodeGzipBytes(meta []byte) ([]byte, error) {
	b := bytes.Buffer{}
	b.Write(meta)
	r, _ := gzip.NewReader(&b)
	defer r.Close()
	datas, readErr := ioutil.ReadAll(r)

	if readErr != nil {
		return nil, readErr
	}

	return datas, nil
}

func EncodeGzipBytes(meta []byte) []byte {
	b := bytes.Buffer{}
	w := gzip.NewWriter(&b)
	defer w.Close()

	w.Write(meta)
	w.Flush()

	return b.Bytes()
}

func JsonToObject(meta string, result interface{}) error {
	return json.Unmarshal(StringToBytes(meta), result)
}

func StringToBytes(value string) []byte {
	return []byte(value)
}

func MetaToJsonContent(value string) string {
	return fmt.Sprintf("{%s}", value)
}

func Int64ToBytes(value int64) []byte {

	result := []byte{}

	buffer := bytes.NewBuffer(result)
	binary.Write(buffer, binary.BigEndian, &value)

	return buffer.Bytes()
}

func IntToBool(value int) bool {
	if value == 1 {
		return true
	} else {
		return false
	}
}

func IsGzipEncode(header http.Header) bool {

	log4go.Debug("header : %v", header)

	value := header["Content-Encoding"]

	return value != nil && strings.EqualFold(value[0], "gzip")
}
