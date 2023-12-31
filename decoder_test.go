package bencoding

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecoder_DecodeString(t *testing.T) {

	// test cases
	var tests = []struct {
		input string
		want  string
		err   error
	}{
		{"4:spam", "spam", nil},
		{"0:", "", nil},
		{"1:a", "a", nil},
		{"2:ab", "ab", nil},
	}

	for _, test := range tests {
		str, err := DecodeString(bytes.NewReader([]byte(test.input)))

		assert.ErrorIs(t, err, test.err)
		assert.Equal(t, test.want, str)
	}
}

func TestDecoder_DecodeInt(t *testing.T) {

	// test cases
	var tests = []struct {
		input string
		want  int
		err   error
	}{
		{"i3e", 3, nil},
		{"i0e", 0, nil},
		{"i-3e", -3, nil},
		{"ie", 0, ErrInvalidFormat},
		{"i-e", 0, ErrInvalidFormat},
		{"i-0e", 0, ErrInvalidFormat},
		{"i03e", 0, ErrInvalidFormat},
		{"i1234e", 1234, nil},
	}

	for _, test := range tests {
		res, err := DecodeInt(bytes.NewReader([]byte(test.input)))

		assert.ErrorIs(t, err, test.err)
		assert.Equal(t, test.want, res)
	}
}

func TestDecoder_DecodeList(t *testing.T) {

	// test cases
	var tests = []struct {
		input string
		want  []any
		err   error
	}{
		{"l4:spam4:eggse", []any{"spam", "eggs"}, nil},
		{"l4:spam4:eggse", []any{"spam", "eggs"}, nil},
		{"l4:spam4:eggsl4:spam4:eggsee", []any{"spam", "eggs", []any{"spam", "eggs"}}, nil},
		{"l4:spami3ee", []any{"spam", 3}, nil},
		{"le", []any{}, nil},
		{"e", nil, ErrInvalidFormat},
	}

	for _, test := range tests {
		res, err := DecodeList(bytes.NewReader([]byte(test.input)))

		assert.ErrorIs(t, err, test.err)
		assert.ElementsMatch(t, test.want, res)
	}
}

func TestDecoder_DecodeDict(t *testing.T) {

	// test cases
	var tests = []struct {
		input string
		want  map[string]any
		err   error
	}{
		{"d3:cow3:moo4:spam4:eggse", map[string]any{"cow": "moo", "spam": "eggs"}, nil},
		{"d4:spaml1:a1:bee", map[string]any{"spam": []any{"a", "b"}}, nil},
		{"d9:publisher3:bob17:publisher-webpage15:www.example.com18:publisher.location4:homee", map[string]any{"publisher": "bob", "publisher-webpage": "www.example.com", "publisher.location": "home"}, nil},
		{"d4:spami3ee", map[string]any{"spam": 3}, nil},
		{"d3:cow3:moo4:spam4:eggse", map[string]any{"cow": "moo", "spam": "eggs"}, nil},
		{"de", map[string]any{}, nil},
		{"e", nil, ErrInvalidFormat},
		{"d3:cow3:moo4:spam4:eggse", map[string]any{"cow": "moo", "spam": "eggs"}, nil},
	}

	for _, test := range tests {
		res, err := DecodeDict(bytes.NewReader([]byte(test.input)))
		assert.ErrorIs(t, err, test.err)
		assert.Equal(t, test.want, res)
	}
}

func TestDecoder_Decode(t *testing.T) {

	// test cases
	var tests = []struct {
		input string
		want  any
		err   error
	}{
		{"4:spam", "spam", nil},
		{"i3e", 3, nil},
		{"l4:spam4:eggse", []any{"spam", "eggs"}, nil},
		{"d3:cow3:moo4:spam4:eggse", map[string]any{"cow": "moo", "spam": "eggs"}, nil},
		{"e", nil, ErrInvalidFormat},
	}

	for _, test := range tests {
		res, err := Decode(bytes.NewReader([]byte(test.input)))
		assert.ErrorIs(t, err, test.err)
		assert.Equal(t, test.want, res)
	}
}
