package bencoding

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncoder_EncodeString(t *testing.T) {

	// test cases
	var tests = []struct {
		input string
		want  string
		err   error
	}{
		{"spam", "4:spam", nil},
		{"", "0:", nil},
		{"a", "1:a", nil},
		{"ab", "2:ab", nil},
	}

	for _, test := range tests {
		res, err := EncodeString(test.input)

		assert.ErrorIs(t, err, test.err)
		assert.Equal(t, test.want, string(res))
	}
}

func TestEncoder_EncodeInt(t *testing.T) {

	// test cases
	var tests = []struct {
		input int
		want  string
		err   error
	}{
		{3, "i3e", nil},
		{0, "i0e", nil},
		{-3, "i-3e", nil},
		{1234, "i1234e", nil},
	}

	for _, test := range tests {
		res, err := EncodeInt(test.input)

		assert.ErrorIs(t, err, test.err)
		assert.Equal(t, test.want, string(res))
	}
}

func TestEncoder_EncodeList(t *testing.T) {

	// test cases
	var tests = []struct {
		input []any
		want  string
		err   error
	}{
		{[]any{"spam", "eggs"}, "l4:spam4:eggse", nil},
		{[]any{"spam", 1234}, "l4:spami1234ee", nil},
		{[]any{"spam", []any{"eggs", 1234}}, "l4:spaml4:eggsi1234eee", nil},
	}

	for _, test := range tests {
		res, err := EncodeList(test.input)

		assert.ErrorIs(t, err, test.err)
		assert.Equal(t, test.want, string(res))
	}
}

func TestEncoder_EncodeDict(t *testing.T) {

	// test cases
	var tests = []struct {
		input map[string]any
		want  string
		err   error
	}{
		{map[string]any{"spam": "eggs"}, "d4:spam4:eggse", nil},
		{map[string]any{"spam": 1234}, "d4:spami1234ee", nil},
		{map[string]any{"spam": []any{"eggs", 1234}}, "d4:spaml4:eggsi1234eee", nil},
	}

	for _, test := range tests {
		res, err := EncodeDict(test.input)

		assert.ErrorIs(t, err, test.err)
		assert.Equal(t, test.want, string(res))
	}
}

func TestEncoder_Encoder(t *testing.T) {

	// test cases
	var tests = []struct {
		input any
		want  string
		err   error
	}{
		{"spam", "4:spam", nil},
		{3, "i3e", nil},
		{[]any{"spam", "eggs"}, "l4:spam4:eggse", nil},
		{map[string]any{"spam": "eggs"}, "d4:spam4:eggse", nil},
	}

	for _, test := range tests {
		res, err := Encode(test.input)

		assert.ErrorIs(t, err, test.err)
		assert.Equal(t, test.want, string(res))
	}
}
