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
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)

  //open database
  db, err := sql.Open("mysql", "mason:pineappleB2@tcp(comp426finalproject.cqu5t9sfyvwq.us-east-2.rds.amazonaws.com:3306)/planner" )
 if err != nil {
   log.Fatal(err)
 }
 //check for dept param
 dept, ok := r.URL.Query()["dept"]
 sqlCond := " and P.dept = '"
 //if dept param is given then only return classes in that dept
 if !ok || len(dept) < 1 {
    sqlCond = ""
  }else{
    sqlCond = sqlCond + dept[0] + "'"
  }
  rows, err := db.Query("SELECT C.id, C.creditHours, C.cNumber, P.dept FROM Courses C, Program P where P.id = C.pid" + sqlCond)
  if err != nil {
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
fmt.Fprintln(w, "Hello world :)")
}
//fully functioning
func Login(w http.ResponseWriter, r *http.Request){
  // example call
  // http://localhost:8080/Login/?username=user&password=pass



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
      if err != nil {
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
   str := "SELECT C.id, C.creditHours, C.cNumber, P.dept FROM Courses C, CoursesTaken CT, UserSessions US, Program P WHERE C.id = CT.cid and CT.uid = US.uid and US.sessionId = ? and C.pid = P.id"
   rows, err := db.Query(str, sessionId[0])
   if err != nil {
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
//needs to be tested but I think this is working
func PostUserInformation(w http.ResponseWriter, r *http.Request){
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
  switch {
  	case err == sql.ErrNoRows:
  		return
  	case err != nil:
  		log.Fatal(err)
  	default:
      if len(user_info.CurrDept) == 1{
        user_info.CurrDept = append(user_info.CurrDept, Program{"",""})
      }
      // Prepare statements for inserting data
      stmtInsUser, err := db.Prepare("UPDATE Users SET semLeft = ?, genEdsLeft = ?, programOne = ?, programTwo = ? WHERE id = ?")
      if err != nil {
        panic(err.Error())
      }
      defer stmtInsUser.Close()
      stmtInsCourses, err := db.Prepare("INSERT INTO CoursesTaken VALUES(?,?)")
      if err != nil {
        panic(err.Error())
      }
      defer stmtInsCourses.Close()

      stmtInsUser.Exec(user_info.SemLeft, user_info.GenEdsLeft, user_info.CurrDept[0], user_info.CurrDept[1], user_info.SessionId)
      for _, dept := range user_info.DTaken {
          for _, course := range dept.CoursesTaken {
            var cid int
            err = db.QueryRow("SELECT C.id FROM Courses C, Program P WHERE C.pid = P.id and C.cNumber = ? and P.dept = ?", course.Number, dept.Name).Scan(&cid)
            if err != nil {
              panic(err.Error())
            }
            stmtInsCourses.Exec(id, cid)
          }
      }


  	}
}

func testGetResult(w http.ResponseWriter, r *http.Request){
  sessionId, ok := r.URL.Query()["sessionId"]
  if !ok || len(sessionId) < 1 {
       fmt.Fprintln(w, "Url Param 'sessionId' is missing")
       return
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
     OrderOfPrereqs          []Courses        `json:"orderOfPrereqs"`
   }

   type Result struct {
     SessionId                  string              `json:"sessionId"`
     StrictRemainingCourses     Courses             `json:"strictRemainingCourses"`
     LooseRemainingCourses      []LooseReqCourse    `json:"looseRemainingCourses"`
     PossibleProg               []PossibleProgram   `json:"possiblePrograms"`
     OrderOfPrereqs             []Courses           `json:"orderOfPrereqs"`
  }
  //connect to db
  db, err := sql.Open("mysql", "mason:pineappleB2@tcp(comp426finalproject.cqu5t9sfyvwq.us-east-2.rds.amazonaws.com:3306)/planner" )
  if err != nil {
   log.Fatal(err)
  }

  //get looseRemainingCourses remaining
  str:= "SELECT PR.req , PR.numCourses, COUNT(DISTINCT C.id) FROM Program P, ProgramRequirements PR, Users U, UserSessions US, Courses C, CoursesTaken CT, CoursesInProgram CP" +
  " WHERE US.sessionId = ? and U.id = US.uid and PR.pid = P.id and (U.programOne = P.id or U.programTwo = P.id) and PR.req != 'required' and CT.uid = U.id and C.id = CT.cid and CP.cid = C.id and  GROUP BY PR.req , PR.numCourses"
  rows1, err := db.Query(str, sessionId[0])
  if err != nil {
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
     looseNum = looseNum + value
     if value > 0 {
       //should return all courses that fit requirement 'key' that user hasnt already taken
       str = "SELECT C.id, C.creditHours, C.cNumber, P.dept FROM Courses C, Program P, ProgramRequirements PR, CoursesInProgram CP " +
             "WHERE C.pid = P.id and C.id = CP.cid and PR.id = CP.prid and PR.pid = P.id and PR.req = ? and C.id NOT IN " +
             "(SELECT C.id from Courses C, CoursesTaken CT, Users U, UserSessions US WHERE C.id = CT.cid and CT.uid = U.id and U.id = US.uid and US.sessionId = ?)"
       rows, err := db.Query(str, key, sessionId[0])
       if err != nil {
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
  str = " SELECT C.id, C.creditHours, C.cNumber, P.dept FROM Courses C, Program P, ProgramRequirements PR, Users U, UserSessions US , CoursesInProgram CP WHERE C.pid = P.id and PR.pid = P.id and C.id = CP.cid and CP.prid = PR.id and US.sessionId = ? and " +
  "US.uid = U.id and (U.programOne = P.id or U.programTwo = P.id) and PR.req = 'required' and C.id NOT IN (Select C.id from Courses C, Users U, UserSessions US, " +
  "CoursesTaken CT WHERE US.sessionId = ? and US.uid = U.id and CT.uid = U.id and CT.cid = C.id)"
  rows2, err := db.Query(str, sessionId[0], sessionId[0])
  if err != nil {
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
   var genEds int
   var semLeft int
   err = db.QueryRow("SELECT U.genEdsLeft, U.semLeft FROM Users U, UserSessions US WHERE U.uid = US.id and US.sessionId = ?", sessionId[0]).Scan(&genEds, &semLeft)
  //used to calculate the number of hours of looseRemaining Courses that the person has left
   count := looseNum + len(strictCourses) + genEds
   classesRemaining := semLeft - count
   rows3, err := db.Query("SELECT P.dept, P.type FROM Program P WHERE P.numClasses <= ?", classesRemaining)
   if err != nil {
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
      //looseRemainingCourses for PossProgram
      str = "SELECT PR.req , PR.numCourses, COUNT(DISTINCT C.id) FROM Program P, ProgramRequirements PR, Users U, UserSessions US, Courses C, CoursesTaken CT, CoursesInProgram CP" +
      " WHERE US.sessionId = ? and U.id = US.uid and PR.pid = P.id and P.dept = ? and P.type = ? and PR.req != 'required' and CT.uid = U.id and C.id = CT.cid and CP.cid = C.id and GROUP BY PR.req , PR.numCourses"
      rows4, err := db.Query(str, possProgram.Dept, possProgram.Type, sessionId[0])
      if err != nil {
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
           str = "SELECT C.id, C.creditHours, C.cNumber, P.dept FROM Courses C, Program P, ProgramRequirements PR, CoursesInProgram CP " +
                 "WHERE C.pid = P.id and C.id = CP.cid and PR.id = CP.prid and PR.pid = P.id and PR.req = ? and C.id NOT IN " +
                 "(SELECT C.id from Courses C, CoursesTaken CT, Users U, UserSessions US WHERE C.id = CT.cid and CT.uid = U.id and U.id = US.uid and US.sessionId = ?)"
           rows, err := db.Query(str, key, sessionId[0])
           if err != nil {
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
            possLooseCourses = append(possLooseCourses, *looseCrs)
           }
         }
       }
       possProgram.LooseRemainingCourses = possLooseCourses

       //get strictRemaining Courses remaining
       str = " SELECT C.id, C.creditHours, C.cNumber, P.dept FROM Courses C, Program P, ProgramRequirements PR, Users U, UserSessions US , CoursesInProgram CP WHERE C.pid = P.id and PR.pid = P.id and C.id = CP.cid and CP.prid = PR.id and US.sessionId = ? and " +
       "US.uid = U.id and P.dept = ? and P.type = ? and PR.req = 'required' and C.id NOT IN (Select C.id from Courses C, Users U, UserSessions US, " +
       "CoursesTaken CT WHERE US.sessionId = ? and US.uid = U.id and CT.uid = U.id and CT.cid = C.id)"
       rows5, err := db.Query(str,sessionId[0],possProgram.Dept, possProgram.Type, sessionId[0])
       if err != nil {
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
          possStrictCourses = append(strictCourses, *crs)
        }
        if err = rows5.Err(); err != nil {
          log.Fatal(err)
        }
        possProgram.StrictRemainingCourses = possStrictCourses

        possNumCourses := possLooseNum + len(possProgram.StrictRemainingCourses)
        possProgram.AvgHoursPerSem = float32(possNumCourses + classesRemaining) / float32(semLeft)

        str = " SELECT PR.prid , C.creditHours, C.cNumber, P.dept from Prereqs PR, Courses C, Program P WHERE PR.cid = ? and C.id = PR.prid and P.id = C.pid"
        prereqs := make([]Courses, 0)
        for _ , c := range possProgram.LooseRemainingCourses {
          rows, err := db.Query(str,c.ReqCourse.Id)
          switch {
          	case err == sql.ErrNoRows:
          		//do nothing
          	case err != nil:
          		log.Fatal(err)
          	default:
              for rows.Next() {
                prereq := make([]Course, 0)
                crs := new(Course)
                err := rows.Scan(&crs.Id, &crs.Hours, &crs.Number, &crs.Program)
                if err != nil {
                  log.Fatal(err)
                }
                prereq = append(prereq, c.ReqCourse)
                prereq = append(prereq, *crs)
                prereqs = append(prereqs, prereq)
              }
          	}
          defer rows.Close()

        }
        for _ , c := range possProgram.StrictRemainingCourses {
          rows, err := db.Query(str,c.Id)
          switch {
          	case err == sql.ErrNoRows:
          		//do nothing
          	case err != nil:
          		log.Fatal(err)
          	default:
              for rows.Next() {
                prereq := make([]Course, 0)
                crs := new(Course)
                err := rows.Scan(&crs.Id, &crs.Hours, &crs.Number, &crs.Program)
                if err != nil {
                  log.Fatal(err)
                }
                prereq = append(prereq, c)
                prereq = append(prereq, *crs)
                prereqs = append(prereqs, prereq)
              }
          	}
          defer rows.Close()

        }
        possProgram.OrderOfPrereqs = prereqs
      if possProgram.AvgHoursPerSem < 18 {
        possPrograms = append(possPrograms, *possProgram)
      }
    }
    if err = rows3.Err(); err != nil {
      log.Fatal(err)
    }
    str = " SELECT PR.prid , C.creditHours, C.cNumber, P.dept from Prereqs PR, Courses C, Program P WHERE PR.cid = ? and C.id = PR.prid and P.id = C.pid"
    prereqs := make([]Courses, 0)
    for _ , c := range looseCourses {
      rows, err := db.Query(str,c.ReqCourse.Id)
      switch {
      	case err == sql.ErrNoRows:
      		//do nothing
      	case err != nil:
      		log.Fatal(err)
      	default:
          for rows.Next() {
            prereq := make([]Course, 0)
            crs := new(Course)
            err := rows.Scan(&crs.Id, &crs.Hours, &crs.Number, &crs.Program)
            if err != nil {
              log.Fatal(err)
            }
            prereq = append(prereq, c.ReqCourse)
            prereq = append(prereq, *crs)
            prereqs = append(prereqs, prereq)
          }
      	}
      defer rows.Close()

    }
    for _ , c := range strictCourses {
      rows, err := db.Query(str,c.Id)
      switch {
      	case err == sql.ErrNoRows:
      		//do nothing
      	case err != nil:
      		log.Fatal(err)
      	default:
          for rows.Next() {
            prereq := make([]Course, 0)
            crs := new(Course)
            err := rows.Scan(&crs.Id, &crs.Hours, &crs.Number, &crs.Program)
            if err != nil {
              log.Fatal(err)
            }
            prereq = append(prereq, c)
            prereq = append(prereq, *crs)
            prereqs = append(prereqs, prereq)
          }
      	}
      defer rows.Close()

    }

    res := Result{
      SessionId : sessionId[0] ,
      StrictRemainingCourses: strictCourses,
      LooseRemainingCourses: looseCourses,
      PossibleProg: possPrograms,
      OrderOfPrereqs: prereqs,
    }


  // res := Result{
  //   SessionId : sessionId[0] ,
  //   StrictRemainingCourses: Courses{Course{1, 3, 550, "COMP"}, Course{2, 3, 455, "COMP"}},
  //   LooseRemainingCourses: LooseReqCourses{{Course{3, 3,426, "COMP"},"Greater than or equal to - 426",5}, {Course{4, 3, 433, "COMP"}, "Greater than or equal to - 426", 5}},
  //   PossibleProg: PossiblePrograms{PossibleProgram{"Math","BS", 14.333, Courses{Course{4,3, 547, "MATH"}, Course{5,3,521, "MATH"}}, LooseReqCourses{LooseReqCourse{Course{5,3, 528, "MATH"}, "Greater than or equal to - 500", 3}},
  //   Prereqs{Courses{Course{6,3, 231, "MATH"},Course{7,3,232,"MATH"}}, Courses{Course{7,3,232, "MATH"}, Course{8,3,233,"MATH"}}}}},
  //   OrderOfPrereqs : Prereqs{Courses{Course{6,3, 231, "MATH"},Course{7,3,232,"MATH"}}, Courses{Course{7,3,232, "MATH"}, Course{8,3,233,"MATH"}}},
  // }

  if err := json.NewEncoder(w).Encode(res); err != nil {
        panic(err)
    }


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
