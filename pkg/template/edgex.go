package template

import (
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

type EdgeXConnectionConfig struct {
	Host string
	// Port int
}

type EdgeX interface {
	SetupEdgeXSource() (config EdgeXConnectionConfig)
	AddDevicesSource(sources []models.Device, f func(string) bool) []models.Device
	Process(data <-chan models.Event)
	AddSink(sources []models.Command)

	// Process2( func input()) func output()
}
