/*
plaincache - a simple in-memory cache with as REST-ful API

Installation:
    go get github.com/gomicroprojects/plaincache

Usage:

    plaincache server_address

Example:

    plaincache :8080
    plaincache 127.0.0.1:http
    plaincache [::1]:http

Will serve HTTP on the given address and exposes three methods:

    GET /my/key - Will return the set value or HTTP staus code 404
    POST /my/key - Will set the value to the POST body
    DELETE /my/key - Will delete the entry if it exists

This is the package doc. It can be placed in a doc.go file, which is a convention.

See http://blog.golang.org/godoc-documenting-go-code
*/
package main
