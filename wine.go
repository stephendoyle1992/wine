package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type Wine struct {
	ID          int           `db:"id" json:"id"`
	Country     string        `db:"country" json:"country"`
	Description string        `db:"description" json:"description"`
	Designation string        `db:"designation" json:"designation"`
	Points      sql.NullInt64 `db:"points" json:"points"`
	Price       sql.NullInt64 `db:"price" json:"price"`
	Province    string        `db:"province" json:"province"`
	Region1     string        `db:"region_1" json:"region1"`
	Region2     string        `db:"region_2" json:"region2"`
	Title       string        `db:"title" json:"title"`
	Variety     string        `db:"variety" json:"variety"`
	Winery      string        `db:"winery" json:"winery"`
}

var Db *sqlx.DB

func main() {
	var err error
	Db, err = sqlx.Open("mysql", "root@tcp(127.0.0.1:3306)/WineApp")
	if err != nil {
		panic(err)
	}
	defer Db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/api/countries/", getCountryList).Methods("GET")
	r.HandleFunc("/api/wines/", getWines).Methods("GET")

	err = http.ListenAndServe(":8888", r)
	fmt.Println(err)

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

func getWines(w http.ResponseWriter, r *http.Request) {
	qvals := r.URL.Query()

	var args []interface{}

	q := `SELECT * FROM wine WHERE PRICE <= ?`

	if qvals["type"] == nil {
		fmt.Println("type missing")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if qvals["price"] != nil {
		args = append(args, qvals["price"][0])
	} else {
		args = append(args, 10000000000000)
	}

	if qvals["type"][0] == "white" {
		q += ` AND (variety LIKE 'sauvignon blanc' OR variety
		LIKE 'verdelho' OR variety LIKE 'semillon' OR variety LIKE 'chardonnay' OR
		variety LIKE 'riesling' OR variety LIKE 'pinot gris' OR variety LIKE
		'pinot grigio' or variety LIKE 'white blend')`
	} else {
		q += ` AND (variety LIKE 'carbenet sauvignon' OR variety LIKE 'shiraz'
		OR variety LIKE 'merlot' OR variety LIKE  'pinot noir' 
		OR variety LIKE 'grenache' OR variety LIKE 'red blend')`
	}

	if qvals["country"] != nil {
		if qvals["country"][0] != "any" {
			q += ` AND COUNTRY = ?`
			args = append(args, qvals["country"][0])
		}
	}
	if qvals["region"] != nil {
		if qvals["region"][0] != "any" {
			q += ` AND REGION = ?`
			args = append(args, qvals["region"][0])
		}
	}
	if qvals["variety"] != nil {
		if qvals["variety"][0] != "any" {
			q += ` AND REGION = ?`
			args = append(args, qvals["variety"][0])
		}
	}

	wines := []Wine{}

	if err := Db.Select(&wines, q, args...); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(wines)
}
