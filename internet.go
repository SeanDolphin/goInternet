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

func getPerson(name string) person {

	me := person{}

	row, err := con.Query("select name,age from test where name = ?", name)
	checkErr(err)
	for row.Next() {
		err = row.Scan(&me.Name,&me.Age)
		checkErr(err)
	}

	return(me)

}

func userHandler(w http.ResponseWriter, r *http.Request) {

	if (r.FormValue("user") != "") {
		me := getPerson(r.FormValue("user"))
		peopleJson, err := json.Marshal(me)
		checkErr(err)
		fmt.Fprintf(w, string(peopleJson))
	} else {
		defaultHandler(w, r)
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the default response")
}


func main() {
	con, err = sql.Open("mysql", "root:goblots@/golang")
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
