package core

import (
	"reflect"
	"testing"
)

func TestParseTodo_ValidCommands(t *testing.T) {
	input := "pick abc123 Commit message\ndrop def456 Another message"
	got := parseTodo(input)
	want := []string{
		"pick abc123 Commit message",
		"drop def456 Another message",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseTodo_InvalidCommands(t *testing.T) {
	input := "foo abc123 Invalid\npick abc123 Valid\nbar def456 Invalid2"
	got := parseTodo(input)
	want := []string{
		"foo abc123 Invalid",
		"pick abc123 Valid",
		"bar def456 Invalid2",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseTodo_EmptyInput(t *testing.T) {
	got := parseTodo("")
	if len(got) != 0 {
		t.Errorf("expected empty result, got %v", got)
	}
}

func TestParseTodo_CommentsAndWhitespace(t *testing.T) {
	input := "# This is a comment\n   \npick abc123 Message\n   # Another comment\ndrop def456 Message2\n   "
	got := parseTodo(input)
	want := []string{
		"pick abc123 Message",
		"drop def456 Message2",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
