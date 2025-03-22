package main

import (
	"log"
	"net/http"
)

// This will be the foundation for the implementation of the server
// This project will attempt to use the stdlib instead of importing external dependencies
type APIServer struct {
	// The APIServer struct will hold dependencies e.g. addresses, database etc.
	addr string 
}

// Constructor function for the APIServer
func NewAPIServer(addr string) *APIServer {
	// The NewAPIServer constructor function takes in an address and returns a pointer to the APIServer struct
	return &APIServer {
		addr: addr,
	}
}

// Run method to run the server
func (s *APIServer) Run() error{
	// Note: ServeMux is an HTTP request multiplexer
	// It matches the URL of each incoming request against a list of registered patterns and calls the
	// handler for the pattern that most closely matches the URL.
	// In general a pattern looks like [METHOD ][HOST]/[PATH]
	router := http.NewServeMux() // NewServeMux allocates and returns a new ServeMux.

	// Declaring a new instance of the http.Server type from the stdlib
	// There are more options to declare in here
	// e.g. for ReadHeaderTimeout TLS config or HTTP/2 connections or Protocols etc. 
	server := http.Server{
		Addr: s.addr,
		Handler: router,
	}
	
	// Adding a log during testing to show if a server has started
	log.Printf("Server has started on address: %s", s.addr)

	return server.ListenAndServe() // Note: remember to upgrade to ListenAndServeTLS for HTTPS
}