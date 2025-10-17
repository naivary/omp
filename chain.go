package main

import (
	"net/http"
	"slices"
)

// chain is a list of middleware functions to execute in the given order.
type chain []func(http.Handler) http.Handler

func (c chain) then(h http.Handler) http.Handler {
	for _, mw := range slices.Backward(c) {
		h = mw(h)
	}
	return h
}
