package cache

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testStruct struct {
	ID       int
	Username string
}

func TestSetGetDelete(t *testing.T) {
	// TODO: add individual unit tests with additional types
	// TODO: add testing to verify cache doesn't exibit race conditions
	c := NewCache()
	testStruct1 := testStruct{ID: 1, Username: "testUser1"}
	mcKey := fmt.Sprintf("ts_%d", testStruct1.ID)

	// test Set and Get
	b, err := json.Marshal(testStruct1)
	if err != nil {
		t.FailNow()
	}
	c.Set(mcKey, b)
	gotBytes, ok := c.Get(mcKey)

	var gotStruct testStruct
	json.Unmarshal(gotBytes, &gotStruct)
	if diff := cmp.Diff(testStruct1, gotStruct); diff != "" || !ok {
		t.Errorf("Set() Get() mismatch (-want +got):\n%s", diff)
	}

	// test Delete and Get
	c.Delete(mcKey)
	var want []byte
	got, ok := c.Get(mcKey)
	if diff := cmp.Diff(want, got); diff != "" || ok {
		t.Errorf("Delete() Get() mismatch (-want +got):\n%s", diff)
	}
}
