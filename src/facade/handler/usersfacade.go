package handler

import (
	"facade/vo"
	"redisclusteradapter"
	"syscommon"

	"github.com/gitstliu/log4go"
)

func GetUserInfoByToken(token string) (*vo.UserInfo, error) {
	log4go.Debug(syscommon.UserLoginInfo_TokenToUser + token)
	userInfoJson, getUserInfoErr := redisclusteradapter.GetAdapter().GET(syscommon.UserLoginInfo_TokenToUser + token)
	result := &vo.UserInfo{}
	if getUserInfoErr != nil {
		log4go.Error(getUserInfoErr)
		return nil, getUserInfoErr
	}
	jsonToUserInfoErr := syscommon.JsonToObject(userInfoJson, result)

	if jsonToUserInfoErr != nil {
		log4go.Error(jsonToUserInfoErr)
		return nil, jsonToUserInfoErr
	}

	return result, nil
}
