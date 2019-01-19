package main

import (
  "database/sql"
  "fmt"
  _ "github.com/go-sql-driver/mysql"
)

type Wine struct {
  Id int
  Country string
  Description string
  Designation string
  Points int
  Price int
  Province string
  Region1 string
  Region2 string
  Title string
  Variety string
  Winery string
}

func main() {
  db, err := sql.Open("mysql", "jordan:1234@tcp(127.0.0.1:3306)/WineApp")
  if err != nil {
    panic(err)
  }
  defer db.Close()

  rows, err := db.Query("SELECT * FROM wine WHERE country='Canada'")
  if err != nil {
    panic(err)
  }
  defer rows.Close()

  for rows.Next() {
    wine := Wine{}
    err = rows.Scan(&wine.Id,
      &wine.Country,
      &wine.Description,
      &wine.Designation,
      &wine.Points,
      &wine.Price,
      &wine.Region1,
      &wine.Region2,
      &wine.Title,
      &wine.Variety,
      &wine.Winery)
    if err != nil {
      panic(err)
    }
    fmt.Println(wine)
  }
  err = rows.Err()
  if err != nil {
    panic(err)
  }

}
