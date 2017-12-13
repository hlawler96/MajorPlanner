# MajorPlanner
The backend program is more or less together at this point, the database just has to be populated with data for the API to actually return non empty results. I haven't updated the aws server yet to reflect these changes so the old examples are still available there however several of the methods have changed (mainly due to implementing logging in). Below I included the schema for the database and information about all the different API calls and what the expected input/output is.

AWS URL: http://ec2-18-217-72-185.us-east-2.compute.amazonaws.com:8080

Ex/ http://ec2-18-217-72-185.us-east-2.compute.amazonaws.com:8080/Courses/?dept=COMP

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
THIS WAS CHANGED A BIT FROM BEFORE WITH PREREQS
Where SESSIONID is replaced by the sessionId you recieved at login. This returns a big JSON with the following information
- SessionId
- A list of "StrictRemainingCourses" that all must be taken
- A list of "LooseRemainingCourses" that fufill certain program requirements but dont all explicitly have ot be taken
- A list of Possible Programs that the user could take in their remaining time
- A list of Prereq Objects that have a list of courses and a description with them. With the list of courses the first course is the one that requires some or all of the remaining courses after the first course. Whether or not the prereqs are strictly requried or loosely required can be figured out if the type = 'required'

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
  type PreReq struct {
  Crs   []Course    `json:"Courses"`
  Des   string      `json:"Type"`
  }
