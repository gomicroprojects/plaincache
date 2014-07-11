plaincache
==========

A very simple in-memory cache with a RESTful API

## The idea

* This will be a very simple and limited in-memory cache service. 
* We will only use the Go stdlib.
* No auth.
* String keys/values.
* No automatic eviction.

### The API

We will create a RESTful API over HTTP.

The whole path will act as the key.

`GET /the/key` Will return the value or a HTTP status code 404 (Not Found) if the key does not exist.

`POST /the/key` Will set the value of the key to the contents of the POST body.

`DELETE /the/key` Will delete the contents of the key.

## Starting

Start with plaincache.go

## www.gomicroprojects.com

This is a project for http://www.gomicroprojects.com