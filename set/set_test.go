package main

import "testing"

func TestInclude(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  bool
	}{
		{
			name:  "Set must include input value",
			input: 3,
			want:  true,
		},
		{
			name:  "Set must not include input value",
			input: 0,
			want:  false,
		},
	}

	for _, tt := range tests {
		set := New(1, 2, 3, 4, 5)
		if set.Includes(tt.input) != tt.want {
			t.Errorf("%v is not included\n", tt.input)
		}
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name  string
		input int
	}{
		{
			name:  "Add 0",
			input: 0,
		},
		{
			name:  "Add 7",
			input: 7,
		},
		{
			name:  "Add -1",
			input: -1,
		},
	}

	for _, tt := range tests {
		set := New(1, 3, 5)
		set.Add(tt.input)
		if !set.Includes(tt.input) {
			t.Errorf("%v is not included\n", tt.input)
		}
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		name  string
		input int
	}{
		{
			name:  "Remove int value",
			input: 1,
		},
	}

	for _, tt := range tests {
		set := New(1, 3, 5)
		set.Remove(tt.input)
		if set.Includes(tt.input) {
			t.Errorf("%v has not been removed\n", tt.input)
		}
	}
}
