package cli

import (
	"regexp"
	"testing"
)

func TestRegex(t *testing.T) {

	data := "dank bank stank"
	x := regexp.MustCompile("^dank bank$")
	res := x.FindAll([]byte(data), 0)
	t.Log("####\n", res)
}
