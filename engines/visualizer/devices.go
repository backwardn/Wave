package visualizer

import (
	log "github.com/Sirupsen/logrus"
	"github.com/hkparker/Wave/models"
	"strings"
)

func updateKnownDevices(frame models.Wireless80211Frame) {
	if _, ok := Devices[frame.Address1]; !ok {
		registerNewDevice(frame.Address1)
	} else if _, ok := Devices[frame.Address2]; !ok {
		registerNewDevice(frame.Address2)
	} else if _, ok := Devices[frame.Address3]; !ok {
		registerNewDevice(frame.Address3)
	} else if _, ok := Devices[frame.Address4]; !ok {
		registerNewDevice(frame.Address4)
	}
}

func registerNewDevice(mac string) {
	DevicesMux.Lock()
	defer DevicesMux.Unlock()
	if broadcast(mac) {
		return
	}
	vendor := "unknown"
	if vendor_string, ok := VendorBytes[strings.ToUpper(mac[0:8])]; ok {
		vendor = vendor_string
	}
	device := models.Device{
		MAC:    mac,
		Vendor: vendor,
	}
	Devices[mac] = device
	device.Save()
	visualizeNewDevice(device)
}

func broadcast(mac string) bool {
	if mac == "ff:ff:ff:ff:ff:ff" {
		return true
	}
	return false
}

func visualizeNewDevice(device models.Device) {
	new_resources := make(VisualEvent)
	new_resources["NewDevices"] = append(
		new_resources["NewDevices"],
		device.VisualData(),
	)
	VisualEvents <- new_resources
	log.WithFields(log.Fields{
		"at":  "visualizeNewDevice",
		"mac": device.MAC,
	}).Info("new device observed")
}
