package mysql

type Student struct {
	StudentId int    `json:"student_id,string"`
	Name      string `json:"name"`
	Stream    string `json:"stream,omitempty"`
	Email_id  string `json:"email_id"`
	Grade     int    `json:"grade"`
	Address   string `json:"address"`
}
