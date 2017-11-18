package main

import (
    "log"
    "encoding/json"
    "net/http"
    "database/sql"
)
import _ "github.com/go-sql-driver/mysql"

func Index(w http.ResponseWriter, r *http.Request) {
  //Tell Response to expect a json
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)

  //open database using root
  //db is currently on my local, not sure how to get around currently
  db, err := sql.Open("mysql", "root@/planner")
 if err != nil {
   log.Fatal(err)
 }

 rows, err := db.Query("SELECT * FROM courses")
 if err != nil {
   log.Fatal(err)
 }
 defer rows.Close()

 courses := make([]*Course, 0)
  for rows.Next() {
    crs := new(Course)
    err := rows.Scan(&crs.Id, &crs.Hours, &crs.Dept, &crs.Number, &crs.Pid)
    if err != nil {
      log.Fatal(err)
    }
    courses = append(courses, crs)
  }
  if err = rows.Err(); err != nil {
    log.Fatal(err)
  }

 if err := json.NewEncoder(w).Encode(courses); err != nil {
       panic(err)
   }
}
