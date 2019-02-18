package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"goji.io"
	"goji.io/pat"

	svc "github.com/gvaduha/rainbgen/pkg/sha2svc"
)

func main() {
	var port = svc.GetEnvOrDefault("SVCPORT", "8800")
	var endpoint = svc.GetEnvOrDefault("SVCENDPOINT", "/gvaduha/api/v1/hash")

	svcRepo := createSvcRepo()
	defer svcRepo.Close()

	// create ctrl-c channel notifier
	stopchan := make(chan os.Signal, 1)
	signal.Notify(stopchan, os.Interrupt)
	// create channel to wait for http server to complete first
	httpsrvdone := make(chan bool, 1)

	// start http server
	mux := goji.NewMux()
	mux.HandleFunc(pat.Post(endpoint), svc.ProcessHashRequest(&svcRepo))
	mux.HandleFunc(pat.Get(endpoint+"/:id"), svc.ProcessResultRequest(&svcRepo))

	log.Printf("rainbgen service @%s%s is ready to accept connections\n", port, endpoint)

	hs := &http.Server{Addr: "localhost:" + port, Handler: mux}

	go func() {
		err := hs.ListenAndServe()
		if err != nil {
			log.Println(err)
		}
		log.Println("http server stopped serving requests")
		close(httpsrvdone)
	}()

	// waiting for stop signal and then to completion of http listeners
	<-stopchan
	log.Println("ctrl-c signal received")
	hs.Shutdown(context.Background())
	<-httpsrvdone

	// wait for worker coroutines (request for hashing) to complete
	log.Println("waiting for hashing tasks to finish")
	svcRepo.WaitGroup.Wait()
}

func setCtrlCInterrupt(f func()) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		f() // do cleanup...
	}()
}

func createSvcRepo() (repo svc.ServiceRepository) {
	repo.Idgen = &svc.SmartHasher{}
	svc.InitMongoStorageSvc(&repo)
	//svc.InitMemorySimpleStorageSvc(&repo)
	repo.Hasher = &svc.Sha2Hasher{Version: svc.Sha512}
	return
}
