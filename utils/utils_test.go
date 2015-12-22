package utils

import (
	"testing"
	"strconv"
)

func TestrndDest(t *testing.T) {
	//going to run 5000 times quickly 
	var dirMap map[string]struct{}
	for x := 0; x < 5000; x++ {
		dirMap[strconv.Itoa(x)] = struct{}{}
	}

	if len(dirMap) != 5000 {
		t.Error("dirMap contains", len(dirMap), "items.  Should contain 5000")
	}

}