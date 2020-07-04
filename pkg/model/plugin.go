package model

type Plugin struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Path    string `json:"path"`
}
