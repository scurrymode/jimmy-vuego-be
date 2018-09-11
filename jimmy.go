package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type ResWithData struct {
	Code int
	Data []Inventory
}
type Res struct {
	Code int
}

type Vins struct {
	Vins []string
}

type Inventory struct {
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
	filePath := "./inventory.json"
	r := mux.NewRouter()

	r.HandleFunc("/inventory", func(w http.ResponseWriter, r *http.Request) {
		inventories, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic(err)
		}

		var invs []Inventory
		json.Unmarshal(inventories, &invs)

		res := ResWithData{200, invs}
		stringRes, rerr := json.Marshal(res)
		if rerr != nil {
			panic(rerr)
		}
		w.Write(stringRes)
	}).Methods("GET")

	r.HandleFunc("/inventory/add", func(w http.ResponseWriter, r *http.Request) {
		inventories, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic(err)
		}

		fmt.Printf("inventories: %v\n", string(inventories))

		var invs []Inventory
		json.Unmarshal(inventories, &invs)

		fmt.Printf("len: %v", len(invs))

		inv := Inventory{
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

		invs = append(invs, inv)

		fmt.Printf("appended invs: %v\n", invs)

		stringBytes, jerr := json.Marshal(invs)
		if jerr != nil {
			panic(jerr)
		}

		//파일 쓰기
		err = ioutil.WriteFile(filePath, stringBytes, 0)
		if err != nil {
			panic(err)
		}

		res := Res{200}
		stringRes, rerr := json.Marshal(res)
		if rerr != nil {
			panic(rerr)
		}
		w.Write(stringRes)
	}).Methods("POST")

	r.HandleFunc("/inventory/delete", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var vins Vins
		err := decoder.Decode(&vins)
		if err != nil {
			panic(err)
		}

		// fmt.Println(vins)
		// fmt.Println(len(vins.Vins))

		//파일 읽기
		inventories, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic(err)
		}

		// fmt.Printf("inventories: %v\n", string(inventories))

		var invs []Inventory
		json.Unmarshal(inventories, &invs)

		// fmt.Printf("len before: %v", len(invs))

		n := 0
		for n < len(vins.Vins) {
			for inv := range invs {
				if invs[inv].Vin == string(vins.Vins[n]) {
					fmt.Printf("find : %v", string(vins.Vins[n]))
					copy(invs[inv:], invs[inv+1:])
					invs = invs[:len(invs)-1]
					break
				}
			}
			n++
		}

		stringBytes, jerr := json.Marshal(invs)
		if jerr != nil {
			panic(jerr)
		}

		err = ioutil.WriteFile(filePath, stringBytes, 0)
		if err != nil {
			panic(err)
		}

		res := Res{200}
		stringRes, rerr := json.Marshal(res)
		if rerr != nil {
			panic(rerr)
		}
		w.Write(stringRes)
	}).Methods("POST")

	//allowedOrigins
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	http.ListenAndServe(":8081", handler)
	// http.ListenAndServe(":8081", r)
}
