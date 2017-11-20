package main



type User struct {
  Id      int     `json:"id"`
  User    string  `json:"user"`
  Pass    string  `json:"pass"`
  SemLeft int     `json:"semLeft"`
}
