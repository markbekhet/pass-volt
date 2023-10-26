package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Pipe struct {
	R   *os.File
	W   *os.File
	Err *os.File
}

func preparePipe(input string) (Pipe, Pipe) {
	old := Pipe{}
	new := Pipe{}
	old.R = os.Stdin
	old.W = os.Stdout
	old.Err = os.Stderr
	new.R, new.W, _ = os.Pipe()
	new.W.Write([]byte(input))
	os.Stdin = new.R
	os.Stdout = new.W
	os.Stderr = new.Err
	return old, new
}

func restoreStd(old, new Pipe) {
	new.W.Close()
	os.Stdin = old.R
	os.Stdout = old.W
	os.Stderr = old.Err

}

func TestAddNewUserWithEncryptedPassword(t *testing.T) {
	m := make(map[string]AccountDetails)
	testId := "gmail"
	testUser := "test_user"
	testPass := "test_password"
	input := fmt.Sprintf("%v\n%v\n%v\n", testId, testUser, testPass)

	old, new := preparePipe(input)
	Add(m)
	restoreStd(old, new)
	value, ok := m[testId]
	if !ok {
		t.Fatal("The value was not added to the map")
	}
	bytesPass := []byte(testPass)
	if reflect.DeepEqual(value.Password, bytesPass) {
		t.Fatal("The password was not encrypted")
	}
	assert.NotEqual(t, value.Password, bytesPass)

}

func TestAddExistingUserDisplaysError(t *testing.T) {
	/*
		Test setup
		Inistialize the map and add the entry before runing the Add method
	*/
	m := make(map[string]AccountDetails)
	testId := "gmail"
	testUser := "test_user"
	testPass := "test_password"
	details := AccountDetails{
		Username: testUser,
		Password: []byte(testPass),
	}
	details.encrypt()
	m[testId] = details
	input := fmt.Sprintf("%v\n%v\n%v\n", testId, testUser, testPass)
	/*===================================================*/
	old, new := preparePipe(input)
	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, new.R)
		outC <- buf.String()
	}()
	Add(m)
	/*===================================================*/

	// back to normal state
	restoreStd(old, new)
	out := <-outC
	assert.True(
		t, strings.Contains(out, "The id already existed use the update keyword instead"),
	)

}
