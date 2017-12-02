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


 courses := make([]*CourseDept, 0)
  for rows.Next() {
    crs := new(CourseDept)
    err := rows.Scan(&crs.Id, &crs.Hours, &crs.Number, &crs.dept)
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
	stmtInsUser, err := db.Prepare("INSERT INTO Users (user, pass, name) VALUES (?,?,?)")
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
      _, err = stmtInsUser.Exec(user[0],pass[0],name[0])

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
   str := "SELECT C.id, C.creditHours, C.cNumber, C.pid FROM Courses C, CoursesTaken CT, UserSessions US WHERE C.id = CT.cid and CT.uid = US.uid and US.sessionId = ?"
   rows, err := db.Query(str, sessionId[0])
   if err != nil {
     log.Fatal(err)
   }
   defer rows.Close()

   //go through sql result and add them to list
   courses := make([]*Course, 0)
    for rows.Next() {
      crs := new(Course)
      err := rows.Scan(&crs.Id, &crs.Hours, &crs.Number, &crs.Pid)
      if err != nil {
        log.Fatal(err)
      }
      courses = append(courses, crs)
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
  Sample Post Message Body Request
  {"sessionId":1, "deptTaken":[{"name":"COMP", "coursesTaken": [{"dept":"COMP","number":"110"},{"dept":"COMP","number":"401"}]}, {"name":"MATH", "coursesTaken": [{"dept":"MATH","number":"233"}]}],"currDept":[{"COMP","BS"},{"MATH","Minor"}], "semLeft": 4, "genEdsLeft": 3}
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
  id, ok := r.URL.Query()["id"]
  if !ok || len(id) < 1 {
       fmt.Fprintln(w, "Url Param 'id' is missing")
       return
   }
   type LooseReqCourse struct {
     ReqCourse     Course     `json:"course"`
     Requirement   string     `json:"requirement"`
   }
   type LooseReqCourses []LooseReqCourse

   type Prereqs []Courses

   type PossibleProgram struct {
     Name                    string           `json:"name"`
     AvgHoursPerSem          float32          `json:"avgHoursPerSem"`
     StrictRemainingCourses  Courses          `json:"strictRemainingCourses"`
     LooseRemainingCourses   LooseReqCourses  `json:"looseRemainingCourses"`
     OrderOfPrereqs          Prereqs          `json:"orderOfPrereqs"`
   }

   type PossiblePrograms []PossibleProgram

   type Result struct {
     Id                         string              `json:"id"`
     StrictRemainingCourses     Courses             `json:"strictRemainingCourses"`
     LooseRemainingCourses      LooseReqCourses     `json:"looseRemainingCourses"`
     PossibleProg               PossiblePrograms    `json:"possiblePrograms"`
     OrderOfPrereqs             Prereqs             `json:"orderOfPrereqs"`
  }

  res := Result{
    Id : id[0] ,
    StrictRemainingCourses: Courses{Course{1, 3, 550, 1}, Course{2, 3, 455, 1}},
    LooseRemainingCourses: LooseReqCourses{{Course{3, 3,426, 1},"Greater than or equal to - 426"}, {Course{4, 3, 433, 1}, "Greater than or equal to - 426"}},
    PossibleProg: PossiblePrograms{PossibleProgram{"Mathematics BS", 14.333, Courses{Course{4,3, 547, 2}, Course{5,3,521, 2}}, LooseReqCourses{LooseReqCourse{Course{5,3, 528, 2}, "Greater than or equal to - 500"}},
    Prereqs{Courses{Course{6,3, 231, 2},Course{7,3,232,2}}, Courses{Course{7,3,232, 2}, Course{8,3,233,2}}}},
    OrderOfPrereqs: Prereqs{Courses{Course{6,3, 231, 2},Course{7,3,232,2}}, Courses{Course{7,3,232, 2}, Course{8,3,233,2}}},
  }}

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
