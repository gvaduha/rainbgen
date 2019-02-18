package sha2svc

import (
	"io"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	dbname = "gvaduha"
	cname  = "rainbgen"
)

type dbRecord struct {
	ID            string
	Payload       string
	HashRoundsCnt uint8
	Status        string
	Hash          string
}

type mongoStorage struct {
	session *mgo.Session

	io.Closer
}

func (s mongoStorage) Close() (err error) {
	log.Println("cleanup mongo storage session")
	s.session.Close()
	return
}

func (s mongoStorage) SaveHashRequest(id string, req RequestData) (err error) {
	session := s.session.Copy()
	defer session.Close()

	c := session.DB(dbname).C(cname)

	rec := dbRecord{ID: id, Payload: req.Payload, HashRoundsCnt: req.HashRoundsCnt, Status: InProgress.String()}
	err = c.Insert(rec)

	return
}

func (s mongoStorage) SaveHashResult(id string, hash string) (err error) {
	session := s.session.Copy()
	defer session.Close()

	c := session.DB(dbname).C(cname)

	var rec dbRecord
	err = c.Find(bson.M{"id": id}).One(&rec)
	if err != nil {
		return
	}

	rec.Hash = hash
	rec.Status = Finished.String()

	err = c.Update(bson.M{"id": id}, rec)

	return
}

func (s mongoStorage) GetHashResult(id string) (resp ResponseData, err error) {
	session := s.session.Copy()
	defer session.Close()

	c := session.DB(dbname).C(cname)

	err = c.Find(bson.M{"id": id}).One(&resp)

	return
}

func ensureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB(dbname).C(cname)

	index := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

// InitMongoStorageSvc prepare mongo storage and assigns it to ServiceRepository fetcher and saver
func InitMongoStorageSvc(repo *ServiceRepository) {
	var storage mongoStorage
	var mongocon = GetEnvOrDefault("MONGOCONNECTION", "root:toor@localhost")
	var err error
	storage.session, err = mgo.Dial(mongocon)
	if err != nil {
		log.Fatalf("cannot connect to mongo @%s [%s]", mongocon, err)
	}

	ensureIndex(storage.session)

	repo.Saver = &storage
	repo.Fetcher = &storage
}
