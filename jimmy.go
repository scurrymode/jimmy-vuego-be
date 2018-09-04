// package main
// import (
//   "net/http"
//   "strings"
//   "github.com/gorilla/mux"
// )
// func sayHello(w http.ResponseWriter, r *http.Request) {
//   message := r.URL.Path
//   message = strings.TrimPrefix(message, "/")
//   message = "Hello " + message
//   w.Write([]byte(message))
// }
// func main() {
//   http.HandleFunc("/", sayHello)
//   if err := http.ListenAndServe(":8080", nil); err != nil {
//     panic(err)
//   }
// }
package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/jimmy/{title}/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
	})

	http.ListenAndServe(":8080", r)
}