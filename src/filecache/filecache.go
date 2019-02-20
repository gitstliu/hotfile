package filecache

import (
	"errors"
	"sync"

	"github.com/gitstliu/log4go"
)

type FileMeta struct {
	Size   int
	Binary []byte
}

type FileCache struct {
	maxSize        int
	oldSize        int
	newSize        int
	originalKeyMap *sync.Map
	oldMap         *sync.Map
	newMap         *sync.Map
	lock           *sync.RWMutex
}

var CurrFileCache = &FileCache{oldSize: 0, newSize: 0, oldMap: &sync.Map{}, newMap: &sync.Map{}, lock: &sync.RWMutex{}, originalKeyMap: &sync.Map{}}

func InitFileCache(maxSize int) {
	CurrFileCache.maxSize = maxSize
}

func (this *FileCache) Read(key string) *FileMeta {
	currOld := this.oldMap
	currNew := this.newMap

	result, isExits := currOld.Load(key)

	if !isExits {

		result, isExits = currNew.Load(key)

		if !isExits {
			return nil
		} else {
			return result.(*FileMeta)
		}

	} else {
		resultMeta := result.(*FileMeta)
		currOld.Delete(key)
		currNew.Store(key, result)
		this.lock.Lock()
		this.oldSize -= resultMeta.Size
		this.newSize += resultMeta.Size
		this.lock.Unlock()
		return resultMeta
	}
}

func (this *FileCache) Write(originalKey string, key string, file *FileMeta) error {

	if this.maxSize < this.newSize+this.oldSize+file.Size {
		//		Flush()
		return errors.New("Cache Out of Mem!!!!!")
	}

	originalKeySet, isOriginalKeySetExits := this.originalKeyMap.Load(originalKey)

	log4go.Debug("originalKey = %v, key = %v", originalKey, key)
	if isOriginalKeySetExits {
		originalKeySet.(*sync.Map).Store(key, "")
	} else {
		originalKeySet := &sync.Map{}
		originalKeySet.Store(key, "")
		this.originalKeyMap.Store(originalKey, originalKeySet)
	}

	currNew := this.newMap

	log4go.Debug("Write Mem cache %v", key)
	_, isExits := currNew.Load(key)

	if !isExits {
		log4go.Debug("Write Mem cache OK")
		currNew.Store(key, file)
		this.lock.Lock()
		this.newSize += file.Size
		this.lock.Unlock()
	}

	return nil
}

func (this *FileCache) Flush(threshold int) {

	if this.maxSize < this.newSize+this.oldSize+threshold {
		log4go.Debug("Threshold is not broken!!! Nothin will be changed!!")
		return
	}

	if this.newSize == 0 || this.oldSize == 0 {
		log4go.Debug("No access change in this flush!!! Nothing will be changed!!")
		return
	}

	this.lock.Lock()

	tempNewSize := this.newSize
	tempNew := this.newMap
	this.newSize = 0
	this.oldSize = tempNewSize
	this.newMap = &sync.Map{}
	this.oldMap = tempNew

	this.lock.Unlock()
}

func (this *FileCache) CleanSepecificCache(originalKey string) {
	originalKeySet, isOriginalKeySetExits := this.originalKeyMap.Load(originalKey)

	if !isOriginalKeySetExits {
		return
	}

	originalKeySet.(*sync.Map).Range(this.remove)
	this.originalKeyMap.Delete(originalKey)
}

func (this *FileCache) remove(key, value interface{}) bool {

	log4go.Debug("Key = %v , value= %v", key, value)
	oldValue, oldValueExits := this.oldMap.Load(key)

	if oldValueExits {
		log4go.Debug("old remove")
		this.lock.Lock()
		this.oldSize -= oldValue.(*FileMeta).Size
		this.lock.Unlock()
		this.oldMap.Delete(key)
	} else {
		log4go.Debug("new remove")
		newValue, newValueExits := this.newMap.Load(key)
		if newValueExits {
			this.lock.Lock()
			this.newSize -= newValue.(*FileMeta).Size
			this.lock.Unlock()
			this.newMap.Delete(key)
		}
	}

	return true

}
