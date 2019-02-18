package sha2svc

import (
	"testing"
)

func TestString(t *testing.T) {
	tables := []struct {
		arg ReqStatus
		res string
	}{
		{Failed, "failed"},
		{InProgress, "in progress"},
		{Finished, "finished"},
	}

	for _, table := range tables {
		if table.res != table.arg.String() {
			t.Errorf("Wrong string value for %s", table.res)
		}
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("wrong status passed defence")
		}
	}()

	var invalid ReqStatus = 10
	_ = invalid.String()
}
