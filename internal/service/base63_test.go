package service

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestEncode(t *testing.T) {
	testTable := []struct {
		input    int64
		expected string
	}{
		{0, "aaaaaaaaaa"},
		{1, "aaaaaaaaab"},
		{10, "aaaaaaaaak"},
		{Base63TenMaxId, "__________"},
		{Base63TenMaxId - 1, "_________9"},
	}

	for _, testCase := range testTable {
		assert.Equal(t, base63Encode(testCase.input), testCase.expected)
	}
}

func TestDecode(t *testing.T) {
	testTable := []struct {
		input    string
		expected int64
	}{
		{"aaaaaaaaaa", 0},
		{"aaaaaaaaab", 1},
		{"aaaaaaaaak", 10},
		{"__________", Base63TenMaxId},
		{"_________9", Base63TenMaxId - 1},
	}

	for _, testCase := range testTable {
		id, _ := base63Decode(testCase.input)
		assert.Equal(t, id, testCase.expected)
	}
}

func TestBase63(t *testing.T) {
	testTable := []int64{0, 1, 2, 10, Base63TenMaxId, Base63TenMaxId - 1, Base63TenMaxId - 2, Base63TenMaxId - 100, 1000, 72, 100, 15, 11111}

	for _, testCase := range testTable {
		encoded := base63Encode(testCase)
		decoded, _ := base63Decode(encoded)
		assert.Equal(t, decoded, testCase)
	}
}
