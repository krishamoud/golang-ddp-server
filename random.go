package server

import (
	"math/rand"
)

var dict = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Id generates a random string of length n and returns it
func (s *Server) Id(n int) string {
	a := make([]rune, n)
	l := len(dict)

	for i := range a {
		a[i] = dict[rand.Intn(l)]
	}

	return string(a)
}
