package handler

import (
	"bytes"
	"error"
	"errors"
	"filecache"
	"imageactions"
	"io/ioutil"
	"net/url"
	"redisclusteradapter"
	"service/business"
	"strings"
	"syscommon"
	"sysconf"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/gitstliu/log4go"
)

var redisCommonKey = "CCDN:Image:"

type FileFacade struct {
	FileService business.FileService
}

func (fileAdapter *FileFacade) SaveFile(w rest.ResponseWriter, r *rest.Request) {

	response := syscommon.CommonResponse{Code: responsecode.Success, Message: "Success"}

	userName, group, decodeUserInfoErr := decodeUserInfo(r)

	if decodeUserInfoErr != nil {
		response.Code = responsecode.Fail
		//		response.Message = decodeUserInfoErr.Error()
		response.Message = "User is invalid"
		log4go.Error(decodeUserInfoErr)
		w.WriteJson(&response)
		return
	}

	buffer, err := ioutil.ReadAll(r.Body)

	log4go.Debug("UpLoad File Size = %d bit", len(buffer))

	if err != nil {
		response.Code = responsecode.Fail
		response.Message = err.Error()
		log4go.Error(err)
		w.WriteJson(&response)
		return
	}

	if buffer == nil || len(buffer) == 0 {
		response.Code = responsecode.Fail
		response.Message = "Buffer size is 0!"
		log4go.Error("Buffer size is 0!")
		w.WriteJson(&response)
		return
	}

	url, saveErr := fileAdapter.FileService.SaveFile(sysconf.GetConfigure().FileContentPath, sysconf.GetConfigure().ServiceId, userName, group, buffer)
	if saveErr != nil {
		response.Code = responsecode.Fail
		response.Message = saveErr.Error()
		log4go.Error(saveErr)
		w.WriteJson(&response)
		return
	}

	response.Result = url

	w.WriteJson(&response)

}

func (fileAdapter *FileFacade) ReadFile(w rest.ResponseWriter, r *rest.Request) {

	response := syscommon.CommonResponse{Code: responsecode.Success, Message: "Success"}

	fileName := r.PathParam("file_name")

	width := 0
	height := 0

	log4go.Debug("r.URL.Query() = %v", r.URL.Query())
	widthMeta, widthOK := r.URL.Query()["width"]
	heightMeta, heighOK := r.URL.Query()["height"]

	if widthOK {
		log4go.Debug("widthMeta[0] = %v", widthMeta[0])
		width = syscommon.StringToInt32(widthMeta[0], 0)
	}

	if heighOK {
		log4go.Debug("heightMeta[0] = %v", heightMeta[0])
		height = syscommon.StringToInt32(heightMeta[0], 0)
	}

	isCustome := widthOK || heighOK

	if isCustome {
		if width <= 0 {
			response.Code = responsecode.Fail
			response.Message = "custome width should be larger then 0"
			w.WriteJson(&response)
			return
		}

		if height <= 0 {
			response.Code = responsecode.Fail
			response.Message = "custome heigh should be larger then 0"
			w.WriteJson(&response)
			return
		}
	}

	decodeFileName, decodeErr := url.QueryUnescape(fileName)

	if decodeErr != nil {
		response.Code = responsecode.Fail
		response.Message = decodeErr.Error()
		log4go.Error(decodeErr)
		w.WriteJson(&response)
		return
	}

	currKey := strings.Join([]string{decodeFileName, syscommon.Int32ToString(width), syscommon.Int32ToString(height)}, ".")
	log4go.Debug("currKey = %v", currKey)

	memCache := filecache.CurrFileCache.Read(currKey)

	if memCache != nil {
		log4go.Debug("Mem Cached !!!")
		w.WriteBytes(memCache.Binary)
		return
	}

	log4go.Debug("Mem UNCached !!!")
	redisKey := getRedisKey(currKey)
	cacheString, readCacheErr := redisclusteradapter.GetAdapter().GET(redisKey)
	cacheBuffer := []byte(cacheString)
	redisCached := true
	if readCacheErr != nil {
		log4go.Debug(readCacheErr)
		log4go.Debug("Key is " + redisKey)
		redisCached = false
	} else if len(cacheBuffer) == 0 {
		redisCached = false
	}

	if redisCached {
		log4go.Debug("Redis Cached !!!")
		writeErr := filecache.CurrFileCache.Write(decodeFileName, currKey, &filecache.FileMeta{Size: len(cacheBuffer), Binary: cacheBuffer})
		if writeErr != nil {
			log4go.Error(writeErr)
		}
		go redisclusteradapter.GetAdapter().EXPIRE(redisKey, sysconf.GetConfigure().RedisCacheTimeOut)
		w.WriteBytes(cacheBuffer)
		return
	}
	log4go.Debug("Redis Uncached !!!")
	log4go.Debug("decodeFileName = %v", decodeFileName)
	buffer, err := fileAdapter.FileService.ReadFile(decodeFileName)

	if err != nil {
		response.Code = responsecode.Fail
		response.Message = err.Error()
		log4go.Error(err)
		w.WriteJson(&response)
		return
	}

	log4go.Debug("sysconf.GetConfigure().CatchTimeOut = %v", sysconf.GetConfigure().RedisCacheTimeOut)
	if !isCustome {
		go setBufferToRedisCache(decodeFileName, redisKey, buffer)
		writeErr := filecache.CurrFileCache.Write(decodeFileName, currKey, &filecache.FileMeta{Size: len(buffer), Binary: buffer})
		if writeErr != nil {
			log4go.Error(writeErr)
		}
		w.WriteBytes(buffer)
		return
	}

	imageReader := bytes.NewReader(buffer)
	log4go.Debug("len(buffer) = %v", len(buffer))

	imageBuffer, imageErr := imageactions.ImageToSpecificSize(imageReader, width, height)

	if imageErr != nil {
		response.Code = responsecode.Fail
		response.Message = imageErr.Error()
		log4go.Error(imageErr)
		w.WriteJson(&response)
		return
	}
	go setBufferToRedisCache(decodeFileName, redisKey, buffer)

	writeErr := filecache.CurrFileCache.Write(decodeFileName, currKey, &filecache.FileMeta{Size: len(imageBuffer), Binary: imageBuffer})
	if writeErr != nil {
		log4go.Error(writeErr)
	}
	w.WriteBytes(imageBuffer)
	return
}

