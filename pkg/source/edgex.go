package source

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/edgexfoundry/go-mod-messaging/pkg/types"
	zmq "github.com/pebbe/zmq4"
)

const (
	devicesPath  = ":48081/api/v1/device"
	commandsPath = ":48082/api/v1/device/name/"
)

type EdgeXConnectionConfig struct {
	Host string
	// Port int
}

type EdgeXSource struct {
	Config     EdgeXConnectionConfig
	DataStream chan models.Event
}

var edgexSource *EdgeXSource

func GetEdgeXSource() *EdgeXSource {
	if edgexSource == nil {
		edgexSource = &EdgeXSource{}
		edgexSource.DataStream = make(chan models.Event, 100)
		return edgexSource
	}
	return edgexSource
}

func (e EdgeXSource) Source() {

}

func EdgeXDataStream() {
	q, _ := zmq.NewSocket(zmq.SUB)
	defer q.Close()
	if err := q.Connect("tcp://192.168.56.4:5563"); err != nil {
		fmt.Printf("Error connect message: %s\n", err.Error())
	}
	q.SetSubscribe("")

	for {
		msg, err := q.RecvMessage(0)
		if err != nil {
			id, _ := q.GetIdentity()
			fmt.Printf("Error getting message %s\n", id)
		} else {
			var envelope types.MessageEnvelope
			var event models.Event
			if err := json.Unmarshal([]byte(msg[1]), &envelope); err != nil {
				fmt.Printf("Error getting message: %s\n", err.Error())
			}
			if err := json.Unmarshal(envelope.Payload, &event); err != nil {
				fmt.Printf("Error getting message: %s\n", err.Error())
			}
			fmt.Printf("getting message %v\n", event)
			if len(GetEdgeXSource().DataStream) == cap(GetEdgeXSource().DataStream) {
				<-GetEdgeXSource().DataStream
			}
			GetEdgeXSource().DataStream <- event
		}
	}

}

func DataPipeline() chan<- models.Event {
	dataStream := make(chan<- models.Event, 100)

	return dataStream
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

	fmt.Printf("send Get Cmd:%s, result: %v\n", url, event)

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
	fmt.Printf("send Put Cmd:%s success\n", url)
	return nil
}

func httpRequest(url string, body []byte) (resp *http.Response, err error) {
	client := &http.Client{Timeout: 5 * time.Second}
	var r *http.Request
	if body != nil {
		paramBody := bytes.NewReader(body)
		if r, err = http.NewRequest(http.MethodPut, url, paramBody); err != nil {
			return nil, err
		}
	} else {
		if r, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
			return nil, err
		}
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
