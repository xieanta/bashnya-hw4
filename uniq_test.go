package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	Name   string
	Input  string
	Output string
	Flags  Flags
}

var testData = []TestCase{
	{
		Name: "01.txt",
		Input: `I love music.
I love music.
I love music.

I love music of Kartik.
I love music of Kartik.
Thanks.`,
		Output: `I love music.

I love music of Kartik.
Thanks.`,
		Flags: Flags{C: false, D: false, U: false, F: 0, S: 0, I: false},
	},
	{
		Name: "02c.txt",
		Input: `I love music.
I love music.
I love music.

I love music of Kartik.
I love music of Kartik.
Thanks.`,
		Output: `3 I love music.
1 
2 I love music of Kartik.
1 Thanks.`,
		Flags: Flags{C: true, D: false, U: false, F: 0, S: 0, I: false},
	},
	{
		Name: "03d.txt",
		Input: `I love music.
I love music.
I love music.

I love music of Kartik.
I love music of Kartik.
Thanks.`,
		Output: `I love music.
I love music of Kartik.`,
		Flags: Flags{C: false, D: true, U: false, F: 0, S: 0, I: false},
	},
	{
		Name: "04u.txt",
		Input: `I love music.
I love music.
I love music.

I love music of Kartik.
I love music of Kartik.
Thanks.`,
		Output: `
Thanks.`,
		Flags: Flags{C: false, D: false, U: true, F: 0, S: 0, I: false},
	},
	{
		Name: "05i.txt",
		Input: `I LOVE MUSIC.
I love music.
I LoVe MuSiC.

I love MuSIC of Kartik.
I love music of kartik.
Thanks.`,
		Output: `I LOVE MUSIC.

I love MuSIC of Kartik.
Thanks.`,
		Flags: Flags{C: false, D: false, U: false, F: 0, S: 0, I: true},
	},
	{
		Name: "06f1.txt",
		Input: `We love music.
I love music.
They love music.

I love music of Kartik.
We love music of Kartik.
Thanks.`,
		Output: `We love music.

I love music of Kartik.
Thanks.`,
		Flags: Flags{C: false, D: false, U: false, F: 1, S: 0, I: false},
	},
	{
		Name: "07s1.txt",
		Input: `I love music.
A love music.
C love music.

I love music of Kartik.
We love music of Kartik.
Thanks.`,
		Output: `I love music.

I love music of Kartik.
We love music of Kartik.
Thanks.`,
		Flags: Flags{C: false, D: false, U: false, F: 0, S: 1, I: false},
	},
	{
		Name: "08f2.txt",
		Input: `I love music.
A love music.
C love music.

I love music of Kartik.
We love music of Kartik.
Thanks.`,
		Output: `I love music.

I love music of Kartik.
Thanks.`,
		Flags: Flags{C: false, D: false, U: false, F: 2, S: 0, I: false},
	},
	{
		Name: "09ics2.txt",
		Input: `I lovE music.
A lovE music.
C love music.

I love music of Kartik.
We love music of Kartik.
Thanks.`,
		Output: `3 I lovE music.
1 
1 I love music of Kartik.
1 We love music of Kartik.
1 Thanks.`,
		Flags: Flags{C: true, D: false, U: false, F: 0, S: 2, I: true},
	},
	{
		Name: "10cif1.txt",
		Input: `I lovE music.
A lovE music.
C love music.

I love music of Kartik.
We love music of Kartik.
Thanks.`,
		Output: `3 I lovE music.
1 
2 I love music of Kartik.
1 Thanks.`,
		Flags: Flags{C: true, D: false, U: false, F: 1, S: 0, I: true},
	},
	{
		Name: "11ui.txt",
		Input: `I Love music.
I Love music.
I love music.

I love music of Kartik.
I love music of Kartik.
Thanks.`,
		Output: `
Thanks.`,
		Flags: Flags{C: false, D: false, U: true, F: 0, S: 0, I: true},
	},
	{
		Name: "12f1s1.txt",
		Input: `I love music.
A love music.
C love music.

I love music of Kartik.
We love music of Kartik.
Thanks.`,
		Output: `I love music.

I love music of Kartik.
Thanks.`,
		Flags: Flags{C: false, D: false, U: false, F: 1, S: 1, I: false},
	},
}

func TestOK(t *testing.T) {
	for _, tc := range testData {
		t.Run(tc.Name, func(t *testing.T) {
			output := &bytes.Buffer{}
			d := strings.Split(tc.Input, "\n")
			processCommand(tc.Flags, d, output)
			result := strings.TrimSpace(output.String())
			expected := strings.TrimSpace(tc.Output)
			assert.Equal(t, expected, result, "Test %s", tc.Name)
		})
	}
}

var badTests = []struct {
	name  string
	flags Flags
	data  []string
}{
	{
		name:  "empty input",
		flags: Flags{},
		data:  []string{},
	},
	{
		name:  "nil input",
		flags: Flags{},
		data:  nil,
	},
	{
		name:  "with empty strings",
		flags: Flags{},
		data:  []string{"hello", "", "world", ""},
	},
	{
		name:  "all same lines",
		flags: Flags{},
		data:  []string{"a", "a", "a", "a"},
	},
	{
		name:  "strange characters",
		flags: Flags{},
		data:  []string{"\x00", "\n", "\t", "норм"},
	},
}

func TestFail(t *testing.T) {
	for _, tc := range badTests {
		t.Run(tc.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			assert.NotPanics(t, func() {
				processCommand(tc.flags, tc.data, output)
			}, "processCommand should not panic on %s", tc.name)
			o := output.String()
			fmt.Println(o)
		})
	}
}
