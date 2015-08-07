package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var con *sql.DB
var err error

type person struct {
	Name string
	Age int
}

type jsonResponse struct {
	Valid bool
	Data []person
}

func getPeople(name string) []person {

	var people []person

	row, err := con.Query("select name,age from test where name = ?", name)
	checkErr(err)
	for row.Next() {
		me := person{}
		err = row.Scan(&me.Name,&me.Age)
		checkErr(err)
		people = append(people, me)
	}

	return(people)

}

func userHandler(w http.ResponseWriter, r *http.Request) {

	somePerson := person{}
	var somePeople []person
	somePeople = append(somePeople,somePerson)
	response := jsonResponse{Valid:false,Data:somePeople}

	if (r.FormValue("user") != "") {
		people := getPeople(r.FormValue("user"))
		if (len(people) > 0) {
			response.Data = people
			response.Valid = true
		}
	}
	_response, err := json.Marshal(response)
	checkErr(err)
	fmt.Fprintf(w,string(_response))
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the default response")
}


func main() {
	con, err = sql.Open("mysql", "root:veggies@/golang")
	checkErr(err)
	http.HandleFunc("/user/", userHandler)
	http.HandleFunc("/", defaultHandler)
	http.ListenAndServe(":8080", nil)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
