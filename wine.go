package main

import (
	"encoding/json"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type Wine struct {
	ID          int    `db:"id" json:"id"`
	Country     string `db:"country" json:"country"`
	Description string `db:"description" json:"description"`
	Designation string `db:"designation" json:"designation"`
	Points      int    `db:"points" json:"points"`
	Price       int    `db:"price" json:"price"`
	Province    string `db:"province" json:"province"`
	Region1     string `db:"region1" json:"region1"`
	Region2     string `db:"region2" json:"region2"`
	Title       string `db:"title" json:"title"`
	Variety     string `db:"variety" json:"variety"`
	Winery      string `db:"winery" json:"winery"`
}

var Db *sqlx.DB

func main() {
	var err error
	Db, err = sqlx.Open("mysql", "jordan:1234@tcp(127.0.0.1:3306)/WineApp")
	if err != nil {
		panic(err)
	}
	defer Db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/api/countries/", getCountryList).Methods("GET")

	http.ListenAndServe("8888", r)

}

func getCountryList(w http.ResponseWriter, r *http.Request) {
	q := `PUT THE DB CALL IN HERE`

	countries := []Wine{}

	if err := Db.Select(&countries, q); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(countries)
}
