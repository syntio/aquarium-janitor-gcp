// Copyright Syntio d.o.o.
// All Rights Reserved
//
// Package main is the starting point for running the server.
package main

import (
	"github.com/syntio/schema-registry/rest"
)

//
// Starting point of the program.
// Just starts the REST server.
//
func main() {
	rest.SetupAndStartServer()
}
