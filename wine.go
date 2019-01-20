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
	Region1     string        `db:"region1" json:"region1"`
	Region2     string        `db:"region2" json:"region2"`
	Title       string        `db:"title" json:"title"`
	Variety     string        `db:"variety" json:"variety"`
	Winery      string        `db:"winery" json:"winery"`
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
	r.HandleFunc("/api/wines/", getWines).Methods("GET")
	r.HandleFunc("/api/variety/", getVarietyList).Methods("GET")
	r.HandleFunc("/api/{countries}/region1/", getRegion1).Methods("GET")

	err = http.ListenAndServe(":8888", r)
	fmt.Println(err)
}

func getVarietyList(w http.ResponseWriter, r *http.Request) {
	q := "select distinct variety from Wine"

	varieties := []Wine{}

	if err := Db.Select(&varieties, q); err != nil {
		fmt.Println(err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(varieties)

}
func getRegion1(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	countries := vars["countries"]
	fmt.Println(countries)
	regions := []Wine{}
	q := `SELECT distinct region1 FROM Wine WHERE country=?`

	if err := Db.Select(&regions, q, countries); err != nil {
		fmt.Println(err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(regions)
}
func getCountryList(w http.ResponseWriter, r *http.Request) {
	q := "select distinct country from Wine"

	countries := []Wine{}

	if err := Db.Select(&countries, q); err != nil {
		fmt.Println(err)

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
		q += ` AND (variety LIKE 'sauvignon blanc' OR variety LIKE 'verdelho' 
		OR variety LIKE 'semillon' OR variety LIKE 'chardonnay' 
		OR variety LIKE 'riesling' OR variety LIKE 'pinot gris' 
		OR variety LIKE	'pinot grigio' or variety LIKE 'white blend'
		 or variety LIKE 'Moscato'or variety LIKE 'Muscat' or variety LIKE 'semillion'
		  or variety LIKE 'viognier' or variety like'Gewürztraminer' or variety like 'Pinot Bianco' )`
	} else {
		q += ` AND variety not in (select variety from Wine where variety LIKE 'sauvignon blanc' OR variety LIKE 'verdelho' 
		OR variety LIKE 'semillon' OR variety LIKE 'chardonnay' 
		OR variety LIKE 'riesling' OR variety LIKE 'pinot gris' 
		OR variety LIKE	'pinot grigio' or variety LIKE 'white blend'
		 or variety LIKE 'Moscato'or variety LIKE 'Muscat' or variety LIKE 'semillion'
		  or variety LIKE 'viognier' or variety like'Gewürztraminer' or variety like 'Pinot Bianco' )`
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
