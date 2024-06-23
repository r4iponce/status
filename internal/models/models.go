package models

type Status struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Success     bool   `json:"success"`
	Error       string `json:"error"`
	Target      string `json:"target"`
}
