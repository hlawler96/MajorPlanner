package main
import (
    "log"
    "database/sql"
    "strings"
)

type PreReq struct {
  Crs   []Course    `json:"Courses"`
  Des   string      `json:"Type"`
}
//gets prereq for strictly required courses
func getLoosePrereqs(courses []LooseReqCourse, db *sql.DB) []PreReq {
  prereqs := make([]PreReq, 0)

  for _ , c := range courses {
    str := " SELECT distinct PR.prid , C.creditHours, C.cNumber, C.dept, PR.req from Prereqs PR, Courses C WHERE PR.cid = ? and C.id = PR.prid ORDER BY PR.req"
    rows, err := db.Query(str,c.ReqCourse.Id)
    switch {
      case err == sql.ErrNoRows:
        //do nothing
      case err != nil:
        log.Fatal(err)
      default:
        prereq := new(PreReq)
        for rows.Next() {
          temp := ""
          crs := new(Course)
          err := rows.Scan(&crs.Id, &crs.Hours, &crs.Number, &crs.Program, &temp)
          if err != nil {
            log.Fatal(err)
          }
          if strings.Compare(temp, prereq.Des) == 0 {
            prereq.Crs = append(prereq.Crs, *crs)
          }else if prereq.Crs == nil{
              prereq.Des = temp
              prereq.Crs = append(prereq.Crs , c.ReqCourse)
              prereq.Crs = append(prereq.Crs, *crs)
          }else {
              prereqs = append(prereqs, *prereq)
              prereq = new(PreReq)
              prereq.Crs = append(prereq.Crs , c.ReqCourse)
              prereq.Crs = append(prereq.Crs, *crs)
              prereq.Des = temp
          }
        }

        if prereq.Crs != nil {
        prereqs = append(prereqs, *prereq)
        }


    defer rows.Close()
  }
}
  return prereqs
}
//gets prereq for strictly required courses
func getStrictPrereqs(courses []Course, db *sql.DB) []PreReq {
  prereqs := make([]PreReq, 0)

  for _ , c := range courses {
    str := " SELECT distinct PR.prid , C.creditHours, C.cNumber, C.dept, PR.req from Prereqs PR, Courses C WHERE PR.cid = ? and C.id = PR.prid ORDER BY PR.req"
    rows, err := db.Query(str,c.Id)
    switch {
      case err == sql.ErrNoRows:
        //do nothing
      case err != nil:
        log.Fatal(err)
      default:
        prereq := new(PreReq)
        for rows.Next() {
          temp := ""
          crs := new(Course)
          err := rows.Scan(&crs.Id, &crs.Hours, &crs.Number, &crs.Program, &temp)
          if err != nil {
            log.Fatal(err)
          }
          if strings.Compare(temp, prereq.Des) == 0 {
            prereq.Crs = append(prereq.Crs, *crs)
          }else if prereq.Crs == nil{
              prereq.Des = temp
              prereq.Crs = append(prereq.Crs , c)
              prereq.Crs = append(prereq.Crs, *crs)
          }else {
              prereqs = append(prereqs, *prereq)
              prereq = new(PreReq)
              prereq.Crs = append(prereq.Crs , c)
              prereq.Crs = append(prereq.Crs, *crs)
              prereq.Des = temp
          }

        }
        if prereq.Crs != nil {
        prereqs = append(prereqs, *prereq)
        }

      }
    defer rows.Close()
  }
  return prereqs
}
