// Package version holds the running tea version for code that cannot take it
// as an argument: the privileged helper reports it, and registration compares
// the installed helper against it.
package version

// Current is set by main from main.version; see main.go.
var Current = "dev"
