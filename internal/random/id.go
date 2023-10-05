package random

import (
	"math/rand"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().Unix()))

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const allBytes = letterBytes + "!#$%&()*+,-./0123456789:;<=>?[\\]^_{|}~"

// ID generates random ID with lenght n.
func ID(n int) string {
	return genFrom(n, letterBytes)
}

// Password generates random password with length n.
//
// The difference between ID and Password is that for password all printable ASCII characters are used,
// while ID is using only the letters.
func Password(n int) string {
	return genFrom(n, allBytes)
}

func genFrom(n int, source string) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = source[r.Intn(len(source))]
	}
	return string(b)
}
