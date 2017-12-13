package main

import (
    "log"
    "encoding/json"
    "net/http"
    "database/sql"
    "fmt"
    "time"
    "math/rand"
)
import _ "github.com/go-sql-driver/mysql"
//fully functioning
func getCourses(w http.ResponseWriter, r *http.Request) {
  //Tell Response to expect a json
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")

  w.WriteHeader(http.StatusOK)

  //open database
  db, err := sql.Open("mysql", "mason:pineappleB2@tcp(comp426finalproject.cqu5t9sfyvwq.us-east-2.rds.amazonaws.com:3306)/planner" )
 if err != nil {
   log.Fatal(err)
 }
 sql := "select C.id, C.creditHours, C.cNumber, C.dept from Courses C"
 //check for dept param
 dept, ok1 := r.URL.Query()["dept"]

 //if dept param is given then only return classes in that dept
 if !ok1 || len(dept) < 1 {

  }

 Ptype, ok := r.URL.Query()["type"]

 if !ok || len(Ptype) < 1 {

  }else{
    sql = "SELECT C.id, C.creditHours, C.cNumber, C.dept FROM Courses C, Program P, ProgramRequirements PR, CoursesInProgram CP where P.id = PR.pid and CP.prid = PR.id and C.id = CP.cid and P.dept = '" + dept[0] + "' and P.type = '" + Ptype[0] + "' "
  }

  rows, err := db.Query(sql)
  if err == sql.ErrNoRows {
    log.Println("empty sql result")
  }else if err != nil {
    log.Fatal(err)
  }
  defer rows.Close()


 courses := make([]Course, 0)
  for rows.Next() {
    crs := new(Course)
    err := rows.Scan(&crs.Id, &crs.Hours, &crs.Number, &crs.Program)
    if err != nil {
      log.Fatal(err)
    }
    courses = append(courses, *crs)
  }
  if err = rows.Err(); err != nil {
    log.Fatal(err)
  }

 if err := json.NewEncoder(w).Encode(courses); err != nil {
       panic(err)
   }
}
//fully functioning lol
func test(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
fmt.Fprintln(w, "Hello world :)")
}
//fully functioning
func Login(w http.ResponseWriter, r *http.Request){
  // example call
  // http://localhost:8080/Login/?username=user&password=pass
 w.Header().Set("Access-Control-Allow-Origin", "*")


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
       SessionId      string     `json:"sessionId"`

   }
   res := new(LoginResponse)

   db, err := sql.Open("mysql", "mason:pineappleB2@tcp(comp426finalproject.cqu5t9sfyvwq.us-east-2.rds.amazonaws.com:3306)/planner" )
  if err != nil {
    log.Fatal(err)
  }

  // Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO UserSessions VALUES (?,?) ON DUPLICATE KEY UPDATE sessionId = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()

  var id int
  err = db.QueryRow("SELECT U.id FROM Users U WHERE U.user = ? and U.pass = ?", user[0], pass[0]).Scan(&id)
  switch {
  	case err == sql.ErrNoRows:
  		res.SessionId = ""
  	case err != nil:
  		log.Fatal(err)
  	default:
      res.SessionId = RandStringGenerator(30)
      stmtIns.Exec(id, res.SessionId, res.SessionId)
  	}
  //return the user information for the user who just logged in
  if err := json.NewEncoder(w).Encode(res); err != nil {
        panic(err)
    }
}
//fully functioning
func SignUp(w http.ResponseWriter, r *http.Request){
  // example call
  //http://localhost:8080/SignUp/?name=hayden&username=hayden&password=password
w.Header().Set("Access-Control-Allow-Origin", "*")

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

   type SignUpResponse struct {
       SessionId      string     `json:"sessionId"`
   }
   res := new(SignUpResponse)

   //connect to db
   db, err := sql.Open("mysql", "mason:pineappleB2@tcp(comp426finalproject.cqu5t9sfyvwq.us-east-2.rds.amazonaws.com:3306)/planner" )
  if err != nil {
    log.Fatal(err)
  }

  // Prepare statement for inserting data
	stmtInsUser, err := db.Prepare("INSERT INTO Users (user, pass) VALUES (?,?)")

	if err != nil {
		panic(err.Error())
	}
	defer stmtInsUser.Close()

  // Prepare statement for inserting data
  stmtInsSession, err := db.Prepare("INSERT INTO UserSessions VALUES (?,?) ON DUPLICATE KEY UPDATE sessionId = ?")
  if err != nil {
    panic(err.Error())
  }
  defer stmtInsSession.Close()

  var id int
  err = db.QueryRow("SELECT U.id FROM Users U WHERE U.user = ?", user[0]).Scan(&id)
  switch  {
    //username isn't taken, insert info into db
  case err == sql.ErrNoRows:
      _, err = stmtInsUser.Exec(user[0],pass[0])

      err = db.QueryRow("SELECT U.id FROM Users U WHERE U.user = ?", user[0]).Scan(&id)
      if err == sql.ErrNoRows {
        log.Println("empty sql result")
      }else if err != nil {
        log.Fatal(err)
      }
      res.SessionId = RandStringGenerator(30)
      _, err = stmtInsSession.Exec(id, res.SessionId)
      if err := json.NewEncoder(w).Encode(res); err != nil {
            panic(err)
        }
  //error in sql query
  case err != nil:
    log.Fatal(err)

  //return 0 if the user is already taken
  default:
    //return the userId of the person who signed up
    res.SessionId = ""
    if err := json.NewEncoder(w).Encode(res); err != nil {
          panic(err)
      }
  }
}
//fully functioning
func GetCoursesTaken(w http.ResponseWriter, r *http.Request){
  //example call http://localhost:8080/CoursesTaken/?id=1
  w.Header().Set("Access-Control-Allow-Origin", "*")

  sessionId, ok := r.URL.Query()["sessionId"]
  if !ok || len(sessionId) < 1 {
       fmt.Fprintln(w, "Url Param 'sessionId' is missing")
       return
   }
   //connect to db
   db, err := sql.Open("mysql", "mason:pineappleB2@tcp(comp426finalproject.cqu5t9sfyvwq.us-east-2.rds.amazonaws.com:3306)/planner" )
  if err != nil {
    log.Fatal(err)
  }

   //get courses that a specific user has taken
   str := "SELECT C.id, C.creditHours, C.cNumber, C.dept FROM Courses C, CoursesTaken CT, UserSessions US WHERE C.id = CT.cid and CT.uid = US.uid and US.sessionId = ?"
   rows, err := db.Query(str, sessionId[0])
   if err == sql.ErrNoRows {
     log.Println("empty sql result")
   }else if err != nil {
     log.Fatal(err)
   }
   defer rows.Close()

   //go through sql result and add them to list
   courses := make([]Course, 0)
    for rows.Next() {
      crs := new(Course)
      err := rows.Scan(&crs.Id, &crs.Hours, &crs.Number, &crs.Program)
      if err != nil {
        log.Fatal(err)
      }
      courses = append(courses, *crs)
    }
    if err = rows.Err(); err != nil {
      log.Fatal(err)
    }

    //return list of courses
   if err := json.NewEncoder(w).Encode(courses); err != nil {
         panic(err)
     }
}
//fully functioning
func PostUserInformation(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Access-Control-Allow-Origin", "*")
  /*
  {"sessionId":"jwGzoQQUmGmONbpqnDBPJeOrncVHbv",
 "deptTaken":[ {"name":"COMP", "coursesTaken": [{"dept":"COMP","number":110}, { "dept":"COMP","number":401}] } ,
			   {"name":"MATH", "coursesTaken": [{"dept":"MATH","number":233} ] } ],
 "currDept":[{"name":"COMP","type":"BS"},{"name":"MATH","type":"Minor"}],
 "semLeft": 4,
 "genEdsLeft": 3
 }
  */
  type AbvCourse struct {
    Dept    string    `json:"dept"`
    Number  int       `json:"number"`
  }
  type AbvCourses []AbvCourse

  type DeptTaken struct{
    Name            string       `json:"name"`
    CoursesTaken    AbvCourses   `json:"coursesTaken"`
  }

  type DeptsTaken []DeptTaken

  type Program  struct{
    Name     string   `json:"name"`
    Type     string   `json:"type"`
  }

  type UserInfo struct {
    SessionId  string        `json:"sessionId"`
    DTaken     DeptsTaken    `json:"deptTaken"`
    CurrDept   []Program     `json:"currDept"`
    SemLeft    int           `json:"semLeft"`
    GenEdsLeft int           `json:"genEdsLeft"`
  }

  decoder := json.NewDecoder(r.Body)

	var user_info UserInfo
	err := decoder.Decode(&user_info)

  if err != nil {
		panic(err)
	}

  db, err := sql.Open("mysql", "mason:pineappleB2@tcp(comp426finalproject.cqu5t9sfyvwq.us-east-2.rds.amazonaws.com:3306)/planner" )
  if err != nil {
   log.Fatal(err)
  }

  var id int
  err = db.QueryRow("SELECT U.id FROM Users U, UserSessions US WHERE US.uid = U.id and US.sessionId = ?", user_info.SessionId).Scan(&id)

  if err == sql.ErrNoRows {
    log.Println("empty sql result")
  }else if err != nil {
    log.Fatal(err)
  }

      // Prepare statements for inserting data
      stmtInsUser, err := db.Prepare("UPDATE Users SET Users.semLeft = ?, Users.genEdsLeft = ?, Users.programOne = ?, Users.programTwo = ? WHERE Users.id = ?")
      if err != nil {
        panic(err.Error())
      }
      defer stmtInsUser.Close()

      stmtInsCourses, err := db.Prepare("INSERT INTO CoursesTaken VALUES(?,?)")
      if err != nil {
        panic(err.Error())
      }
      defer stmtInsCourses.Close()

      remCourses, err := db.Prepare("DELETE FROM CoursesTaken WHERE CoursesTaken.uid = ?")
      if err != nil {
        panic(err.Error())
      }
      defer remCourses.Close()

      remCourses.Exec(id)
      pid1, pid2 := 0, 0
      err = db.QueryRow("SELECT P.id FROM Program P WHERE P.dept = ? and P.type = ?", user_info.CurrDept[0].Name , user_info.CurrDept[0].Type).Scan(&pid1)
      if err == sql.ErrNoRows {
        log.Println("empty sql result")
      }else if err != nil {
        log.Fatal(err)
      }
      if len(user_info.CurrDept) == 1{
        user_info.CurrDept = append(user_info.CurrDept, Program{"",""})
      }else {
        err = db.QueryRow("SELECT P.id FROM Program P WHERE P.dept = ? and P.type = ?", user_info.CurrDept[1].Name , user_info.CurrDept[1].Type).Scan(&pid2)
        if err == sql.ErrNoRows {
          log.Println("empty sql result")
        }else if err != nil {
          log.Fatal(err)
        }
      }


      stmtInsUser.Exec(user_info.SemLeft, user_info.GenEdsLeft, pid1, pid2, id)
      for _, dept := range user_info.DTaken {
          for _, course := range dept.CoursesTaken {
            var cid int
            err = db.QueryRow("SELECT C.id FROM Courses C WHERE C.cNumber = ? and C.dept = ?", course.Number, dept.Name).Scan(&cid)
            if err == sql.ErrNoRows {
              log.Println("empty sql result")
            }else if err != nil {
              log.Fatal(err)
            }
            stmtInsCourses.Exec(id, cid)
          }
      }

}

