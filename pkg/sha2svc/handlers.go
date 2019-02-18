package sha2svc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"goji.io/pat"
)

func writeErrorResponse(w http.ResponseWriter, code int, msg string) {
	writeOneArgResponse(w, code, "message", msg)
}

func writeOneArgResponse(w http.ResponseWriter, code int, param string, val string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{%q: %q}", param, val)
}

func writeResponse(w http.ResponseWriter, code int, json []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

// ProcessHashRequest handler for incoming request for hashing given value.
// Incoming request example {"payload": "test", "hash_rounds_cnt": 1} (id=ee89026a6c5603c51b4504d218ac60f6874b7750)
// Constraints: len(payload) > 0, hash_rounds_cnt > 0
func ProcessHashRequest(repo *ServiceRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Decode and check request
		var req RequestData
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			writeErrorResponse(w, http.StatusBadRequest, "Incorrect request")
			return
		}

		if req.HashRoundsCnt < 1 || len(req.Payload) < 1 {
			writeErrorResponse(w, http.StatusPreconditionFailed, "payload and hash_rounds_cnt should be greater than zero")
			return
		}

		// Generate uniq id for the job
		jobid := repo.Idgen.GenerateID(req)
		log.Printf("incoming request #%s for hash from %s", jobid, r.RemoteAddr)

		// Save request to storage
		err = repo.Saver.SaveHashRequest(jobid, req)
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, "Error while saving request")
			log.Printf("error saving request #%s for hash", jobid)
			return
		}

		// Write response with job id
		writeOneArgResponse(w, http.StatusCreated, "id", jobid)

		// Calculate and save hash. Also put coroutine to waitgroup to complete on exit
		go func() {
			repo.WaitGroup.Add(1)
			defer repo.WaitGroup.Done()
			hash := repo.Hasher.Hash(req.Payload, req.HashRoundsCnt)
			repo.Saver.SaveHashResult(jobid, hash)
		}()
	}
}

// ProcessResultRequest handler for request for result of hashing
func ProcessResultRequest(repo *ServiceRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		jobid := pat.Param(r, "id")
		log.Printf("incoming request for result of job #%s from %s", jobid, r.RemoteAddr)

		// fetch result from database
		resp, err := repo.Fetcher.GetHashResult(jobid)
		if err != nil {
			log.Printf("error fetching result for request #%s", jobid)
			writeErrorResponse(w, http.StatusNotFound, jobid+" not found")
			return
		}

		// write response
		respBody, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			log.Printf("error generation json for request #%s", jobid)
			writeErrorResponse(w, http.StatusInternalServerError, jobid+" incorrect record")
			return
		}

		writeResponse(w, http.StatusOK, respBody)
	}
}
