package main



type Course struct {
  Id      int    `json:"id"`
  Hours   int     `json:"hours"`
  Number  int     `json:"number"`
  Pid     int     `json:"pid"`
}

type Courses []Course


type CourseDept struct {
  Id      int    `json:"id"`
  Hours   int     `json:"hours"`
  Number  int     `json:"number"`
  dept    string  `json:"dept"`
}

type CoursesDept []CourseDept
