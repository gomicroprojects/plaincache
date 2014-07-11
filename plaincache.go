// The package clause (http://golang.org/ref/spec#Package_clause)
// "main" indicates that this file is a command and will generate an executable
// binary.
// Take a look at http://golang.org/doc/code.html for an introduction on
// "How to Write Go Code"
package main

// The import declaration (http://golang.org/ref/spec#Import_declarations)
// This is in the multiline format.
// You could also write:
//     import a
//     import b
import (
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

	http.ListenAndServe(addr.String(), nil)
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
