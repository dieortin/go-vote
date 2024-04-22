package main

import (
	"context"
	"fmt"
	"net/http"
)

// AuthCookie Name of the cookie used for authenticating users
const AuthCookie = "auth"

const LoggedUser = "loggedUser"

func SetAuthCookie(w http.ResponseWriter, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:  AuthCookie,
		Value: value,
	})
}

func GetClient(r *http.Request) (*Client, error) {
	cookie, err := r.Cookie(AuthCookie)
	if err != nil {
		return nil, err
	}

	client, err := clientMap.GetClientByID(cookie.Value)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Client with name %v found with provided id\n", client.Name)
	return client, nil
}

func Authenticate(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client, err := GetClient(r)
		if err != nil {
			fmt.Println("Authentication failed")
			http.Redirect(w, r, "/register", http.StatusFound)
			return
		}

		fmt.Println("Authentication succeeded")
		r = r.WithContext(context.WithValue(r.Context(), LoggedUser, client))
		handler.ServeHTTP(w, r)
	}
}
