package utils

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

func TestReadLines(t *testing.T) {
	/*** setUp ***/
	file, err := os.CreateTemp("", "gotest")
	if err != nil {
		t.Errorf("ERROR: %s", err.Error())
	}
	defer func() {
		err := os.Remove(file.Name())
		if err != nil {
			t.Errorf("ERROR: %s", err.Error())
		}
	}()

	for i := 0; i < 3; i++ {
		_, err := file.Write([]byte(fmt.Sprintf("%d\n", i)))
		if err != nil {
			t.Errorf("ERROR: %s", err.Error())
		}
	}

	/*** test ***/
	got := ReadLines(file.Name(), func(i int64, line string) bool {
		if line != strconv.FormatInt(i, 10) {
			t.Errorf("ERROR: %d", i)
		}
		return true
	})
	if got != nil {
		t.Errorf("ERROR: %s", got.Error())
	}
}

func TestReadFileErrNotFounc(t *testing.T) {
	/*** setUp ***/
	file, err := os.CreateTemp("", "gotest")
	if err != nil {
		t.Errorf("ERROR: %s", err.Error())
	}

	err = os.Remove(file.Name())
	if err != nil {
		t.Errorf("ERROR: %s", err.Error())
	}

	/*** test ***/
	got := ReadLines(file.Name(), func(i int64, line string) bool {
		return true
	})
	if got == nil {
		t.Errorf("Should return an error")
	}
}
