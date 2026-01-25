package main

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"time"
)

type UrlMap interface {
	AddMapping(s string)
	RetrieveUrl(s string) (bool, string)
	RemoveMapping(s string) bool
}

type UrlMapping struct {
	originalUrl string
	createdAt   time.Time
}

const SHORT_URL_LENGTH int = 10

var CHAR_TYPES = [3]rune{48, 65, 97} // ascii start value for numbers, upper & lower letters

func Shorten() string {
	var result strings.Builder
	var charPos int
	var charType int

	for range SHORT_URL_LENGTH {

		if charType = rand.IntN(3); charType == 0 {
			charPos = rand.IntN(10)
		} else {
			charPos = rand.IntN(26)
		}

		char := CHAR_TYPES[charType] + rune(charPos)
		fmt.Fprintf(&result, "%c", char)

	}

	return result.String()
}
