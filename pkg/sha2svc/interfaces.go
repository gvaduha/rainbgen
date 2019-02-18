package sha2svc

import (
	"io"
	"log"
	"sync"
)

// RequestData incoming hash request
type RequestData struct {
	// Payload data to hash
	Payload string `json:"payload"`
	// HashRoundsCnt number of times to apply hash function
	HashRoundsCnt uint8 `json:"hash_rounds_cnt"`
}

// ResponseData response to hashing request
type ResponseData struct {
	// ID uniq identifier of request
	ID string `json:"id"`
	// Payload data to hash
	Payload string `json:"payload"`
	// HashRoundsCnt number of times to apply hash function
	HashRoundsCnt uint8 `json:"hash_rounds_cnt"`
	// Status of processing request
	Status ReqStatus `json:"status"`
	// Hash sha2 hash of payload applied HashRoundsCnt number of times
	Hash string `json:"hash"`
}

// IDGenerator interface for generation of uniq id from given data
type IDGenerator interface {
	// GenerateId generates uniq id from given data
	GenerateID(req RequestData) string
}

// StorageSaver interface for saving hash requests and processing results
type StorageSaver interface {
	// SaveHashRequest save incoming request
	SaveHashRequest(id string, req RequestData) error
	// SaveHashResult save result of hashing given data. In case of failure status is failed an no hash supplied
	SaveHashResult(id string, hash string) error

	io.Closer
}

// StorageFetcher interface for retrieving information about hash requests
type StorageFetcher interface {
	// GetHashResult retrieve proc
	GetHashResult(id string) (ResponseData, error)

	io.Closer
}

// Hasher interface for doing main job of hashing
type Hasher interface {
	// Hash creates and returns sha2 hash cnt number of times from given payload
	Hash(payload string, cnt uint8) string
}

// ServiceRepository collection of "service" objects to perform job
type ServiceRepository struct {
	Idgen   IDGenerator
	Saver   StorageSaver
	Fetcher StorageFetcher
	Hasher  Hasher

	WaitGroup sync.WaitGroup

	io.Closer
}

// Close collection of objects with external resources in CleanUpObjects member
func (s *ServiceRepository) Close() (err error) {
	log.Println("service repository cleanup")
	s.Saver.Close()
	s.Fetcher.Close()
	return
}
