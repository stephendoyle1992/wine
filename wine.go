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
  TasterName string
  TasterTwitterHandle string
  Title string
  Variety string
  Winery string
}

func main() {
  db, err := sql.Open("mysql", "theUser: thePassword@/theDbName")
  if err != nil {
    panic(err)
  }
  defer db.Close()

  rows, err := db.Query("SELECT ...")
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
      &wine.TasterName,
      &wine.TasterTwitterHandle,
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
