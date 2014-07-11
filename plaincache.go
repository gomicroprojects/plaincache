package main

// package main
//
// The package clause (http://golang.org/ref/spec#Package_clause)
// "main" indicates that this file is a command and will generate an executable
// binary.
//
// Take a look at http://golang.org/doc/code.html for an introduction on
// "How to Write Go Code"

// The import declaration (http://golang.org/ref/spec#Import_declarations)
// This is in the multiline format.
// You could also write:
//     import a
//     import b
import (
	// The io/ioutil package (http://golang.org/pkg/io/ioutil/)
	// We will use this to read the POST body
	"io/ioutil"
	// The time package (http://golang.org/pkg/time/)
	// We use it to set our timeouts
	"time"
	// The net package (http://golang.org/pkg/net/)
	// We will use this package to check for the correctness of the address argument
	"net"
	// The fmt package (http://golang.org/pkg/fmt/) for formatted I/O
	"fmt"
	// We need exactly one argument for our program: the server address.
	// The os package (http://golang.org/pkg/os/) gives us access to the
	// args through os.Args []string
	"os"
	// The net/http package (http://golang.org/pkg/net/http/)
	// We will create our HTTP server with this package
	"net/http"
	// The sync package (http://golang.org/pkg/sync/) for synchronization primitives
	// We will use a sync.RWMutex to protect our cache from race conditions.
	"sync"
)

// We declare a global variable "m". The type is a sync.RWMutex (an RWMutex in the sync package)
//
// We need to protect our cache from concurrent writes, since we will be serving it through
// HTTP. Every request is served in its own goroutine.
//
// If you want to see where that happens, see the http.Server.Serve() method:
//
//    func (srv *Server) Serve(l net.Listener) error {
//        ...
//        go c.serve()
//    }
var m sync.RWMutex

// The cache. We will use a map (http://golang.org/ref/spec#Map_types)
// This is just a declaration of the variable "cache".
// If you would use the "cache" map, you would receive a run-time panic
// since the map is not initialized yet.
var cache map[string]string

// The main function. This is the starting point when your program executes.
func main() {
	// Check for the server address argument
	// os.Args[0] is the program name
	if len(os.Args) == 1 {
		// See the printUsage() function definition below
		printUsage()
		// Exit the program with an error code (!= 0)
		os.Exit(1)
	}
	// net.ResolveTCPAddr has multiple (here two) return values
	// The second one is of type error
	addr, err := net.ResolveTCPAddr("tcp", os.Args[1])
	// Always check for errors if a function returns an error
	if err != nil {
		fmt.Printf("error resolving address: %s\n\n", err)
		printUsage()
		os.Exit(1)
	}

	// Uncomment this for a run-time panic.
	// panic: runtime error: assignment to entry in nil map
	// cache["a"] = "b"

	// Now we initialize the map
	cache = make(map[string]string)

	// Create an HTTP server
	// You can read this as: server is a pointer to(take the address of) an
	// http.Server with fields...
	//
	// This is equivalent to:
	// server := new(http.Server)
	// server.Addr = addr.String()
	// server.Handler = handler()
	// etc.
	server := &http.Server{
		Addr:         addr.String(),
		Handler:      handler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	// Now start the server
	// This method call will block until an error occurs (which is usually fatal)
	//
	// The net/http server will serve every request in a separate goroutine
	// (http://golang.org/doc/effective_go.html#goroutines)
	// This means that our server is implicitly concurrent, therefore
	// we will have to protect our state from race conditions
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("fatal error while serving: %s", err)
		os.Exit(1)
	}
}

// This function will print a simple help text
func printUsage() {
	// http://golang.org/pkg/fmt/
	fmt.Println("Usage:")
	// use the program name
	fmt.Printf("%s server_address\n", os.Args[0])
	fmt.Println("\nExample:")
	fmt.Printf("%s :8080\n", os.Args[0])
	fmt.Printf("%s 127.0.0.1:8080\n", os.Args[0])
	fmt.Printf("%s [::1]:http\n", os.Args[0])
}

// This will create an http.Handler to be used in our HTTP server
func handler() http.Handler {
	// http.HandlerFunc takes a function value and turns it into a type
	// which implements the http.Handler interface
	// We turn our anonymous function func(http.ResponseWriter, *http.Request)
	// into an http.Handler here
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// switch between our supported HTTP methods
		switch r.Method {
		case "GET":
			// call our GET implementation
			get(w, r)
		case "POST":
			post(w, r)
		case "DELETE":
			deleteCache(w, r)
		default:
			// or send the status Method Not Allowed in any other case
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}

// implementing the HTTP GET method
func get(w http.ResponseWriter, r *http.Request) {
	// set the Content-Type header
	w.Header().Set("Content-Type", "text/plain")
	// set a read lock on our RWMutex
	// while we hold this lock, we cannot get a write lock (m.Lock())
	m.RLock()
	// Look up the entry in our map
	// We assign the index expression to two variables
	// ok will be false if the key does not exist
	value, ok := cache[r.URL.Path]
	// release the read lock
	m.RUnlock()
	if !ok {
		// no entry for given key
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// http.ResponseWrite implements the io.Writer interface
	// we will just cast the string value to a byte slice []byte
	// We are allowed to do this, because a string is a sequence of bytes
	// see http://golang.org/ref/spec#String_types
	_, err := w.Write([]byte(value))
	if err != nil {
		fmt.Printf("error writing response on GET from %s: %s\n", r.RemoteAddr, err)
	}
}

// implementing the HTTP POST method, setting values
func post(w http.ResponseWriter, r *http.Request) {
	// read the full body
	// the ioutil helps us with the handling of the io.Reader (the r.Body)
	// we want to read everything until EOF
	value, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error reading the POST body from %s: %s\n", r.RemoteAddr, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// lock our RWMutex for writing
	// while we hold this lock, no other read or write lock can be obtained
	m.Lock()
	// set the value
	cache[r.URL.Path] = string(value)
	// release the lock
	m.Unlock()
	// set the Content-Type header
	// it will be text/plain and we will echo the value back
	w.Header().Set("Content-Type", "text/plain")
	_, err = w.Write(value)
	if err != nil {
		fmt.Printf("error on writing POST response: %s", err)
	}
}

// deleting values
// delete() is a reserved keyword, so we have to name it differently
func deleteCache(w http.ResponseWriter, r *http.Request) {
	// this is pretty simple by now
	m.Lock()
	// http://golang.org/ref/spec#Deletion_of_map_elements
	delete(cache, r.URL.Path)
	m.Unlock()
}
