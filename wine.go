package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	//_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Wine struct {
	ID          int            `db:"id" json:"id"`
	Country     sql.NullString `db:"country" json:"country"`
	Description sql.NullString `db:"description" json:"description"`
	Designation sql.NullString `db:"designation" json:"designation"`
	Points      sql.NullInt64  `db:"points" json:"points"`
	Price       sql.NullInt64  `db:"price" json:"price"`
	Province    sql.NullString `db:"province" json:"province"`
	Region1     sql.NullString `db:"region1" json:"region1"`
	Region2     sql.NullString `db:"region2" json:"region2"`
	Title       sql.NullString `db:"title" json:"title"`
	Variety     sql.NullString `db:"variety" json:"variety"`
	Winery      sql.NullString `db:"winery" json:"winery"`
}

var Db *sqlx.DB

func main() {
	var err error

	//Db, err = sqlx.Open("mysql", "root@tcp(127.0.0.1:3306)/WineApp")
	Db, err = sqlx.Open("postgres", os.Getenv("DATABASE_URL"))
	//Db, err = sqlx.Open("postgres", "user=postgres dbname=postgres port=5432 sslmode=disable")

	if err != nil {
		panic(err)
	}
	defer Db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/api/countries/", getCountryList).Methods("GET")
	r.HandleFunc("/api/wines/", getWines).Methods("GET")
	r.HandleFunc("/api/variety/", getVarietyList).Methods("GET")
	r.HandleFunc("/api/{countries}/region1/", getRegion1).Methods("GET")

	err = http.ListenAndServe(":"+os.Getenv("PORT"), r)
	//err = http.ListenAndServe(":8888", r)
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
	q := `SELECT distinct region1 FROM Wine WHERE country=($1)`

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
	fieldVal := 2

	q := `SELECT * FROM wine WHERE PRICE <= ($1)`

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
			q += ` AND COUNTRY = ($` + strconv.Itoa(fieldVal) + `)`
			args = append(args, qvals["country"][0])
			fieldVal++
		}
	}
	if qvals["region"] != nil {
		if qvals["region"][0] != "any" {
			q += ` AND REGION = ($` + strconv.Itoa(fieldVal) + `)`
			args = append(args, qvals["region"][0])
			fieldVal++
		}
	}
	if qvals["variety"] != nil {
		if qvals["variety"][0] != "any" {
			q += ` AND REGION = ($` + strconv.Itoa(fieldVal) + `)`
			args = append(args, qvals["variety"][0])
			fieldVal++
		}
	}

	if qvals["status"] != nil {
		if qvals["status"][0] != "any" {
			if qvals["status"][0] == "value" {
				q += ` ORDER BY (points+1/price+1)`
			} else if qvals["status"][0] == "points" {
				q += ` ORDER BY points` 
			} else if qvals["status"][0] == "cheap" {
				q += ` ORDER BY price`
			}
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
