package main

import (
	"fmt"
	"net/http"
)

func userHandler(w http.ResponseWriter, r *http.Request) {
	if (r.FormValue("user") != "") {
		fmt.Fprintf(w, "Your username is %s", r.FormValue("user"))
	} else {
		defaultHandler(w, r)
	}
}

func itemHandler(w http.ResponseWriter, r *http.Request) {
	if (r.FormValue("item") != "") {
		fmt.Fprintf(w, "The item you requested is %s", r.FormValue("item"))
	} else {
		defaultHandler(w, r)
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the default response")
}


func main() {
	http.HandleFunc("/item/", itemHandler)
	http.HandleFunc("/user/", userHandler)
	http.HandleFunc("/", defaultHandler)
	http.ListenAndServe(":8080", nil)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
