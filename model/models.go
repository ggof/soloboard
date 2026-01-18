package model

type Task struct {
	ID          int
	Name        string
	Description string
}

type Section struct {
	ID    int
	Name  string
	Tasks []Task
}

type Board struct {
	ID       int
	Name     string
	Sections []Section
}
