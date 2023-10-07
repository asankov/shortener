package random

import (
	"math/rand"
	"time"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	allBytes    = letterBytes + "!#$%&*+,-./0123456789:;<=>?[]]^_~"
)

var (
	def = New()
)

type Random struct {
	rand *rand.Rand
}

func New() *Random {
	return &Random{
		rand: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

// ID generates random ID with lenght n using the default Random.
func ID(n int) string {
	return def.ID(n)
}

// Password generates random password with length n using the default Random.
func Password(n int) string {
	return def.Password(n)
}

// ID generates random ID with lenght n.
func (r *Random) ID(n int) string {
	return r.genFrom(n, letterBytes)
}

// Password generates random password with length n.
//
// The difference between ID and Password is that for password all printable ASCII characters are used,
// while ID is using only the letters.
func (r *Random) Password(n int) string {
	return r.genFrom(n, allBytes)
}

func (r *Random) genFrom(n int, source string) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = source[r.rand.Intn(len(source))]
	}
	return string(b)
}
