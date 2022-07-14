package controllers

import (
	"fmt"
	"net/http"

	"authentication.com/workers"
)

/**
HandleHome is method that handle the request to home page (/)
**/
func HandleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Authentication Service")
}

/**
HandleLogin is method that call DoLogin that handle the login process
**/
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	workers.DoLogin(w, r)
}

/**
HandleAuthentication is method that call DoAuthenticationRequest that handle the authentication
**/
func HandleAuthentication(w http.ResponseWriter, r *http.Request) {
	workers.DoAuthenticationRequest(w, r)
}

/**
HandleRefreshToken is method that call DoHandleRefreshToken that handle the authentication
**/
func HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	workers.DoHandleRefreshToken(w, r)
}
