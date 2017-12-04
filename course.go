package main



type Course struct {
  Id      int    `json:"id"`
  Hours   int     `json:"hours"`
  Number  int     `json:"number"`
  Program string  `json:"program"`
}

type Courses []Course