```

An example result of this query is :
```
{"sessionId":"FZxduUieKLhLMAFASLxzUXvvPZBPOp","strictRemainingCourses":[{"id":4,"hours":3,"number":410,"program":"COMP"},{"id":5,"hours":4,"number":411,"program":"COMP"},{"id":19,"hours":4,"number":231,"program":"MATH"},{"id":20,"hours":4,"number":232,"program":"MATH"},{"id":6,"hours":3,"number":455,"program":"COMP"},{"id":7,"hours":3,"number":550,"program":"COMP"},{"id":65,"hours":3,"number":435,"program":"STOR"}],"looseRemainingCourses":[{"course":{"id":23,"hours":4,"number":116,"program":"PHYS "},"requirement":"1 of 2 Mechanics Courses","number":1},{"course":{"id":24,"hours":4,"number":118,"program":"PHYS "},"requirement":"1 of 2 Mechanics Courses","number":1},{"course":{"id":11,"hours":3,"number":426,"program":"COMP"},"requirement":"5 courses \u003e= 426","number":5},{"course":{"id":12,"hours":4,"number":541,"program":"COMP"},"requirement":"5 courses \u003e= 426","number":5},{"course":{"id":13,"hours":3,"number":521,"program":"COMP"},"requirement":"5 courses \u003e= 426","number":5},{"course":{"id":14,"hours":3,"number":431,"program":"COMP"},"requirement":"5 courses \u003e= 426","number":5},{"course":{"id":15,"hours":3,"number":530,"program":"COMP"},"requirement":"5 courses \u003e= 426","number":5},{"course":{"id":16,"hours":3,"number":435,"program":"COMP"},"requirement":"5 courses \u003e= 426","number":5},{"course":{"id":17,"hours":3,"number":433,"program":"COMP"},"requirement":"5 courses \u003e= 426","number":5},{"course":{"id":25,"hours":4,"number":101,"program":"ASTR"},"requirement":"1 of 11 additional Science with Lab","number":1},{"course":{"id":26,"hours":4,"number":101,"program":"BIOL"},"requirement":"1 of 11 additional Science with Lab","number":1},{"course":{"id":27,"hours":4,"number":202,"program":"BIOL"},"requirement":"1 of 11 additional Science with Lab","number":1},{"course":{"id":28,"hours":4,"number":205,"program":"BIOL"},"requirement":"1 of 11 additional Science with Lab","number":1},{"course":{"id":29,"hours":4,"number":101,"program":"CHEM"},"requirement":"1 of 11 additional Science with Lab","number":1},{"course":{"id":30,"hours":4,"number":102,"program":"CHEM"},"requirement":"1 of 11 additional Science with Lab","number":1},{"course":{"id":31,"hours":4,"number":101,"program":"GEOL"},"requirement":"1 of 11 additional Science with Lab","number":1},{"course":{"id":32,"hours":4,"number":117,"program":"PHYS "},"requirement":"1 of 11 additional Science with Lab","number":1},{"course":{"id":33,"hours":4,"number":119,"program":"PHYS "},"requirement":"1 of 11 additional Science with Lab","number":1},{"course":{"id":34,"hours":4,"number":351,"program":"PHYS "},"requirement":"1 of 11 additional Science with Lab","number":1},{"course":{"id":35,"hours":4,"number":352,"program":"PHYS "},"requirement":"1 of 11 additional Science with Lab","number":1},{"course":{"id":10,"hours":3,"number":283,"program":"COMP"},"requirement":"1 of 2 Discrete Courses","number":1},{"course":{"id":22,"hours":3,"number":381,"program":"MATH"},"requirement":"1 of 2 Discrete Courses","number":1},{"course":{"id":8,"hours":3,"number":547,"program":"MATH"},"requirement":"1 of 2 Linear Courses","number":1},{"course":{"id":18,"hours":3,"number":577,"program":"MATH"},"requirement":"1 of 2 Linear Courses","number":1}],"possiblePrograms":[{"dept":"MATH  ","type":"BA","avgHoursPerSem":13,"strictRemainingCourses":[{"id":22,"hours":3,"number":381,"program":"MATH"},{"id":40,"hours":3,"number":383,"program":"MATH"},{"id":49,"hours":3,"number":521,"program":"MATH"}],"looseRemainingCourses":[{"course":{"id":20,"hours":4,"number":232,"program":"MATH"},"requirement":"Math 232 or Math 283","number":1},{"course":{"id":48,"hours":4,"number":283,"program":"MATH"},"requirement":"Math 232 or Math 283","number":1},{"course":{"id":8,"hours":3,"number":547,"program":"MATH"},"requirement":"1 of 2 Linear Courses","number":1},{"course":{"id":18,"hours":3,"number":577,"program":"MATH"},"requirement":"1 of 2 Linear Courses","number":1},{"course":{"id":44,"hours":3,"number":524,"program":"MATH"},"requirement":"3 Courses \u003e= 500","number":3},{"course":{"id":45,"hours":4,"number":528,"program":"MATH"},"requirement":"3 Courses \u003e= 500","number":3},{"course":{"id":46,"hours":4,"number":529,"program":"MATH"},"requirement":"3 Courses \u003e= 500","number":3},{"course":{"id":50,"hours":3,"number":533,"program":"MATH"},"requirement":"3 Courses \u003e= 500","number":3},{"course":{"id":19,"hours":4,"number":231,"program":"MATH"},"requirement":"Math 231 or Math 241","number":1},{"course":{"id":47,"hours":4,"number":241,"program":"MATH"},"requirement":"Math 231 or Math 241","number":1}],"orderOfPrereqs":[{"Courses":[{"id":22,"hours":3,"number":381,"program":"MATH"},{"id":20,"hours":4,"number":232,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":40,"hours":3,"number":383,"program":"MATH"},{"id":21,"hours":4,"number":233,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":49,"hours":3,"number":521,"program":"MATH"},{"id":21,"hours":4,"number":233,"program":"MATH"},{"id":22,"hours":3,"number":381,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":20,"hours":4,"number":232,"program":"MATH"},{"id":19,"hours":4,"number":231,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":48,"hours":4,"number":283,"program":"MATH"},{"id":47,"hours":4,"number":241,"program":"MATH"},{"id":19,"hours":4,"number":231,"program":"MATH"}],"Type":"1 of Calc 1"},{"Courses":[{"id":8,"hours":3,"number":547,"program":"MATH"},{"id":21,"hours":4,"number":233,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":18,"hours":3,"number":577,"program":"MATH"},{"id":22,"hours":3,"number":381,"program":"MATH"},{"id":40,"hours":3,"number":383,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":44,"hours":3,"number":524,"program":"MATH"},{"id":40,"hours":3,"number":383,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":45,"hours":4,"number":528,"program":"MATH"},{"id":40,"hours":3,"number":383,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":46,"hours":4,"number":529,"program":"MATH"},{"id":40,"hours":3,"number":383,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":50,"hours":3,"number":533,"program":"MATH"},{"id":22,"hours":3,"number":381,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":19,"hours":4,"number":231,"program":"MATH"},{"id":43,"hours":3,"number":130,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":47,"hours":4,"number":241,"program":"MATH"},{"id":43,"hours":3,"number":130,"program":"MATH"}],"Type":"required"}]},{"dept":"MATH ","type":"Minor","avgHoursPerSem":12,"strictRemainingCourses":[{"id":22,"hours":3,"number":381,"program":"MATH"},{"id":40,"hours":3,"number":383,"program":"MATH"}],"looseRemainingCourses":[{"course":{"id":44,"hours":3,"number":524,"program":"MATH"},"requirement":"3 Courses \u003e= 500","number":3},{"course":{"id":45,"hours":4,"number":528,"program":"MATH"},"requirement":"3 Courses \u003e= 500","number":3},{"course":{"id":46,"hours":4,"number":529,"program":"MATH"},"requirement":"3 Courses \u003e= 500","number":3},{"course":{"id":50,"hours":3,"number":533,"program":"MATH"},"requirement":"3 Courses \u003e= 500","number":3},{"course":{"id":19,"hours":4,"number":231,"program":"MATH"},"requirement":"Math 231 or Math 241","number":1},{"course":{"id":47,"hours":4,"number":241,"program":"MATH"},"requirement":"Math 231 or Math 241","number":1},{"course":{"id":41,"hours":3,"number":155,"program":"STOR"},"requirement":"Math 231 or Math 241","number":1},{"course":{"id":20,"hours":4,"number":232,"program":"MATH"},"requirement":"Math 232 or Math 283","number":1},{"course":{"id":48,"hours":4,"number":283,"program":"MATH"},"requirement":"Math 232 or Math 283","number":1}],"orderOfPrereqs":[{"Courses":[{"id":22,"hours":3,"number":381,"program":"MATH"},{"id":20,"hours":4,"number":232,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":40,"hours":3,"number":383,"program":"MATH"},{"id":21,"hours":4,"number":233,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":44,"hours":3,"number":524,"program":"MATH"},{"id":40,"hours":3,"number":383,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":45,"hours":4,"number":528,"program":"MATH"},{"id":40,"hours":3,"number":383,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":46,"hours":4,"number":529,"program":"MATH"},{"id":40,"hours":3,"number":383,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":50,"hours":3,"number":533,"program":"MATH"},{"id":22,"hours":3,"number":381,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":19,"hours":4,"number":231,"program":"MATH"},{"id":43,"hours":3,"number":130,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":47,"hours":4,"number":241,"program":"MATH"},{"id":43,"hours":3,"number":130,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":41,"hours":3,"number":155,"program":"STOR"},{"id":36,"hours":3,"number":110,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":20,"hours":4,"number":232,"program":"MATH"},{"id":19,"hours":4,"number":231,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":48,"hours":4,"number":283,"program":"MATH"},{"id":47,"hours":4,"number":241,"program":"MATH"},{"id":19,"hours":4,"number":231,"program":"MATH"}],"Type":"1 of Calc 1"}]},{"dept":"ECON","type":"BA","avgHoursPerSem":13.5,"strictRemainingCourses":[{"id":51,"hours":3,"number":101,"program":"ECON"},{"id":52,"hours":3,"number":400,"program":"ECON"},{"id":53,"hours":3,"number":410,"program":"ECON"},{"id":55,"hours":3,"number":112,"program":"STOR"},{"id":41,"hours":3,"number":155,"program":"STOR"}],"looseRemainingCourses":[{"course":{"id":46,"hours":4,"number":529,"program":"MATH"},"requirement":"4 courses \u003e= 400","number":4},{"course":{"id":57,"hours":3,"number":445,"program":"ECON"},"requirement":"4 courses \u003e= 400","number":4},{"course":{"id":58,"hours":3,"number":511,"program":"ECON"},"requirement":"4 courses \u003e= 400","number":4},{"course":{"id":59,"hours":3,"number":545,"program":"ECON"},"requirement":"4 courses \u003e= 400","number":4},{"course":{"id":60,"hours":3,"number":560,"program":"ECON"},"requirement":"4 courses \u003e= 400","number":4},{"course":{"id":61,"hours":3,"number":460,"program":"ECON"},"requirement":"4 courses \u003e= 400","number":4},{"course":{"id":19,"hours":4,"number":231,"program":"MATH"},"requirement":"1 calculus from approved list","number":1},{"course":{"id":55,"hours":3,"number":112,"program":"STOR"},"requirement":"1 calculus from approved list","number":1}],"orderOfPrereqs":[{"Courses":[{"id":52,"hours":3,"number":400,"program":"ECON"},{"id":51,"hours":3,"number":101,"program":"ECON"}],"Type":"required"},{"Courses":[{"id":53,"hours":3,"number":410,"program":"ECON"},{"id":19,"hours":4,"number":231,"program":"MATH"},{"id":55,"hours":3,"number":112,"program":"STOR"}],"Type":"1 of Econ maths"},{"Courses":[{"id":53,"hours":3,"number":410,"program":"ECON"},{"id":51,"hours":3,"number":101,"program":"ECON"}],"Type":"required"},{"Courses":[{"id":55,"hours":3,"number":112,"program":"STOR"},{"id":36,"hours":3,"number":110,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":41,"hours":3,"number":155,"program":"STOR"},{"id":36,"hours":3,"number":110,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":46,"hours":4,"number":529,"program":"MATH"},{"id":40,"hours":3,"number":383,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":57,"hours":3,"number":445,"program":"ECON"},{"id":53,"hours":3,"number":410,"program":"ECON"}],"Type":"required"},{"Courses":[{"id":58,"hours":3,"number":511,"program":"ECON"},{"id":53,"hours":3,"number":410,"program":"ECON"},{"id":21,"hours":4,"number":233,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":59,"hours":3,"number":545,"program":"ECON"},{"id":57,"hours":3,"number":445,"program":"ECON"}],"Type":"required"},{"Courses":[{"id":60,"hours":3,"number":560,"program":"ECON"},{"id":61,"hours":3,"number":460,"program":"ECON"}],"Type":"required"},{"Courses":[{"id":61,"hours":3,"number":460,"program":"ECON"},{"id":53,"hours":3,"number":410,"program":"ECON"}],"Type":"required"},{"Courses":[{"id":19,"hours":4,"number":231,"program":"MATH"},{"id":43,"hours":3,"number":130,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":55,"hours":3,"number":112,"program":"STOR"},{"id":36,"hours":3,"number":110,"program":"MATH"}],"Type":"required"}]},{"dept":"STOR","type":"Minor","avgHoursPerSem":10.5,"strictRemainingCourses":[],"looseRemainingCourses":[{"course":{"id":63,"hours":3,"number":305,"program":"STOR"},"requirement":"3 From approved list","number":3},{"course":{"id":64,"hours":3,"number":415,"program":"STOR"},"requirement":"3 From approved list","number":3},{"course":{"id":65,"hours":3,"number":435,"program":"STOR"},"requirement":"3 From approved list","number":3},{"course":{"id":66,"hours":3,"number":445,"program":"STOR"},"requirement":"3 From approved list","number":3},{"course":{"id":67,"hours":3,"number":455,"program":"STOR"},"requirement":"3 From approved list","number":3},{"course":{"id":62,"hours":3,"number":215,"program":"STOR"},"requirement":"Stor 215 or Math 381","number":1},{"course":{"id":22,"hours":3,"number":381,"program":"MATH"},"requirement":"Stor 215 or Math 381","number":1}],"orderOfPrereqs":[{"Courses":[{"id":63,"hours":3,"number":305,"program":"STOR"},{"id":41,"hours":3,"number":155,"program":"STOR"}],"Type":"required"},{"Courses":[{"id":64,"hours":3,"number":415,"program":"STOR"},{"id":8,"hours":3,"number":547,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":65,"hours":3,"number":435,"program":"STOR"},{"id":21,"hours":4,"number":233,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":66,"hours":3,"number":445,"program":"STOR"},{"id":65,"hours":3,"number":435,"program":"STOR"}],"Type":"required"},{"Courses":[{"id":67,"hours":3,"number":455,"program":"STOR"},{"id":41,"hours":3,"number":155,"program":"STOR"}],"Type":"required"},{"Courses":[{"id":62,"hours":3,"number":215,"program":"STOR"},{"id":36,"hours":3,"number":110,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":22,"hours":3,"number":381,"program":"MATH"},{"id":20,"hours":4,"number":232,"program":"MATH"}],"Type":"required"}]}],"orderOfPrereqs":[{"Courses":[{"id":4,"hours":3,"number":410,"program":"COMP"},{"id":3,"hours":4,"number":401,"program":"COMP"}],"Type":"required"},{"Courses":[{"id":19,"hours":4,"number":231,"program":"MATH"},{"id":43,"hours":3,"number":130,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":20,"hours":4,"number":232,"program":"MATH"},{"id":19,"hours":4,"number":231,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":6,"hours":3,"number":455,"program":"COMP"},{"id":10,"hours":3,"number":283,"program":"COMP"},{"id":22,"hours":3,"number":381,"program":"MATH"}],"Type":"1 discrete"},{"Courses":[{"id":6,"hours":3,"number":455,"program":"COMP"},{"id":3,"hours":4,"number":401,"program":"COMP"}],"Type":"required"},{"Courses":[{"id":7,"hours":3,"number":550,"program":"COMP"},{"id":10,"hours":3,"number":283,"program":"COMP"},{"id":22,"hours":3,"number":381,"program":"MATH"}],"Type":"1 discrete"},{"Courses":[{"id":7,"hours":3,"number":550,"program":"COMP"},{"id":4,"hours":3,"number":410,"program":"COMP"}],"Type":"required"},{"Courses":[{"id":65,"hours":3,"number":435,"program":"STOR"},{"id":21,"hours":4,"number":233,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":23,"hours":4,"number":116,"program":"PHYS "},{"id":19,"hours":4,"number":231,"program":"MATH"},{"id":20,"hours":4,"number":232,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":24,"hours":4,"number":118,"program":"PHYS "},{"id":20,"hours":4,"number":232,"program":"MATH"},{"id":19,"hours":4,"number":231,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":11,"hours":3,"number":426,"program":"COMP"},{"id":5,"hours":4,"number":411,"program":"COMP"},{"id":4,"hours":3,"number":410,"program":"COMP"}],"Type":"required"},{"Courses":[{"id":12,"hours":4,"number":541,"program":"COMP"},{"id":5,"hours":4,"number":411,"program":"COMP"},{"id":4,"hours":3,"number":410,"program":"COMP"}],"Type":"required"},{"Courses":[{"id":13,"hours":3,"number":521,"program":"COMP"},{"id":5,"hours":4,"number":411,"program":"COMP"},{"id":4,"hours":3,"number":410,"program":"COMP"}],"Type":"required"},{"Courses":[{"id":14,"hours":3,"number":431,"program":"COMP"},{"id":5,"hours":4,"number":411,"program":"COMP"},{"id":4,"hours":3,"number":410,"program":"COMP"}],"Type":"required"},{"Courses":[{"id":15,"hours":3,"number":530,"program":"COMP"},{"id":5,"hours":4,"number":411,"program":"COMP"},{"id":4,"hours":3,"number":410,"program":"COMP"}],"Type":"required"},{"Courses":[{"id":16,"hours":3,"number":435,"program":"COMP"},{"id":4,"hours":3,"number":410,"program":"COMP"}],"Type":"required"},{"Courses":[{"id":17,"hours":3,"number":433,"program":"COMP"},{"id":4,"hours":3,"number":410,"program":"COMP"}],"Type":"required"},{"Courses":[{"id":27,"hours":4,"number":202,"program":"BIOL"},{"id":26,"hours":4,"number":101,"program":"BIOL"}],"Type":"required"},{"Courses":[{"id":28,"hours":4,"number":205,"program":"BIOL"},{"id":27,"hours":4,"number":202,"program":"BIOL"}],"Type":"required"},{"Courses":[{"id":30,"hours":4,"number":102,"program":"CHEM"},{"id":29,"hours":4,"number":101,"program":"CHEM"}],"Type":"required"},{"Courses":[{"id":32,"hours":4,"number":117,"program":"PHYS "},{"id":20,"hours":4,"number":232,"program":"MATH"},{"id":23,"hours":4,"number":116,"program":"PHYS "}],"Type":"required"},{"Courses":[{"id":33,"hours":4,"number":119,"program":"PHYS "},{"id":20,"hours":4,"number":232,"program":"MATH"},{"id":24,"hours":4,"number":118,"program":"PHYS "}],"Type":"required"},{"Courses":[{"id":34,"hours":4,"number":351,"program":"PHYS "},{"id":33,"hours":4,"number":119,"program":"PHYS "},{"id":32,"hours":4,"number":117,"program":"PHYS "}],"Type":"1 of electro"},{"Courses":[{"id":34,"hours":4,"number":351,"program":"PHYS "},{"id":21,"hours":4,"number":233,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":35,"hours":4,"number":352,"program":"PHYS "},{"id":34,"hours":4,"number":351,"program":"PHYS "}],"Type":"required"},{"Courses":[{"id":10,"hours":3,"number":283,"program":"COMP"},{"id":19,"hours":4,"number":231,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":22,"hours":3,"number":381,"program":"MATH"},{"id":20,"hours":4,"number":232,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":8,"hours":3,"number":547,"program":"MATH"},{"id":21,"hours":4,"number":233,"program":"MATH"}],"Type":"required"},{"Courses":[{"id":18,"hours":3,"number":577,"program":"MATH"},{"id":22,"hours":3,"number":381,"program":"MATH"},{"id":40,"hours":3,"number":383,"program":"MATH"}],"Type":"required"}]}

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
 "genEdsLeft": 3
 ```
