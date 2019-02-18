package sha2svc

import (
	"errors"
	"io"
	"log"
	"strconv"
	"sync"
)

type memorySimpleStorage struct {
	removeAfterGetResult bool
	store                map[string]ResponseData
	mtx                  sync.Mutex

	io.Closer
}

func (s *memorySimpleStorage) Close() (err error) {
	log.Println("cleanup memory storage")
	return
}

func (s *memorySimpleStorage) SaveHashRequest(id string, req RequestData) (err error) {
	var res ResponseData
	res.ID = id
	res.Payload = req.Payload
	res.HashRoundsCnt = req.HashRoundsCnt
	res.Status = InProgress
	s.mtx.Lock()
	s.store[id] = res
	s.mtx.Unlock()
	return
}

func (s *memorySimpleStorage) SaveHashResult(id string, hash string) (err error) {
	if res, ok := s.store[id]; ok {
		res.Hash = hash
		res.Status = Finished
		s.mtx.Lock()
		s.store[id] = res
		s.mtx.Unlock()
	} else {
		err = errors.New("value not exist")
	}
	return
}

func (s *memorySimpleStorage) GetHashResult(id string) (resp ResponseData, err error) {
	var ok bool
	if resp, ok = s.store[id]; !ok {
		err = errors.New("value not exist")
	} else {
		if s.removeAfterGetResult {
			s.mtx.Lock()
			delete(s.store, id)
			s.mtx.Unlock()
		}
	}
	return
}

// InitMemorySimpleStorageSvc prepare memory simple storage and assigns it to ServiceRepository fetcher and saver
func InitMemorySimpleStorageSvc(repo *ServiceRepository) {
	var storage memorySimpleStorage
	removeAfterFlag := GetEnvOrDefault("MEMORYSIMPLESTORAGE_REMOVE_AFTER_GETRESULT", "true")

	storage.removeAfterGetResult, _ = strconv.ParseBool(removeAfterFlag)
	storage.store = make(map[string]ResponseData)

	repo.Saver = &storage
	repo.Fetcher = &storage
}