func GetResult(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Access-Control-Allow-Origin", "*")
  sessionId, ok := r.URL.Query()["sessionId"]
  if !ok || len(sessionId) < 1 {
       fmt.Fprintln(w, "Url Param 'sessionId' is missing")
       return
   }

   type Result struct {
     SessionId                  string              `json:"sessionId"`
     StrictRemainingCourses     Courses             `json:"strictRemainingCourses"`
     LooseRemainingCourses      []LooseReqCourse    `json:"looseRemainingCourses"`
     PossibleProg               []PossibleProgram   `json:"possiblePrograms"`
     OrderOfPrereqs             []PreReq            `json:"orderOfPrereqs"`
  }
  //connect to db
  db, err := sql.Open("mysql", "mason:pineappleB2@tcp(comp426finalproject.cqu5t9sfyvwq.us-east-2.rds.amazonaws.com:3306)/planner" )
  if err != nil {
   log.Fatal(err)
  }

  //get looseRemainingCourses remaining
  str:= "SELECT distinct PR.req , PR.numCourses, (SELECT COUNT(C.id) FROM Courses C, ProgramRequirements PR2, CoursesInProgram CP, Users U, UserSessions US, " +
  "CoursesTaken CT WHERE C.id = CP.cid and CP.prid = PR2.id and PR2.id = PR.id and U.id = US.uid and US.sessionId =   ? " +
  "and (C.id = CT.cid and CT.uid = U.id )) as count FROM  Program P, ProgramRequirements PR, Users U, UserSessions US WHERE P.id = PR.pid and U.id = US.uid and " +
  "US.sessionId = ?  and  (U.programOne = P.id or U.programTwo = P.id ) and PR.req != 'required' GROUP BY PR.req"
  rows1, err := db.Query(str, sessionId[0], sessionId[0])
  if err == sql.ErrNoRows {
    log.Println("empty sql result")
  }else if err != nil {
    log.Fatal(err)
  }
  defer rows1.Close()

  //go through sql result and add them to list
  reqMap := make(map[string]int)
   for rows1.Next() {
     req := ""
     numReq := 0
     numTaken := 0
     err := rows1.Scan(&req, &numReq, &numTaken)
     // log.Printf("%s\t%d\t%d", req, numReq, numTaken,)
     if err != nil {
       log.Fatal(err)
     }
     reqMap[req] = numReq - numTaken
   }

   if err = rows1.Err(); err != nil {
     log.Fatal(err)
   }
   looseCourses := make([]LooseReqCourse, 0)
   looseNum := 0
   for key, value := range reqMap {
     // log.Printf("%s\t%d", key, value)
     looseNum = looseNum + value
     if value > 0 {
       //should return all courses that fit requirement 'key' that user hasnt already taken
       str = "SELECT C.id, C.creditHours, C.cNumber, C.dept FROM Courses C, Program P, ProgramRequirements PR, CoursesInProgram CP, UserSessions US, Users U " +
             "WHERE  C.id = CP.cid and PR.id = CP.prid and PR.pid = P.id and (P.id = U.programOne or P.id = U.programTwo) and U.id = US.uid and US.sessionId = ? and PR.req = ? and C.id NOT IN " +
             "(SELECT C.id from Courses C, CoursesTaken CT, Users U, UserSessions US WHERE C.id = CT.cid and CT.uid = U.id and U.id = US.uid and US.sessionId = ?)"
       rows, err := db.Query(str, sessionId[0], key, sessionId[0])
       if err == sql.ErrNoRows {
         log.Println("empty sql result")
       }else if err != nil {
         log.Fatal(err)
       }
       defer rows.Close()
       for rows.Next() {
         looseCrs := new(LooseReqCourse)
         crs := new(Course)
         err := rows.Scan(&crs.Id, &crs.Hours, &crs.Number, &crs.Program)
         if err != nil {
           log.Fatal(err)
         }
        looseCrs.ReqCourse = *crs
        looseCrs.Requirement = key
        looseCrs.Number = value
        looseCourses = append(looseCourses, *looseCrs)
       }
     }
   }

  //get strictRemaining Courses remaining
  str = " SELECT C.id, C.creditHours, C.cNumber, C.dept FROM Courses C, Program P, ProgramRequirements PR, Users U, UserSessions US , CoursesInProgram CP WHERE PR.pid = P.id and C.id = CP.cid and CP.prid = PR.id and US.sessionId = ? and " +
  "US.uid = U.id and (U.programOne = P.id or U.programTwo = P.id) and PR.req = 'required' and C.id NOT IN (Select C.id from Courses C, Users U, UserSessions US, " +
  "CoursesTaken CT WHERE US.sessionId = ? and US.uid = U.id and CT.uid = U.id and CT.cid = C.id)"
  rows2, err := db.Query(str, sessionId[0], sessionId[0])
  if err == sql.ErrNoRows {
    log.Println("empty sql result")
  }else if err != nil {
    log.Fatal(err)
  }
  defer rows2.Close()

  //go through sql result and add them to list
  strictCourses := make([]Course, 0)
   for rows2.Next() {
     crs := new(Course)
     err := rows2.Scan(&crs.Id, &crs.Hours, &crs.Number, &crs.Program)
     if err != nil {
       log.Fatal(err)
     }
     strictCourses = append(strictCourses, *crs)
   }
   if err = rows2.Err(); err != nil {
     log.Fatal(err)
   }

   semLeft, genEds:= 0,0
   err = db.QueryRow("SELECT U.genEdsLeft, U.semLeft FROM Users U, UserSessions US WHERE U.id = US.uid and US.sessionId = ?", sessionId[0]).Scan(&genEds, &semLeft)
   if err == sql.ErrNoRows {
     log.Println("empty sql result")
   }else if err != nil {
     log.Fatal(err)
   }
  //used to calculate the number of hours of looseRemaining Courses that the person has left
   count := looseNum + len(strictCourses) + genEds
   classesRemaining := semLeft*6 - count

   rows3, err := db.Query("SELECT P.dept, P.type FROM Program P, Users U, UserSessions US WHERE P.numClasses <= ?  and U.id = US.uid and US.sessionId = ? and P.dept NOT IN (SELECT P.dept FROM Program P Where P.id = U.programOne or P.id = U.programTwo)", classesRemaining, sessionId[0])
   if err == sql.ErrNoRows {
     log.Println("empty sql result")
   }else if err != nil {
     log.Fatal(err)
   }
   defer rows3.Close()


   possPrograms := make([]PossibleProgram, 0)

    for rows3.Next() {
      possProgram := new(PossibleProgram)
      err := rows3.Scan(&possProgram.Dept, &possProgram.Type)


      if err != nil {
        log.Fatal(err)
      }
      // log.Printf("Poss Program %s \t %s" , possProgram.Dept, possProgram.Type)
      //looseRemainingCourses for PossProgram
      str= "SELECT PR.req , PR.numCourses, (SELECT COUNT(C.id) FROM Courses C, ProgramRequirements PR2, CoursesInProgram CP, Users U, UserSessions US, " +
      "CoursesTaken CT WHERE C.id = CP.cid and CP.prid = PR2.id and PR2.id = PR.id and U.id = US.uid and US.sessionId =   ? " +
      "and (C.id = CT.cid and CT.uid = U.id ) ) as count FROM  Program P, ProgramRequirements PR, Users U, UserSessions US WHERE P.id = PR.pid and U.id = US.uid and " +
      "US.sessionId = ?  and P.dept = ? and P.type = ? and PR.req != 'required' GROUP BY PR.req"

      rows4, err := db.Query(str, sessionId[0], sessionId[0],  possProgram.Dept, possProgram.Type, )
      if err == sql.ErrNoRows {
        log.Println("empty sql result")
      }else if err != nil {
        log.Fatal(err)
      }
      defer rows4.Close()
      reqMap2 := make(map[string]int)
       for rows4.Next() {
         req := ""
         numReq := 0
         numTaken := 0
         err := rows4.Scan(&req, &numReq, &numTaken)
         if err != nil {
           log.Fatal(err)
         }
         reqMap2[req] = numReq - numTaken
       }
       if err = rows4.Err(); err != nil {
         log.Fatal(err)
       }
       possLooseCourses := make([]LooseReqCourse, 0)
       possLooseNum := 0

       for key, value := range reqMap2 {
         possLooseNum = possLooseNum + value
         if value > 0 {
           //should return all courses that fit requirement 'key' that user hasnt already taken
           str = "SELECT distinct C.id, C.creditHours, C.cNumber, C.dept FROM Courses C, Program P, ProgramRequirements PR, CoursesInProgram CP, UserSessions US, Users U " +
                 "WHERE P.dept = ? and P.type = ? and C.id = CP.cid and PR.id = CP.prid and PR.pid = P.id and PR.req = ? and C.id NOT IN " +
                 "(SELECT C.id from Courses C, CoursesTaken CT, Users U, UserSessions US WHERE C.id = CT.cid and CT.uid = U.id and U.id = US.uid and US.sessionId = ?)"
           rows, err := db.Query(str, possProgram.Dept, possProgram.Type, key, sessionId[0])
           if err == sql.ErrNoRows {
             log.Println("empty sql result")
           }else if err != nil {
             log.Fatal(err)
           }
           defer rows.Close()
           for rows.Next() {
             looseCrs := new(LooseReqCourse)
             crs := new(Course)
             err := rows.Scan(&crs.Id, &crs.Hours, &crs.Number, &crs.Program)
             if err != nil {
               log.Fatal(err)
             }
            looseCrs.ReqCourse = *crs
            looseCrs.Requirement = key
            looseCrs.Number = value
             // log.Printf("Possible Loose Course %s \t %d", looseCrs.ReqCourse.Program, looseCrs.ReqCourse.Number)
            possLooseCourses = append(possLooseCourses, *looseCrs)
           }
         }
       }
       possProgram.LooseRemainingCourses = possLooseCourses

       //get strictRemaining Courses remaining
       str = " SELECT C.id, C.creditHours, C.cNumber, C.dept FROM Courses C, Program P, ProgramRequirements PR, CoursesInProgram CP WHERE PR.pid = P.id and " +
       "C.id = CP.cid and CP.prid = PR.id and P.dept = ? and P.type = ? and PR.req = 'required' AND C.id NOT IN ( Select CT.cid FROM Users U, UserSessions US, " +
         "CoursesTaken CT WHERE US.sessionId = ? and US.uid = U.id and CT.uid = U.id) "
       rows5, err := db.Query(str, possProgram.Dept, possProgram.Type, sessionId[0])
       if err == sql.ErrNoRows {
         log.Println("empty sql result")
       }else if err != nil {
         log.Fatal(err)
       }
       defer rows5.Close()

       //go through sql result and add them to list
       possStrictCourses := make([]Course, 0)
        for rows5.Next() {
          crs := new(Course)
          err := rows5.Scan(&crs.Id, &crs.Hours, &crs.Number, &crs.Program)
          if err != nil {
            log.Fatal(err)
          }
          // log.Printf("Possible Strict Course %s \t %d", crs.Program, crs.Number)
          possStrictCourses = append(possStrictCourses, *crs)
        }
        if err = rows5.Err(); err != nil {
          log.Fatal(err)
        }
        possProgram.StrictRemainingCourses = possStrictCourses
        possNumCourses := possLooseNum + len(possProgram.StrictRemainingCourses)
        possProgram.AvgHoursPerSem = float32((possNumCourses + classesRemaining)*3.0) / float32(semLeft)

        //prereqs
        prereqs := getStrictPrereqs(possProgram.StrictRemainingCourses, db)
        prereqs = append(prereqs, getLoosePrereqs(possProgram.LooseRemainingCourses, db)...)
        possProgram.OrderOfPrereqs = prereqs
      if possProgram.AvgHoursPerSem <= 18 {
        possPrograms = append(possPrograms, *possProgram)
      }
    }
    if err = rows3.Err(); err != nil {
      log.Fatal(err)
    }

    prereqs := getStrictPrereqs(strictCourses, db)
    prereqs = append(prereqs, getLoosePrereqs(looseCourses, db)...)


    res := Result{
      SessionId : sessionId[0] ,
      StrictRemainingCourses: strictCourses,
      LooseRemainingCourses: looseCourses,
      PossibleProg: possPrograms,
      OrderOfPrereqs: prereqs,
    }


  if err := json.NewEncoder(w).Encode(res); err != nil {
        panic(err)
    }
  }


