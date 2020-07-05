package markov

import "testing"

func TestNgram(t *testing.T) {
	s := []string{"some", "random", "string", "content"}
	res := ngram(s, 3)
	if res[0][0] != "some" && res[0][1] != "random" {
		t.Error("incorrect result:", res[0])
	}
	t.Error(res)
	// if res[0][0] != "some" && res[0][1] != "random" {
	// 	t.Error("incorrect result:", res[0])
	// }
}

func TestReadFileContent(t *testing.T) {
	s := readFileContent("content.txt")      // # This file should be in the current directory.
	if s != "some random string content\n" { // # Note pesky newline.
		t.Error("incorrect:" + s)
	}

}

func TestWeidhtedChoice(t *testing.T) {
	// # The random natura of this test makes it fail randomly (though unlikely).

	{
		strs := []string{"some", "random", "string", "content"}
		nums := []float32{20, 1, 1, 1}
		if res, _ := weightedChoice(strs, nums); res != "some" {
			t.Errorf("unlikely error with result '%s'. try again?", res)
		}
	}
	{
		strs := []string{"some", "random", "string", "content"}
		nums := []float32{1, 1, 20, 1}
		if res, _ := weightedChoice(strs, nums); res != "string" {
			t.Errorf("unlikely error with result '%s'. try again?", res)
		}
	}

}
