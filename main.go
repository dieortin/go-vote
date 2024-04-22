package main

import (
	"fmt"
	"golang.org/x/exp/maps"
	"html/template"
	"log"
	"net/http"
)

var clientMap = make(ClientMap)

func registerHandlerPost(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "No name provided", http.StatusBadRequest)
		return
	}

	newUuid, err := clientMap.AddClient(name)
	if err != nil {
		http.Error(w, "Could not add new client", http.StatusInternalServerError)
		return
	}
	SetAuthCookie(w, newUuid.String())
	http.Redirect(w, r, "/", http.StatusFound)
}

func registerHandlerGet(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/register.html")
}

func votePageHandler(w http.ResponseWriter, r *http.Request) {
	client := r.Context().Value(LoggedUser)
	if client == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	err := templates.ExecuteTemplate(w, "vote.html", client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func allClientsHandler(w http.ResponseWriter, r *http.Request) {
	allClients := maps.Values(clientMap)
	err := templates.ExecuteTemplate(w, "allclients.html", allClients)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var templates = template.Must(template.ParseGlob("views/*.html"))

func main() {
	fmt.Println("Starting voting server")

	//http.HandleFunc("/", mainpageHandler)
	http.HandleFunc("GET /register", registerHandlerGet)
	http.HandleFunc("POST /register", registerHandlerPost)
	http.HandleFunc("/", Authenticate(votePageHandler))
	http.HandleFunc("/allclients", Authenticate(allClientsHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
