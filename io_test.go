package dockerlib_test

import (
	"github.com/ATenderholt/dockerlib"
	"strings"
	"testing"
)

func TestReadLinesAsBytes(t *testing.T) {
	str := "{\"status\":\"Pulling from library/python\",\"id\":\"3.9.11-alpine3.14\"}\r\n{\"status\":\"Digest: sha256:da20794e1b03c80c6a21918ba8b7958886783ad94268a970fc34534e4e577a72\"}\r\n{\"status\":\"Status: Image is up to date for python:3.9.11-alpine3.14\"}\r\n\"}"
	reader := strings.NewReader(str)
	c := dockerlib.ReadLinesAsBytes(reader)

	b := <-c

	got := string(b)
	exp := "{\"status\":\"Pulling from library/python\",\"id\":\"3.9.11-alpine3.14\"}"
	if got != exp {
		t.Errorf("Expected %s, got %s", exp, got)
	}
}
