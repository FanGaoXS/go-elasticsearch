package entity

type User struct {
	UUID string   `json:"uuid"`
	Name string   `json:"name"`
	Age  int64    `json:"age"`
	Desc string   `json:"desc"`
	Tags []string `json:"tags"`
}
