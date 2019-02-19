package sha2svc

import (
	"bytes"
	"net/http"
	. "net/http"
	"testing"
	"time"
)

func createTestRepo() (repo ServiceRepository) {
	repo.Idgen = &SmartHasher{}
	InitMemorySimpleStorageSvc(&repo)
	repo.Hasher = &Sha2Hasher{Version: Sha512}
	return
}

type brokenHasher struct{}

func (brokenHasher) Hash(payload string, cnt uint8) string {
	panic("brake process for test purposes")
}

type responseWriterMock struct {
	status int
	h      Header
}

func (m *responseWriterMock) Header() Header {
	return m.h
}
func (m *responseWriterMock) Write([]byte) (n int, e error) {
	return
}
func (m *responseWriterMock) WriteHeader(statusCode int) {
	m.status = statusCode
}

func createPostRequest(body string) (req *Request) {
	var byteReq bytes.Buffer
	byteReq.WriteString(body)
	req, _ = http.NewRequest("POST", "", &byteReq)
	return
}

func createGetRequest(id string) (req *Request) {
	req, _ = http.NewRequest("GET", "localhost/gvaduha/api/v1/hash/"+id, nil)
	return
}

const (
	correctReq     = `{"payload": "test", "hash_rounds_cnt": 1}`
	incorrectReq   = `{"payload": "test", "hash_rounds_cnt": "NOTNUM"}`
	zeropayloadReq = `{"payload": "", "hash_rounds_cnt": 1}`
	zerocntReq     = `{"payload": "test", "hash_rounds_cnt": 0}`
)

func TestProcessHashRequestErrors(t *testing.T) {

	//IMPORTANT: Tables ONE TIME ONLY! arg request body will be void after first use
	tables := []struct {
		arg *Request
		res int
	}{
		{createPostRequest(correctReq), StatusCreated},
		{createPostRequest(incorrectReq), StatusBadRequest},
		{createPostRequest(zeropayloadReq), StatusPreconditionFailed},
		{createPostRequest(zerocntReq), StatusPreconditionFailed},
	}

	repo := createTestRepo()
	handler := ProcessHashRequest(&repo)
	w := responseWriterMock{status: -1, h: make(Header)}

	for _, tab := range tables {
		handler(&w, tab.arg)
		if tab.res != w.status {
			t.Errorf("Status differ waiting for %d got %d", tab.res, w.status)
		}
		w.status = -1
	}
}

func TestProcessHashRequestSavedStatuses(t *testing.T) {
	repo := createTestRepo()
	handler := ProcessHashRequest(&repo)
	w := responseWriterMock{status: -1, h: make(Header)}

	hash := "ee89026a6c5603c51b4504d218ac60f6874b7750"

	// correct save
	handler(&w, createPostRequest(correctReq))
	res, err := repo.Fetcher.GetHashResult(hash)
	if err != nil || res.Status != InProgress {
		t.Error("Status of request should be in progress")
	}

	// correct save and process
	handler(&w, createPostRequest(correctReq))
	time.Sleep(time.Second)
	res, err = repo.Fetcher.GetHashResult(hash)
	if err != nil || res.Status != Finished {
		t.Error("Status of request should be finished")
	}
}

/*
func TestProcessResultRequest(t *testing.T) {
	repo := createTestRepo()
	createhandler := ProcessHashRequest(&repo)
	testhandler := ProcessResultRequest(&repo)
	w := responseWriterMock{status: -1, h: make(Header)}

	hash := "ee89026a6c5603c51b4504d218ac60f6874b7750"

	// create request
	createhandler(&w, createPostRequest(correctReq))

	// test cases
	tables := []struct {
		arg *Request
		res int
	}{
		{createGetRequest(hash), StatusCreated},
		{createGetRequest("0000"), StatusNotFound},
	}

	pat := struct{ s int }{s: 1}
	_ = pat

	for _, tab := range tables {
		testhandler(&w, tab.arg)
		if w.status != StatusOK {
			t.Errorf("wrong status %d", tab.res)
		}
	}
}
*/
