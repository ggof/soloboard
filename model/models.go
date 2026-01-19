package model

import "github.com/oklog/ulid/v2"

type Task struct {
	ID          string
	Name        string
	Description string
}

type Section struct {
	ID    string
	Name  string
	Tasks []Task
}

type Board struct {
	ID       string
	Name     string
	Sections []Section
}

func NewTask(name string, description string) Task {
	return Task{ID: ulid.Make().String(), Name: name, Description: description}
}

func NewSection(name string) Section {
	return Section{ID: ulid.Make().String(), Name: name}
}

func NewBoard(name string) Board {
	return Board{
		ID:   ulid.Make().String(),
		Name: name,
		Sections: []Section{
			NewSection("TODO"),
			NewSection("IN PROGRESS"),
			NewSection("DONE"),
		},
	}
}	
