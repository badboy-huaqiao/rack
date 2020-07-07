package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"rack/pkg/internal/core"
	"rack/pkg/model"
	"strings"
)

const (
	pluginPath   = "plugins"
	pluginSuffix = ".so"
)

func Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var plug model.Plugin
	if err := json.NewDecoder(r.Body).Decode(&plug); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func Upload(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	source, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer source.Close()
	if !strings.HasSuffix(header.Filename, pluginSuffix) {
		http.Error(w, errors.New("unsupported file type.").Error(), http.StatusBadRequest)
		return
	}
	if err := os.Mkdir(pluginPath, os.ModePerm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	path := filepath.Join(pluginPath, header.Filename)

	dst, err := os.Create(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := io.Copy(dst, source); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// plugin := model.Plugin{
	// 	Name:    strings.TrimSuffix(header.Filename, pluginSuffix),
	// 	Version: r.FormValue("version"),
	// 	Path:    path,
	// }
}

func Load(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var plug model.Plugin
	if err := json.NewDecoder(r.Body).Decode(&plug); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	go core.Load(plug)
}
