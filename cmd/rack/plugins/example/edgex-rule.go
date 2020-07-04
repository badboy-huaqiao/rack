// +build ignore
package main

type EdgeXRule struct {
	// SetupEdgeXSource() (config model.EdgeXConnectionConfig)
	// AddDevicesSource(sources []models.Device, f func(string) bool) []models.Device
	// Process(data <-chan models.Event)
	// AddSink(sources []models.Command)
}

// func (e EdgeXRule) SetupEdgeXSource() (config model.EdgeXConnectionConfig) {
// 	config.Host = "192.168.56.4"
// 	return config
// }

// func (e EdgeXRule) AddSource(sources []models.Device) (result []map[string][]string) {
// 	for _, d := range sources {
// 		if d.Name == "device-01" {
// 			for _, p := range d.Profile.DeviceResources {
// 				if p.Name == "temple" {

// 				}
// 			}
// 		}
// 	}
// }