func (fileAdapter *FileFacade) ReplaceFile(w rest.ResponseWriter, r *rest.Request) {

	response := syscommon.CommonResponse{Code: responsecode.Success, Message: "Success"}

	userName, group, decodeUserInfoErr := decodeUserInfo(r)

	if decodeUserInfoErr != nil {
		response.Code = responsecode.Fail
		//		response.Message = decodeUserInfoErr.Error()
		response.Message = "User is invalid"
		log4go.Error(decodeUserInfoErr)
		w.WriteJson(&response)
		return
	}

	fileName := r.PathParam("file_name")

	decodeFileName, decodeErr := url.QueryUnescape(fileName)

	if decodeErr != nil {
		response.Code = responsecode.Fail
		response.Message = decodeErr.Error()
		log4go.Error(decodeErr)
		w.WriteJson(&response)
		return
	}

	buffer, err := ioutil.ReadAll(r.Body)

	log4go.Debug("Upload File Size = %d bit", len(buffer))

	if err != nil {
		response.Code = responsecode.Fail
		response.Message = err.Error()
		log4go.Error(err)
		w.WriteJson(&response)
		return
	}

	if buffer == nil || len(buffer) == 0 {
		response.Code = responsecode.Fail
		response.Message = "Buffer size is 0!"
		log4go.Error("Buffer size is 0!")
		w.WriteJson(&response)
		return
	}

	url, replaceErr := fileAdapter.FileService.ReplaceFile(sysconf.GetConfigure().FileContentPath, sysconf.GetConfigure().ServiceId, userName, group, decodeFileName, buffer)
	if replaceErr != nil {
		response.Code = responsecode.Fail
		response.Message = replaceErr.Error()
		log4go.Error(replaceErr)
		w.WriteJson(&response)
		return
	}

	go cleanSpecificRedisCache(decodeFileName)
	go filecache.CurrFileCache.CleanSepecificCache(decodeFileName)

	response.Result = url

	w.WriteJson(&response)
}

func (fileAdapter *FileFacade) DeleteFile(w rest.ResponseWriter, r *rest.Request) {

	response := syscommon.CommonResponse{Code: responsecode.Success, Message: "Success"}

	userName, group, decodeUserInfoErr := decodeUserInfo(r)

	if decodeUserInfoErr != nil {
		response.Code = responsecode.Fail
		//		response.Message = decodeUserInfoErr.Error()
		response.Message = "User is invalid"
		log4go.Error(decodeUserInfoErr)
		w.WriteJson(&response)
		return
	}

	fileName := r.PathParam("file_name")

	decodeFileName, decodeErr := url.QueryUnescape(fileName)

	if decodeErr != nil {
		response.Code = responsecode.Fail
		response.Message = decodeErr.Error()
		log4go.Error(decodeErr)
		w.WriteJson(&response)
		return
	}

	deleteErr := fileAdapter.FileService.DeleteFile(sysconf.GetConfigure().FileContentPath, sysconf.GetConfigure().ServiceId, userName, group, decodeFileName)
	if deleteErr != nil {
		response.Code = responsecode.Fail
		response.Message = deleteErr.Error()
		log4go.Error(deleteErr)
		w.WriteJson(&response)
		return
	}

	go cleanSpecificRedisCache(decodeFileName)
	go filecache.CurrFileCache.CleanSepecificCache(decodeFileName)

	response.Result = decodeFileName

	w.WriteJson(&response)
}

