package lock

import (
	"proxy-fileserver/enums"
	"sync"
)

type MapLock = map[string]*sync.RWMutex

var _MapLock MapLock

func InitMapLock() {
	_MapLock = make(MapLock)
}

func RLockWithKey(key string) error {
	mutex, ok := _MapLock[key]
	if ok {
		mutex.RLock()
		return nil
	}
	return enums.ErrMutexNotFound
}

func RUnLockWithKey(key string) error {
	mutex, ok := _MapLock[key]
	if ok {
		mutex.RUnlock()
		return nil
	}
	return enums.ErrMutexNotFound
}

func WLockWithKey(key string) error {
	mutex, ok := _MapLock[key]
	if ok {
		mutex.Lock()
		return nil
	}
	return enums.ErrMutexNotFound
}

func WUnLockWithKey(key string) error {
	mutex, ok := _MapLock[key]
	if ok {
		mutex.Unlock()
		return nil
	}
	return enums.ErrMutexNotFound
}

func AddLock(key string) error {
	_, ok := _MapLock[key]
	if ok {
		return enums.ErrMutexExisted
	}
	_MapLock[key] = &sync.RWMutex{
	}
	return nil
}
