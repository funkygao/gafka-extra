package main

// Varnish stores content in pluggable modules called storage backends.
// It does this via its internal stevedore interface.
// Tt allows Varnish Cache to run multiple storage backends concurrently, even
// multiple instances of the same type.
type istevedore interface {
	init()
	open()
	alloc()
	trim()
	free()
	close()
	allocobj()
	signal_close()
}

type stevedore struct {
	istevedore

	magic uint32
	name  string

	lru *lru
}
