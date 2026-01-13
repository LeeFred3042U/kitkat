package core

import (
	"reflect"
	"testing"
)

func TestParseTodo(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "ValidCommands",
			input: "pick abc123 Commit message\ndrop def456 Another message",
			want: []string{
				"pick abc123 Commit message",
				"drop def456 Another message",
			},
		},
		{
			name:  "InvalidCommands",
			input: "foo abc123 Invalid\npick abc123 Valid\nbar def456 Invalid2",
			want: []string{
				"foo abc123 Invalid",
				"pick abc123 Valid",
				"bar def456 Invalid2",
			},
		},
		{
			name:  "EmptyInput",
			input: "",
			want:  []string{},
		},
		{
			name:  "CommentsAndWhitespace",
			input: "# This is a comment\n   \npick abc123 Message\n   # Another comment\ndrop def456 Message2\n   ",
			want: []string{
				"pick abc123 Message",
				"drop def456 Message2",
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := parseTodo(tc.input)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
