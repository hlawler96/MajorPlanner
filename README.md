# MajorPlanner
## Database Schema

There is an Excel File with the Schema laid out but for simplicity I copied it here as well:

### Users		
- id	INT	Primary Key 
- user	VARCHAR(25)	
- pass	VARCHAR(25)	
- semLeft	INT	
- genEdsLeft	INT	
- programOne	INT	Foreign Key to Programs
- programTwo	INT	Foreign Key to Programs
		
		
### Courses		
- id	INT	Primary Key
- creditHours	INT	
- cNumber	INT	
- pid	INT	Foreign Key for Program
		
### CoursesTaken		
- uid	INT	Foreign Key for User
- cid	INT	Foreign Key for Courses
		
### Program		
- id	INT	Primary Key
- name	VARCHAR(25)	
- dept	VARCHAR(10)	
- type	ENUM('BS','BA','Minor')	
- numClasses	INT	
		
### ProgramRequirements		
- id	INT	Primary Key 
- req	VARCHAR(25)	
- numCourses	INT	
- pid		Foreign Key to Program 
		
### CoursesInProgram		
- cid	INT	Foreign Key to Courses
- prid	INT	Foreign Key to ProgramRequirements
		
### Prereqs		
- cid	INT	Foreign Key to Courses
- prid	INT	Foreign Key to Courses
		
### UserSessions		
- uid	INT	Primary Key
- sessionId	VARCHAR(30)	

## GET API Calls

### Test
Use `[hostname]:8080/`
This returns a test message. This can be used to check client side connection to server.

### All Courses
Use `[hostname]:8080/Courses`
This returns all of the Courses in the Database.
  
### Courses in Dept
Use `[hostname]:8080/Courses/?dept=DEPT`
Where DEPT is replaced with the dept Code such as COMP. This returns all of the Courses in that Dept in the Database.
  
### Login
Use `[hostname]:8080//Login/?username=USERNAME&password=PASSWORD`
Where USERNAME and PASSWORD are replaced with real usernames and passwords. This returns a 30 character Session Id that is needed for all user specific API calls. This method will return an empty array if username or password are not correct.
  
### SignUp
Use `[hostname]:8080/SignUp/?username=USERNAME&password=PASSWORD`
Where USERNAME and PASSWORD are replaced with new usernames and passwords. This also returns a 30 character Session Id that is needed for all user specific API calls. This will return an empty array if username is already taken.
  
### Courses Taken By User
Use `[hostname]:8080/CoursesTaken/?sessionId=SESSIONID`
Where SESSIONID is replaced by the sessionId you recieved at login. This returns an array of Course objects that the given user has already taken.

### User Info
Use `[hostname]:8080/UserInfo/sessionId=SESSIONID`
Where SESSIONID is replaced by the sessionId you recieved at login. This returns information about the user profile including 
- id
- username
- semesters left
- gen ed classes left
- Their primary major
- Their secondary major/minor (if exists, empty string if it doesn't exist)

### Application Result
Use `[hostname]:8080/PossiblePrograms/?sessionId=SESSIONID`
Where SESSIONID is replaced by the sessionId you recieved at login. This returns a big JSON with the following information
- SessionId
- A list of "StrictRemainingCourses" that all must be taken
- A list of "LooseRemainingCourses" that fufill certain program requirements but dont all explicitly have ot be taken
- A list of Possible Programs that the user could take in their remaining time
- A list of Course pairs that represent prereqs (the second is a prereq on the first)

The following code represents the object structure in go that is converted to the result to help you decipher the format.

```
   type Course struct {
     Id      int    `json:"id"`
     Hours   int     `json:"hours"`
     Number  int     `json:"number"`
     Program string  `json:"program"`
   }

   type Courses []Course
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
```

An example result of this query is :
```
{"sessionId":"jwGzoQQUmGmONbpqnDBPJeOrncVHbv","strictRemainingCourses":[{"id":1,"hours":3,"number":550,"dept":"COMP"},{"id":2,"hours":3,"number":455,"dept":"COMP"}],"looseRemainingCourses":[{"course":{"id":3,"hours":3,"number":426,"dept":"COMP"},"requirement":"Greater than or equal to - 426","number":6},{"course":{"id":4,"hours":3,"number":433,"dept":"COMP"},"requirement":"Greater than or equal to - 426","number":6}],"PossiblePrograms":[{"dept":"MATH","type":"BA","avgHoursPerSem":14.333,"strictRemainingCourses":[{"id":4,"hours":3,"number":547,"dept":"MATH"},{"id":5,"hours":3,"number":521,"dept":"MATH"}],"looseRemainingCourses":[{"course":{"id":5,"hours":3,"number":528,"dept":"MATH"},"requirement":"Greater than or equal to - 500"},"number":5],"orderOfPrereqs":[[{"id":6,"hours":3,"number":231,"dept":"MATH"},{"id":7,"hours":3,"number":232,"dept":"MATH"}],[{"id":7,"hours":3,"number":232,"dept":"MATH"},{"id":8,"hours":3,"number":233,"dept":"MATH"}]]}]}
```
## POST API Calls

### Provide Extra User Info
Use `[hostname]:8080/PossiblePrograms`
Here's an example of the message body with the user information to post
```{
 "sessionId":"jwGzoQQUmGmONbpqnDBPJeOrncVHbv",
 "deptTaken":[ {"name":"COMP", 
 		"coursesTaken": [{"dept":"COMP","number":110}, { "dept":"COMP","number":401}] } ,
	       {"name":"MATH", 
	        "coursesTaken": [{"dept":"MATH","number":233} ] } ],
 "currDept":[{"name":"COMP","type":"BS"},{"name":"MATH","type":"Minor"}],
 "semLeft": 4,
 "genEdsLeft": 3 ```
