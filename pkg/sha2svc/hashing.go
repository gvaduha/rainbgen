package sha2svc

import (
	"crypto/sha1"
	"crypto/sha512"
	"fmt"
	"strconv"
	"time"
)

// sha1Hasher implements IDGenerator interface for generation of uniq id from given data
type sha1Hasher struct {
}

func (sha1Hasher) GenerateID(req RequestData) string {
	data := strconv.Itoa(int(req.HashRoundsCnt)) + req.Payload
	res := sha1.Sum([]byte(data))
	return fmt.Sprintf("%x", res[:sha1.Size])
}

// stringHasher implements IDGenerator interface for generation of uniq id from given data
type stringHasher struct {
}

func (stringHasher) GenerateID(req RequestData) string {
	return strconv.Itoa(int(req.HashRoundsCnt)) + req.Payload
}

// SmartHasher implements IDGenerator interface for generation of uniq id from given data
type SmartHasher struct {
}

// GenerateID create "uniq" hash based on data (ad-hoc "smart" implementation needs to be improved)
func (SmartHasher) GenerateID(req RequestData) string {
	var data string
	if len(req.Payload) < 4096 {
		data = strconv.Itoa(int(req.HashRoundsCnt)) + req.Payload
	} else {
		data = time.Now().Format(time.RFC3339) + strconv.Itoa(int(req.HashRoundsCnt)) + req.Payload
	}

	res := sha1.Sum([]byte(data))
	return fmt.Sprintf("%x", res[:sha1.Size])
}

// Sha2Version type of SHA2
type Sha2Version uint8

const (
	// Sha512 SHA-2 512
	Sha512 Sha2Version = iota
	// Sha384 SHA-2 384
	Sha384
	// Sha512_256 SHA-2 512/224
	Sha512_256
	// Sha512_224 SHA-2 512/224
	Sha512_224
)

// Sha2Hasher implements Hasher to make SHA-2 hashes
type Sha2Hasher struct {
	Version Sha2Version
}

// Hash compute SHA-2 hash of payload number of times
func (Sha2Hasher) Hash(payload string, cnt uint8) string {
	hashfun := sha512.Sum512
	data := []byte(payload)

	hash := hashfun(data)

	for i := 1; i < int(cnt); i++ {
		hash = hashfun(hash[:])
	}

	return fmt.Sprintf("%x", hash[:])
}
