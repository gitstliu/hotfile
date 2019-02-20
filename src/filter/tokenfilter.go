package filter

import (
	"error"
	"fmt"
	"net/http"
	"redisadapter"
	"syscommon"

	"github.com/gitstliu/log4go"
)

func CheckToken(w http.ResponseWriter, r *http.Request) (bool, interface{}) {
	fmt.Println("CheckToken")
	response := syscommon.CommonResponse{Code: responsecode.FilerNotAllowed, Message: "FilerNotAllowed"}
	token := r.Header["Authtoken"]
	fmt.Println(r.Header)

	if token == nil || len(token) == 0 {
		response.Message = "Token missed!!"
		return false, response
	}

	userIDResult, userIDErr := redisadapter.GetAdapter().Get(redisadapter.UserLoginInfo_TokenToUser + token[0])

	if userIDErr != nil {
		log4go.Error(userIDErr)
		response.Message = "Token Check Failed!!"
		return false, response
	}

	if userIDResult == "" {
		log4go.Debug("User info should not be nil or empty!! For token " + token[0])
		response.Message = "Token Check Failed!!"
		return false, response
	}

	userInfoResult, userInfoErr := redisadapter.GetAdapter().Get(redisadapter.UserLoginInfo_UserNameToUserInfo + userIDResult)

	if userInfoErr != nil {
		log4go.Error(userInfoErr)
		response.Message = "Token Check Failed!!"
		return false, response
	}

	r.Header["IntraTokenUserInfo"] = []string{userInfoResult}
	return true, ""
}

func CheckTokenAfter(w http.ResponseWriter, r *http.Request) (bool, interface{}) {
	fmt.Println("CheckToken")
	response := syscommon.CommonResponse{Code: responsecode.FilerNotAllowed, Message: "FilerNotAllowed"}
	token := r.Header["Authtoken"]
	fmt.Println(r.Header)

	if token == nil || len(token) == 0 {
		response.Message = "Token missed!!"
		return false, response
	}

	redisResult, redisErr := redisadapter.GetAdapter().Get(redisadapter.UserLoginInfo_TokenToUser + token[0])
	if redisErr != nil {
		log4go.Error(redisErr)
		response.Message = "Token Check Failed!!"
		return false, response
	}

	if redisResult == "" {
		log4go.Debug("User info should not be nil or empty!! For token " + token[0])
		response.Message = "Token Check Failed!!"
		return false, response
	}

	r.Header["IntraTokenUserInfo"] = []string{redisResult}
	return true, ""
}
