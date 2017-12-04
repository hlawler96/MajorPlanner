# MajorPlanner
## Database Schema

There is an Excel File with the Schema laid out but for simplicity I copied it here as well:

Users		
id	INT	Primary Key
user	VARCJAR(25)	
pass	VARCJAR(25)	
semLeft	INT	
genEdsLeft	INT	
programOne	INT	Foreign Key to Programs
programTwo	INT	Foreign Key to Programs
		
		
Courses		
id	INT	Primary Key
creditHours	INT	
cNumber	INT	
pid	INT	Foreign Key for Program
		
CoursesTaken		
uid	INT	Foreign Key for User
cid	INT	Foreign Key for Courses
		
Program		
id	INT	Primary Key
name	VARCHAR(25)	
dept	VARCHAR(10)	
type	ENUM('BS','BA','Minor')	
numClasses	INT	
		
ProgramRequirements		
id	INT	Primary Key 
req	VARCHAR(25)	
numCourses	INT	
pid		Foreign Key to Program 
		
CoursesInProgram		
cid	INT	Foreign Key to Courses
prid	INT	Foreign Key to ProgramRequirements
		
Prereqs		
cid	INT	Foreign Key to Courses
prid	INT	Foreign Key to Courses
		
UserSessions		
uid	INT	Primary Key
sessionId	VARCHAR(30)	

## GET API Calls

### Test
Use '<hostname>:8080/'
This returns a test message. This can be used to check client side connection to server.

### All Courses
Use '<hostname>:8080/Courses'
This returns all of the Courses in the Database.
  
### Courses in Dept
Use '<hostname>:8080/Courses/?dept=DEPT'
Where DEPT is replaced with the dept Code such as COMP. This returns all of the Courses in that Dept in the Database.
  
### Login
Use '<hostname>:8080/Login/?username=USERNAME&password=PASSWORD'
Where USERNAME and PASSWORD are replaced with real usernames and passwords. This returns a 30 character Session Id that is needed for all user specific API calls. This method will return an empty array if username or password are not correct.
  
### SignUp
Use '<hostname>:8080/SignUp/?username=USERNAME&password=PASSWORD'
Where USERNAME and PASSWORD are replaced with new usernames and passwords. This also returns a 30 character Session Id that is needed for all user specific API calls. This will return an empty array if username is already taken.
  
