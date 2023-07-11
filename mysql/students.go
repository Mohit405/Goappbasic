package mysql

type Student struct {
	StudentId int    `json:"student_id"`
	Name      string `json:"name"`
	Stream    string `json:"stream"`
}
