package main

import (
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
	c := New[any]()
	testStruct1 := testStruct{ID: 1, Username: "testUser1"}
	mcKey := fmt.Sprintf("ts_%d", testStruct1.ID)

	// test Set and Get
	c.Set(mcKey, testStruct1)
	got, ok := c.Get(mcKey)
	if diff := cmp.Diff(testStruct1, got); diff != "" || !ok {
		t.Errorf("Get() mismatch (-want +got):\n%s", diff)
	}

	// test Delete and Get
	c.Delete(mcKey)
	got, ok = c.Get(mcKey)
	if diff := cmp.Diff(nil, got); diff != "" || ok {
		t.Errorf("Get() mismatch (-want +got):\n%s", diff)
	}
}
