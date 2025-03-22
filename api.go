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
	// e.g. GET /users/{userID} // The single whitespace is important
	// A pattern with no method matches every method
	// A pattern with the method GET matches both GET and HEAD requests.
	// Otherwise, the method must match exactly.
	router := http.NewServeMux() // NewServeMux allocates and returns a new ServeMux.
	router.HandleFunc("GET /users/{userID}", func(w http.ResponseWriter, r *http.Request) {
		userID := r.PathValue("userID")
		w.Write([]byte("User ID: " + userID))
	})

	// Implementing subroutes
	// The subroutes will of course need their own handler functions
	// This is just an illustration of how to make it accessible to the /users/{userID}
	v1 := http.NewServeMux()
	v1.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	// Implementing the middlewares
	middlewareChain := MiddlewareChain( 
		// If we swap the order here it will affect the output of logs
		// e.g. in this case, swapping the order would mean that no logs will be written if the user
		// is unauthorized
		RequestLoggerMiddleware,
		RequireAuthMiddleware,
	)
	// Declaring a new instance of the http.Server type from the stdlib
	// There are more options to declare in here
	// e.g. for ReadHeaderTimeout TLS config or HTTP/2 connections or Protocols etc. 
	server := http.Server{
		Addr: s.addr,
		//Handler: router, // We now wrap the handle with the middleware
		//Handler: RequireAuthMiddleware(RequestLoggerMiddleware(router)), // To combine the two wrappers we can do it this way but it is very ugly
		Handler: middlewareChain(router),
	}

	// Adding a log during testing to show if a server has started
	log.Printf("Server has started on address %s\n", s.addr)

	return server.ListenAndServe() // Note: remember to upgrade to ListenAndServeTLS for HTTPS
}

// Making sample middleware
// Middlewares can be useful for tasks like logging or authentication etc.

// In this case our middleware will be a logger
// This logger will print out the request url and the method when the user makes a server request
func RequestLoggerMiddleware(next http.Handler) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("method %s, path: %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

// Making middleware to authenticate a route
// e.g. making some routes only accessible to signed in users with jwt
func RequireAuthMiddleware(next http.Handler) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		// check if the user is authenticated
		token := r.Header.Get("Authorization")
		if token != "Bearer token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

type Middlware func(http.Handler) http.HandlerFunc

func MiddlewareChain(middlewares ...Middlware) Middlware{
	return func(next http.Handler) http.HandlerFunc {
		for i := len(middlewares)-1 ; i>= 0 ; i-- { // working from the end to the first because the order matters
			next = middlewares[i](next) // making sure to call the next
		} 
		return next.ServeHTTP
	}
}