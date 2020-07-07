package template

import (
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

type EdgeX interface {
	AddDevicesSource(sources []models.Device, f func(string) bool) []models.Device
	Process(data <-chan models.Event)
	AddSink(sources []models.Command)

	// Process2( func input()) func output()
}
