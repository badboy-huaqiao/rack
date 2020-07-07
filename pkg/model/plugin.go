package model

// type Concrete struct {
// 	Name    string `json:"Concrete"`
// 	Version string `json:"version"`
// }

type Plugin struct {
	Concrete string `json:"Concrete"`
	Version  string `json:"version"`
	Path     string `json:"path"`
}
