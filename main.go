package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/exp/maps"
)

var clientMap = make(ClientMap)
var pollStorage = NewPollStorage()

var pageMap = NewPageMap("base.tmpl", []string{"newpoll.tmpl", "polls.tmpl", "register.tmpl"})

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

func renderTemplate(tmpl string, value any) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Rendering template %q\n", tmpl)
		err := pageMap.ExecuteTemplate(w, tmpl, value)
		if err != nil {
			log.Printf("Error while executing template %q", tmpl)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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

	poll := NewPoll(client, name, []string{option1, option2, option3, option4})
	id, err := pollStorage.AddPoll(poll)
	if err != nil {
		fmt.Println("Error while adding new poll to storage")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Added new poll to storage with UUID %v, redirecting to polls overview", id)
	http.Redirect(w, r, "/polls", http.StatusFound)
}

func allPollsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate("polls.tmpl", maps.Values(pollStorage))(w, r)
}

func redirectHandler(url string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Redirecting to %q", url)
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func main() {
	log.Println("Starting go-vote server")

	http.HandleFunc("GET /register", renderTemplate("register.tmpl", nil))
	http.HandleFunc("POST /register", registerHandlerPost)
	http.HandleFunc("POST /polls/new", Authenticate(newPollHandler))
	http.HandleFunc("GET /polls/new", Authenticate(renderTemplate("newpoll.tmpl", nil)))
	http.HandleFunc("GET /polls", allPollsHandler)
	http.HandleFunc("/", redirectHandler("/polls/new"))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
