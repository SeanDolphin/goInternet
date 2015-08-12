package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var con *sql.DB
var r *render.Render

var ErrNoPeople = errors.New("")

type person struct {
	//person is a person with a name, and an age
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type jsonResponse struct {
	//jsonResponse is the actual response we will send to the client.
	//Valid is false if we are returning no useful data, and true if we are
	Valid bool     `json:"valid"`
	Data  []person `json:"data"`
}

func getPeople(name string) ([]person, error) {

	//an array to hold the results we get from our query
	var people []person

	row, err := con.Query("select name,age from test where name = ?", name)
	if err != nil {
		panic(err)
	}

	for row.Next() {
		//and empty person to hold the data from this row
		me := person{}
		if err := row.Scan(&me.Name, &me.Age); err != nil {
			return people, err
		}

		//append this person to the people array
		people = append(people, me)
	}

	if len(people) == 0 {
		return people, ErrNoPeople
	}

	//return all fo the results [which may be empty]
	return people, nil
}

func userHandler(writer http.ResponseWriter, req *http.Request) {

	//crafting the default response.  Create an empty person struct
	//then add it to an array, and set the response struct's data
	//container to that array [with the empty person in it]
	//there MUST be a better way of doing this...

	response := jsonResponse{Valid: false, Data: []person{}}

	//if we got a user in the uesr field, query the database
	if userName := req.FormValue("user"); userName != "" {

		if people, err := getPeople(userName); err == nil {
			//if we got soem results, put them in the response
			response.Data = people
			response.Valid = true
		}
	}

	//turn our response struct into a json
	r.JSON(writer, http.StatusOK, response)
}

func defaultHandler(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(writer, "This is the default response")
}

func main() {
	var err error
	con, err = sql.Open("mysql", "root:fruits@/golang")
	if err != nil {
		log.Fatal(err)
	}

	r = render.New(render.Options{})

	mux := mux.NewRouter()
	mux.HandleFunc("/user", userHandler).Methods("POST")
	mux.HandleFunc("/", defaultHandler).Methods("POST")

	PORT := ":" + os.Getenv("PORT")
	if PORT == ":" {
		PORT = ":8080"
	}

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(PORT)
}
