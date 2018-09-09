package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

type Res struct {
    Code   int
    Data   []Inventory
}
type Res2 struct {
	Code int
}

type Inventory struct {
    No     string `json:"no"`
    Vin    string `json:"vin"`
	Model  string `json:"model"`
	Make   string `json:"make"`
	Year   string `json:"year"`
	Msrp   string `json:"msrp"`
	Status string `json:"status"`
	Booked string `json:"booked"`
	Listed string `json:"listed"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/jimmy/{title}/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
	})

	r.HandleFunc("/inventory", func(w http.ResponseWriter, r *http.Request) {
		//파일 읽기
		inventories, err := ioutil.ReadFile("./inventory.json")
		if err != nil {
			panic(err)
		}

		var invs []Inventory
		json.Unmarshal( inventories, &invs )

		// jsonBytes, jerr := json.Marshal( invs )
		// if jerr != nil {
		// 	panic(jerr)
		// }

		res := Res{200, invs}
		stringRes, rerr := json.Marshal(res)
		if rerr != nil {
			panic(rerr)
		}
		w.Write( stringRes )
	})

	r.HandleFunc("/inventory/add", func(w http.ResponseWriter, r *http.Request) {
		//파일 읽기
		inventories, err := ioutil.ReadFile("./inventory.json")
		if err != nil {
			panic(err)
		}

		fmt.Printf("inventories: %v\n", string(inventories))

		var invs []Inventory
		json.Unmarshal( inventories, &invs )

		fmt.Printf("len: %v", len(invs))

		inv := Inventory{
			strconv.Itoa(len(invs)+1),
		 	r.FormValue("vin"),
			r.FormValue("model"),
			r.FormValue("make"),
			r.FormValue("year"),
			r.FormValue("msrp"),
			r.FormValue("status"),
			r.FormValue("booked"),
			r.FormValue("listed")}

		fmt.Printf("inv: %v\n", inv)

		fmt.Printf("invs: %v\n", invs)

		invs = append( invs, inv )

		fmt.Printf("appended invs: %v\n", invs)

		stringBytes, jerr := json.Marshal( invs )
		if jerr != nil {
			panic(jerr)
		}

		//파일 쓰기
		err = ioutil.WriteFile("./inventory.json", stringBytes, 0)
		if err != nil {
			panic(err)
		}

		res := Res2{200}
		stringRes, rerr := json.Marshal(res)
		if rerr != nil {
			panic(rerr)
		}
		w.Write(stringRes);
	}).Methods("POST")

	r.HandleFunc("/inventory/delete/{no}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		no := vars["no"]
		
		//파일 읽기
		inventories, err := ioutil.ReadFile("./inventory.json")
		if err != nil {
			panic(err)
		}

		fmt.Printf("inventories: %v\n", string(inventories))

		var invs []Inventory
		json.Unmarshal( inventories, &invs )

		fmt.Printf("len before: %v", len(invs))

		i, err := strconv.Atoi(no)
		invs = invs[:i]

		fmt.Printf("len after: %v", len(invs))

		fmt.Printf("deleted invs: %v\n", invs)

		stringBytes, jerr := json.Marshal( invs )
		if jerr != nil {
			panic(jerr)
		}

		//파일 쓰기
		err = ioutil.WriteFile("./inventory.json", stringBytes, 0)
		if err != nil {
			panic(err)
		}

		res := Res2{200}
		stringRes, rerr := json.Marshal(res)
		if rerr != nil {
			panic(rerr)
		}
		w.Write(stringRes);
	}).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8080"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	http.ListenAndServe(":8081", handler )
	// http.ListenAndServe(":8081", r)
}