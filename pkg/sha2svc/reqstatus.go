package sha2svc

/***
//TODO: Design: we really don't need this and perform "status changing" on storage level
//      We can use flag, hash value presence or even 2 different record to infer the status
***/

// ReqStatus type of service incoming request processing status
type ReqStatus uint8

const (
	// InProgress request to service in processing and no result available yet
	InProgress ReqStatus = iota
	// Finished request to service processed and there is a result saved and ready to consume
	Finished
	// Failed request processing failed or didn't start
	Failed
)

// String function converts ReqStatus constants into string representation
func (stat ReqStatus) String() string {
	names := [...]string{
		"in progress",
		"finished",
		"failed",
	}

	if !stat.IsValid() {
		panic("Unknown request status")
	}

	return names[stat]
}

// IsValid checks validity of request status
func (stat ReqStatus) IsValid() bool {
	return stat >= InProgress || stat <= Failed
}
