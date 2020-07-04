package source

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"rack/pkg/template"
	"strings"
	"time"

	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

const (
	devicesPath  = ":48081/api/v1/device"
	commandsPath = ":48082/api/v1/device/name/"
)

type EdgeXSource struct {
	Config template.EdgeXConnectionConfig
}

func (e EdgeXSource) Source() {

}

func (e EdgeXSource) SendCommand(device string, property string, param map[string]interface{}) (result models.Event, err error) {
	cmds, err := e.Commands(device)
	if err != nil {
		return result, err
	}
	var cmdUrl string
	for _, cmd := range cmds {
		if cmd.Name == property && param != nil {
			cmdUrl = cmd.Put.URL
			break
		} else {
			cmdUrl = cmd.Get.URL
			break
		}
	}
	cmdUrl = strings.Replace(cmdUrl, "edgex-core-command", e.Config.Host, 1)
	if param == nil {
		if result, err = sendGetCmd(cmdUrl); err != nil {
			return result, err
		}
	} else {
		if err = sendPutCmd(cmdUrl, param); err != nil {
			return result, err
		}
	}
	return result, nil
}

func sendGetCmd(url string) (event models.Event, err error) {
	resp, err := httpRequest(url, nil)
	if err != nil {
		return event, err
	}
	if err = json.NewDecoder(resp.Body).Decode(&event); err != nil {
		return event, err
	}
	return event, nil
}

func sendPutCmd(url string, param map[string]interface{}) (err error) {
	data, err := json.Marshal(param)
	if err != nil {
		return err
	}
	_, err = httpRequest(url, data)
	if err != nil {
		return err
	}
	return nil
}

func httpRequest(url string, body []byte) (resp *http.Response, err error) {
	client := &http.Client{Timeout: 5 * time.Second}
	var r *http.Request
	if body != nil {
		paramBody := bytes.NewReader(body)
		r, err = http.NewRequest(http.MethodPut, url, paramBody)
	} else {
		r, err = http.NewRequest(http.MethodGet, url, nil)
	}

	if err != nil {
		return nil, err
	}
	if resp, err = client.Do(r); err != nil {
		return nil, err
	}
	return resp, nil
}

func (e EdgeXSource) Commands(device string) (commands []models.Command, err error) {
	url := fmt.Sprintf("http://%s%s%s", e.Config.Host, commandsPath, device)
	res, err := httpRequest(url, nil)
	if err != nil {
		return nil, err
	}
	if err = json.NewDecoder(res.Body).Decode(&commands); err != nil {
		return nil, err
	}
	return commands, nil
}

func (e EdgeXSource) Devices() (devices []models.Device, err error) {
	url := fmt.Sprintf("http://%s%s", e.Config.Host, devicesPath)
	res, err := httpRequest(url, nil)
	if err != nil {
		return nil, err
	}
	if err = json.NewDecoder(res.Body).Decode(&devices); err != nil {
		return nil, err
	}
	return devices, nil
}

type EdgeXSink struct {
}

func (es EdgeXSource) Filter(d []string, f func(string) bool) []string {
	var devices []string
	for _, v := range d {
		if f(v) {
			d = append(d, v)
		}
	}
	return devices
}

func (es EdgeXSink) Sink() {

}
