package main

import (
    "math/rand"
	  "time"
	  "net/http"
	  "github.com/go-martini/martini"
  	"html/template"
    "gopkg.in/mgo.v2"
	  "gopkg.in/mgo.v2/bson"
	  "fmt"
)

var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

type Paste struct {
	Id   string   `json: "id"`
	Title string `json: "title"`
	Time int64 `json: "time"`
	Content string `json: "content"`
}

func main() {
    rand.Seed(time.Now().UTC().UnixNano())
	  connection, err := mgo.Dial("localhost") // Connect to the database
    if err != nil {
	    // uh oh
	    panic(err) // I AM PANICKING!
	 }
	  coll := connection.DB("pastes").C("pastes")

	  m := martini.Classic()

	  m.Get("/", func(w http.ResponseWriter,r *http.Request) {

		http.ServeFile(w, r, "create.html")

	})

	m.Post("/save", func(w http.ResponseWriter,r *http.Request) {
	//	http.ServeFile(w, r, "create.html")
		err := r.ParseForm()

		if err != nil {
			panic(err)
		}
		paste := &Paste{Id: GenHash(), Title: r.FormValue("title"), Time: time.Now().Unix(), Content: r.FormValue("content")}

		coll.Insert(paste)

		http.Redirect(w, r, "view/" + paste.Id, http.StatusFound) // Redirect cause i am lazy
	})

	m.Get("/view/:id", func(w http.ResponseWriter,r *http.Request, params martini.Params) {
        id := params["id"] // Get the id

        if id == "" {
					w.WriteHeader(404)
				}

			  var result *Paste
				err := coll.Find(bson.M{"id": id}).One(&result)
				if err != nil {
				   w.WriteHeader(404)
				}

				if result == nil {
           w.WriteHeader(404)
				} else {
            temp,_ := template.ParseFiles("view.html")
					  temp.Execute(w, &result)
				}
	})

	m.NotFound(func(w http.ResponseWriter,r *http.Request) {
		 fmt.Fprintln(w, "404 nigga")
	})

	http.ListenAndServe(":80", m) // Starts the webserver
}

func GenHash() string {
	b := make([]rune, 5)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
