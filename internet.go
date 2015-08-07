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
	//person is a person with a name, and an age
	Name string
	Age int
}

type jsonResponse struct {
	//jsonResponse is the actual response we will send to the client.
	//Valid is false if we are returning no useful data, and true if we are
	Valid bool
	Data []person
}

func getPeople(name string) []person {

	//an array to hold the results we get from our query
	var people []person

	row, err := con.Query("select name,age from test where name = ?", name)
	checkErr(err)
	for row.Next() {
		//and empty person to hold the data from this row
		me := person{}
		err = row.Scan(&me.Name,&me.Age)
		checkErr(err)

		//append this person to the people array
		people = append(people, me)
	}

	//return all fo the results [which may be empty]
	return(people)
}

func userHandler(w http.ResponseWriter, r *http.Request) {

	//crafting the default response.  Create an empty person struct
	//then add it to an array, and set the response struct's data
	//container to that array [with the empty person in it]
	//there MUST be a better way of doing this...

	somePerson := person{}
	var somePeople []person
	somePeople = append(somePeople,somePerson)
	response := jsonResponse{Valid:false,Data:somePeople}

	//if we got a user in the uesr field, query the database
	if (r.FormValue("user") != "") {
		people := getPeople(r.FormValue("user"))
		if (len(people) > 0) {
			//if we got soem results, put them in the response
			response.Data = people
			response.Valid = true
		}
	}

	//turn our response struct into a json
	_response, err := json.Marshal(response)
	checkErr(err)

	//send it to the client
	fmt.Fprintf(w,string(_response))
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the default response")
}


func main() {
	con, err = sql.Open("mysql", "root:fruits@/golang")
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
