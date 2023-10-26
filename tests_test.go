package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockStdin(input string) (*os.File, *os.File) {
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write([]byte(input)); err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	oldStdin := os.Stdin

	os.Stdin = tmpfile

	return oldStdin, tmpfile

}

func restoreStdin(oldStdin, tmp *os.File) {
	if err := tmp.Close(); err != nil {
		log.Fatal(err)
	}
	os.Stdin = oldStdin
}

func TestAddNewUserWithEncryptedPassword(t *testing.T) {
	m := make(map[string]AccountDetails)
	testId := "gmail"
	testUser := "test_user"
	testPass := "test_password"
	input := fmt.Sprintf("%v\n%v\n%v\n", testId, testUser, testPass)
	oldStdin, tmpFile := mockStdin(input)
	defer restoreStdin(oldStdin, tmpFile) // Restore original Stdin
	Add(m)
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
	oldStdin, tmpFile := mockStdin(input)
	defer restoreStdin(oldStdin, tmpFile) // Restore original Stdin
	/*===================================================*/
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()
	Add(m)
	/*===================================================*/

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC
	assert.True(
		t, strings.Contains(out, "The id already existed use the update keyword instead"),
	)

}
