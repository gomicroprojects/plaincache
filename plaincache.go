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
	// Uncomment this for a run-time panic.
	// panic: runtime error: assignment to entry in nil map
	// cache["a"] = "b"
	cache = make(map[string]string)
}
