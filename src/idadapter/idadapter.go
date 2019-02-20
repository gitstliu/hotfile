package idadapter

import (
	"errors"
	"sync"

	"sysconf"

	"github.com/gitstliu/go-id-worker"
)

var idWorkerMap = &sync.Map{}

func CreateIdWorker(key string) {
	currWoker := &idworker.IdWorker{}
	currWoker.InitIdWorker(sysconf.GetConfigure().ServiceId, sysconf.GetConfigure().DatacenterId)
	idWorkerMap.Store(key, currWoker)
}

//func IsIdWorderExits(key string) bool {
//	_, isExits := idWorkerMap.Load(key)
//	return isExits
//}

func GetNewId(key string) (int64, error) {
	currWoker, isExits := idWorkerMap.Load(key)

	if !isExits {
		CreateIdWorker(key)
		currWoker, isExits = idWorkerMap.Load(key)

		if !isExits {
			return -1, errors.New("Could not create Id Worker")
		}
	}

	return currWoker.(*idworker.IdWorker).NextId()
}