func (fileAdapter *FileFacade) CleanCache(w rest.ResponseWriter, r *rest.Request) {

	response := syscommon.CommonResponse{Code: responsecode.Success, Message: "Success"}

	_, _, decodeUserInfoErr := decodeUserInfo(r)

	if decodeUserInfoErr != nil {
		response.Code = responsecode.Fail
		//		response.Message = decodeUserInfoErr.Error()
		response.Message = "User is invalid"
		log4go.Error(decodeUserInfoErr)
		w.WriteJson(&response)
		return
	}

	fileName := r.PathParam("file_name")

	decodeFileName, decodeErr := url.QueryUnescape(fileName)

	if decodeErr != nil {
		response.Code = responsecode.Fail
		response.Message = decodeErr.Error()
		log4go.Error(decodeErr)
		w.WriteJson(&response)
		return
	}

	go cleanSpecificRedisCache(decodeFileName)
	go filecache.CurrFileCache.CleanSepecificCache(decodeFileName)

	response.Result = decodeFileName

	w.WriteJson(&response)
}

func decodeUserInfo(r *rest.Request) (string, string, error) {
	loginTokenMeta, loginTokenParamExist := r.Header["Login-Token"]
	groupName := r.PathParam("group_name")

	log4go.Debug("groupName %v", groupName)
	if groupName == "" {
		groupName = "YHYS"
	}

	log4go.Debug("r.Header %v", r.Header)

	if !loginTokenParamExist {
		return "", "", errors.New("Upload file must have a token")
	}
	if len(loginTokenMeta) > 0 {
		userInfo, getUserInfoErr := GetUserInfoByToken(loginTokenMeta[0])
		if getUserInfoErr != nil {
			log4go.Error(getUserInfoErr)
			return "", "", getUserInfoErr
		}
		//		if loginTokenMeta[0] != "8080eafa-d00c-4e54-a51d-03d25119186a" {
		//			return "", "", errors.New("Token Check Failed!!!")
		//		}
		currGroup, groupExist := userInfo.Groups[groupName]
		if !groupExist {
			return "", "", errors.New("Group is invalid!!!")
		}
		//		userInfo.Username
		//		userName = "YHYS"
		//		group = "YHTS"
		return userInfo.Username, currGroup.GroupName, nil

	} else {
		return "", "", errors.New("Upload file must have a token")
	}
}

func setBufferToRedisCache(redisKey string, redisKeyExtend string, buffer []byte) error {
	setRedisKey := redisclusteradapter.RedisPipelineCommand{}
	setRedisKey.Key = redisKey
	setRedisKey.CommandName = "SADD"
	setRedisKey.Args = []interface{}{redisKeyExtend}

	setRedisKeyExtend := redisclusteradapter.RedisPipelineCommand{}
	setRedisKeyExtend.Key = redisKeyExtend
	setRedisKeyExtend.CommandName = "SET"
	setRedisKeyExtend.Args = []interface{}{buffer, "EX", sysconf.GetConfigure().RedisCacheTimeOut}

	piplineResults, piplineErrors := redisclusteradapter.GetAdapter().SendPipelineCommands([]redisclusteradapter.RedisPipelineCommand{setRedisKey, setRedisKeyExtend})

	isErr := false
	for index, currErr := range piplineErrors {
		if currErr != nil {
			log4go.Error(currErr)
			log4go.Error("Pipline Command result = %v", piplineResults[index])
			isErr = true
		}
	}

	if isErr {
		return errors.New("Redis Error, Please Check The Error Log!!!")
	}

	return nil
}

func cleanSpecificRedisCache(redisKey string) error {
	smembersResult, smembersErr := redisclusteradapter.GetAdapter().SMEMBERS(redisKey)

	if smembersErr != nil {
		//		log4go.Error(smembersErr)
		return smembersErr
	}

	commands := []redisclusteradapter.RedisPipelineCommand{}
	for _, currResult := range smembersResult {
		currCommand := redisclusteradapter.RedisPipelineCommand{}
		currCommand.Key = currResult
		currCommand.CommandName = "DEL"
		commands = append(commands, currCommand)
	}

	piplineResults, piplineErrors := redisclusteradapter.GetAdapter().SendPipelineCommands(commands)

	isErr := false
	log4go.Debug("len(piplineResults) = %v", len(piplineResults))
	log4go.Debug("len(piplineErrors) = %v", len(piplineErrors))
	for index, currErr := range piplineErrors {
		if currErr != nil {
			log4go.Error(currErr)
			log4go.Error("Pipline Command result = %v", piplineResults[index])
			isErr = true
		}
	}

	if isErr {
		return errors.New("Redis Error, Please Check The Error Log!!!")
	}

	return nil
}

func getRedisKey(key string) string {
	return strings.Join([]string{redisCommonKey, key}, "")
}