func GetUserInfo(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Access-Control-Allow-Origin", "*")
  sessionId, ok := r.URL.Query()["sessionId"]
  if !ok || len(sessionId) < 1 {
       fmt.Fprintln(w, "Url Param 'sessionId' is missing")
       return
   }
   //connect to db
   db, err := sql.Open("mysql", "mason:pineappleB2@tcp(comp426finalproject.cqu5t9sfyvwq.us-east-2.rds.amazonaws.com:3306)/planner" )
  if err != nil {
    log.Fatal(err)
  }
  type User struct{
    Id         int     `json:"id"`
    Username   string  `json:"username"`
    SemLeft    int     `json:"semLeft"`
    GenEdsLeft int     `json:"genEdsLeft"`
    ProgramOne string  `json:"programOne"`
    ProgramTwo string  `json:"programTwo"`
  }
  usr := new(User)
  str := "SELECT U.id, U.user, U.semLeft, U.genEdsLeft, U.programOne, U.programTwo from Users U, UserSessions US WHERE US.uid = U.id and US.sessionId = ?"
  err = db.QueryRow(str , sessionId[0]).Scan(&usr.Id, &usr.Username, &usr.SemLeft, &usr.GenEdsLeft, &usr.ProgramOne, &usr.ProgramTwo)
  if err == sql.ErrNoRows {
    log.Println("empty sql result")
  }else if err != nil {
    log.Fatal(err)
  }

  if err := json.NewEncoder(w).Encode(usr); err != nil {
        panic(err)
    }
}
type LooseReqCourse struct {
  ReqCourse     Course     `json:"course"`
  Requirement   string     `json:"requirement"`
  Number        int        `json:"number"`
}

type PossibleProgram struct {
  Dept                    string           `json:"dept"`
  Type                    string           `json:"type"`
  AvgHoursPerSem          float32          `json:"avgHoursPerSem"`
  StrictRemainingCourses  Courses          `json:"strictRemainingCourses"`
  LooseRemainingCourses   []LooseReqCourse `json:"looseRemainingCourses"`
  OrderOfPrereqs          []PreReq         `json:"orderOfPrereqs"`
}
//helper methods below for generating user Session string
var src = rand.NewSource(time.Now().UnixNano())
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
func RandStringGenerator(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
