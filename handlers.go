package main

import (
    "log"
    "encoding/json"
    "net/http"
    "database/sql"
    "fmt"
)
import _ "github.com/go-sql-driver/mysql"

func getCourses(w http.ResponseWriter, r *http.Request) {
  //Tell Response to expect a json
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)

  //open database using root
  //db is currently on my local, not sure how to get around currently
  // db, err := sql.Open("mysql", "root@/planner")
  db, err := sql.Open("mysql", "mason:pineappleB2@tcp(comp426finalproject.cqu5t9sfyvwq.us-east-2.rds.amazonaws.com:3306)/planner" )
 if err != nil {
   log.Fatal(err)
 }

 rows, err := db.Query("SELECT * FROM Courses")
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

func test(w http.ResponseWriter, r *http.Request) {
fmt.Fprintln(w, "Hello world :)")

}

func testLogin(w http.ResponseWriter, r *http.Request){
  //example call
  // http://localhost:8080/Login/?username=user&password=pass

  //TODO Actually Query Database to get user info
  //TODO Maybe send back a set Cookie to allow for pesistance in login session

  user, ok := r.URL.Query()["username"]
  if !ok || len(user) < 1 {
       log.Println("Url Param 'username' is missing")
       return
   }
  pass, ok := r.URL.Query()["password"]
  if !ok || len(pass) < 1 {
       log.Println("Url Param 'password' is missing")
       return
   }


   type LoginResponse struct {
       Id      int     `json:"id"`
   }
   res := LoginResponse{
     Id: 1,
   }
  //return the user information for the user who just logged in
  if err := json.NewEncoder(w).Encode(res); err != nil {
        panic(err)
    }
}

func testSignUp(w http.ResponseWriter, r *http.Request){
  // example call
  //http://localhost:8080/SignUp/?name=hayden&username=hayden&password=password

  name, ok := r.URL.Query()["name"]
  if !ok || len(name) < 1 {
      fmt.Fprintln(w, "Url Param 'name' is missing")
       return
   }

  user, ok := r.URL.Query()["username"]
  if !ok || len(user) < 1 {
       fmt.Fprintln(w, "Url Param 'username' is missing")
       return
   }
  pass, ok := r.URL.Query()["password"]
  if !ok || len(pass) < 1 {
       fmt.Fprintln(w, "Url Param 'password' is missing")
       return
   }
   //TODO Actually insert new user into DB and return the user id in the json response
   type SignUpResponse struct {
       Id      int     `json:"id"`
   }
   res := SignUpResponse{
     Id: 1,
   }

   //return the userId of the person who signed up
   if err := json.NewEncoder(w).Encode(res); err != nil {
         panic(err)
     }

}

func testGetCoursesTaken(w http.ResponseWriter, r *http.Request){
  //example call http://localhost:8080/CoursesTaken/?id=1
  //if you give the method an id of 10 it will return empty
  // as if it was a new user

  id, ok := r.URL.Query()["id"]
  if !ok || len(id) < 1 {
       fmt.Fprintln(w, "Url Param 'id' is missing")
       return
   }

    courses := Courses {}
    if id[0] != "10" {

    //TODO Actually return courses user has taken
     courses = Courses{
     Course {"1", "3", "COMP", "110", "1",},
     Course {"2", "3", "COMP", "401", "1",},
     Course {"3", "3", "COMP", "410", "1",},
     Course {"4", "3", "COMP", "411", "1",},
     Course {"5", "3", "MATH", "233", "2",},
   }
 }

   if err := json.NewEncoder(w).Encode(courses); err != nil {
         panic(err)
     }

}

func testPostUserInformation(w http.ResponseWriter, r *http.Request){


  /*
  Sample Post Message Body Request
  {"id":1, "deptTaken":[{"name":"COMP", "coursesTaken": [{"id":"1","hours":"3","dept":"COMP","number":"110","pid":"1"},{"id":"2","hours":"3","dept":"COMP","number":"401","pid":"1"}]}, {"name":"MATH", "coursesTaken": [{"id":"1","hours":"3","dept":"MATH","number":"233","pid":"2"}]}],"currDept":["COMP", "MATH"], "semLeft": 4, "genEdsLeft": 3}
  */

  type DeptTaken struct{
    Name            string    `json:"name"`
    CoursesTaken    Courses   `json:"coursesTaken"`
  }

  type DeptsTaken []DeptTaken

  type UserInfo struct {
    Id         int           `json:"id"`
    DTaken     DeptsTaken    `json:"deptTaken"`
    CurrDept   []string      `json:"currDept"`
    SemLeft    int           `json:"semLeft"`
    GenEdsLeft int           `json:"genEdsLeft"`
  }

  decoder := json.NewDecoder(r.Body)

	var user_info UserInfo
	err := decoder.Decode(&user_info)

  if err != nil {
		panic(err)
	}

//TODO Actually store all of this information in the DB instead of just printing it

  fmt.Println(user_info.Id)
  for _, dept := range user_info.DTaken {
     fmt.Println(dept.Name)
     for _, course := range dept.CoursesTaken {
          fmt.Println(string(course.Id) + " " + string(course.Dept) + " " + string(course.Number) + " " + string(course.Hours) + " " + string(course.Pid))
     }
  }
  for _, dept := range user_info.CurrDept {
    fmt.Println(dept)
  }
  fmt.Println(user_info.SemLeft)
  fmt.Println(user_info.GenEdsLeft)

}


func testGetResult(w http.ResponseWriter, r *http.Request){
  id, ok := r.URL.Query()["id"]
  if !ok || len(id) < 1 {
       fmt.Fprintln(w, "Url Param 'id' is missing")
       return
   }


}
