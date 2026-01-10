package diff_test

import (
	"reflect"
	"testing"

	"github.com/LeeFred3042U/kitcat/internal/diff"
)

func TestMyersDiff(t *testing.T) {
	tests := []struct {
		name     string
		text1    []string
		text2    []string
		expected []diff.Diff[string]
	}{
		{
			name:     "Empty vs Empty",
			text1:    []string{},
			text2:    []string{},
			expected: nil, // Or empty slice depending on implementation details
		},
		{
			name:  "Insert Only",
			text1: []string{},
			text2: []string{"a", "b"},
			expected: []diff.Diff[string]{
				{Operation: diff.INSERT, Text: []string{"a", "b"}},
			},
		},
		{
			name:  "Delete Only",
			text1: []string{"a", "b"},
			text2: []string{},
			expected: []diff.Diff[string]{
				{Operation: diff.DELETE, Text: []string{"a", "b"}},
			},
		},
		{
			name:  "Replace",
			text1: []string{"a"},
			text2: []string{"b"},
			expected: []diff.Diff[string]{
				{Operation: diff.DELETE, Text: []string{"a"}},
				{Operation: diff.INSERT, Text: []string{"b"}},
			},
		},
		{
			name:  "Identical",
			text1: []string{"a", "b"},
			text2: []string{"a", "b"},
			expected: []diff.Diff[string]{
				{Operation: diff.EQUAL, Text: []string{"a", "b"}},
			},
		},
		{
			name:  "Mixed",
			text1: []string{"a", "b", "c"},
			text2: []string{"a", "d", "c"},
			expected: []diff.Diff[string]{
				{Operation: diff.EQUAL, Text: []string{"a"}},
				{Operation: diff.DELETE, Text: []string{"b"}},
				{Operation: diff.INSERT, Text: []string{"d"}},
				{Operation: diff.EQUAL, Text: []string{"c"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := diff.NewMyersDiff(tt.text1, tt.text2)
			got := d.Diffs()

			if len(got) != len(tt.expected) {
				t.Errorf("Diffs() length = %v, want %v\nGot: %v\nWant: %v", len(got), len(tt.expected), got, tt.expected)
				return
			}

			for i := range got {
				if got[i].Operation != tt.expected[i].Operation {
					t.Errorf("Diffs()[%d].Operation = %v, want %v", i, got[i].Operation, tt.expected[i].Operation)
				}
				if !reflect.DeepEqual(got[i].Text, tt.expected[i].Text) {
					// Handle nil vs empty slice distinction if helpful, but DeepEqual handles []string{} vs nil as different.
					// We should probably allow nil/empty mismatch if both are effective empty.
					if len(got[i].Text) == 0 && len(tt.expected[i].Text) == 0 {
						continue
					}
					t.Errorf("Diffs()[%d].Text = %v, want %v", i, got[i].Text, tt.expected[i].Text)
				}
			}
		})
	}
}
