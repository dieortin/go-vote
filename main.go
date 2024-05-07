package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/exp/maps"
)

var clientMap = make(ClientMap)
var pollStorage = NewPollStorage()

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

func newTemplateServer(template string, value any) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := templates.ExecuteTemplate(w, template, value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func serveFile(file string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, file)
	}
}

func newPollHandler(w http.ResponseWriter, r *http.Request) {
	client := GetCtxClient(r.Context())

	name := r.FormValue("name")
	option1 := r.FormValue("option1")
	option2 := r.FormValue("option2")
	option3 := r.FormValue("option3")
	option4 := r.FormValue("option4")

	log.Printf("Received POST for form with values:\n")
	log.Printf("  Name: %s\n", name)
	log.Printf("  Option 1: %s\n", option1)
	log.Printf("  Option 2: %s\n", option2)
	log.Printf("  Option 3: %s\n", option3)
	log.Printf("  Option 4: %s\n", option4)

	poll := NewPoll(name, []string{option1, option2, option3, option4})
	id, err := pollStorage.AddPoll(poll)
	if err != nil {
		fmt.Println("Error while adding new poll to storage")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Added new poll to storage with UUID %v, redirecting to polls overview", id)
	http.Redirect(w, r, "/polls", http.StatusFound)
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
	http.HandleFunc("POST /poll", Authenticate(newPollHandler))
	http.HandleFunc("GET /poll", Authenticate(serveFile("views/newpoll.html")))
	http.HandleFunc("/", Authenticate(serveFile("views/newpoll.html")))
	http.HandleFunc("/allclients", Authenticate(allClientsHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
