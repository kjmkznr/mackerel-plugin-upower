package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/godbus/dbus"
	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mackerelio/mackerel-agent/logging"
)

const (
	UPOWER_BUS_NAME         = "org.freedesktop.UPower"
	UPOWER_DEVICE_BUS_NAME  = "org.freedesktop.UPower.Device"
	UPOWER_BASE_OBJECT_PATH = "/org/freedesktop/UPower"
)

const (
	DEVICE_TYPE_UNKNOWN = iota
	DEVICE_TYPE_LINE_POWER
	DEVICE_TYPE_BATTERY
	DEVICE_TYPE_UPS
	DEVICE_TYPE_MONITOR
	DEVICE_TYPE_MOUSE
	DEVICE_TYPE_KEYBOARD
	DEVICE_TYPE_PDA
	DEVICE_TYPE_PHONE
)

var powerProperties = []string{".NativePath", ".Type", ".State"}
var batteryProperties = []string{".EnergyRate", ".EnergyFull", ".Energy", ".EnergyFullDesign", ".Voltage"}

var logger = logging.GetLogger("metrics.plugin.upower")

type UPowerPlugin struct {
}

type DeviceProperties struct {
	NativePath       string
	DeviceType       uint32
	EnergyFullDesign float64
	EnergyFull       float64
	Energy           float64
	EnergyRate       float64
	Voltage          float64
	State            uint32
}

func getDeviceProperties(conn *dbus.Conn, path dbus.ObjectPath) *DeviceProperties {
	device := conn.Object(UPOWER_BUS_NAME, path)

	result := &DeviceProperties{}

	for _, p := range powerProperties {
		val, err := device.GetProperty(UPOWER_DEVICE_BUS_NAME + p)

		if err != nil {
			logger.Warningf("Failed to get property: '%v'", err)
		}

		switch p {
		case powerProperties[0]:
			result.NativePath = val.Value().(string)
		case powerProperties[1]:
			result.DeviceType = val.Value().(uint32)
		case powerProperties[2]:
			result.State = val.Value().(uint32)
		}
	}

	if result.DeviceType == DEVICE_TYPE_BATTERY {
		for _, p := range batteryProperties {
			val, err := device.GetProperty(UPOWER_DEVICE_BUS_NAME + p)

			if err != nil {
				logger.Warningf("Failed to get property: '%v'", err)
			}

			switch p {
			case batteryProperties[0]:
				result.EnergyRate = val.Value().(float64)
			case batteryProperties[1]:
				result.EnergyFull = val.Value().(float64)
			case batteryProperties[2]:
				result.Energy = val.Value().(float64)
			case batteryProperties[3]:
				result.EnergyFullDesign = val.Value().(float64)
			case batteryProperties[4]:
				result.Voltage = val.Value().(float64)
			}
		}
	}

	return result
}

func (p UPowerPlugin) GraphDefinition() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		"upower.energy.#": mp.Graphs{
			Label: "Energy",
			Unit:  "float",
			Metrics: []mp.Metrics{
				mp.Metrics{Name: "current", Label: "Current Wh"},
				mp.Metrics{Name: "full", Label: "Full Wh"},
				mp.Metrics{Name: "full_design", Label: "Full Design Wh"},
				mp.Metrics{Name: "rate", Label: "Rate W"},
			},
		},
		"upower.voltage.#": mp.Graphs{
			Label: "Battery Voltage",
			Unit:  "float",
			Metrics: []mp.Metrics{
				mp.Metrics{Name: "voltage", Label: "Voltage"},
			},
		},
		"upower.state.#": mp.Graphs{
			Label: "Power State",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				mp.Metrics{Name: "state", Label: "State"},
			},
		},
	}
}

func (p UPowerPlugin) FetchMetrics() (map[string]interface{}, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		logger.Warningf("Failed to connect to session bus: '%v'", err)
		return nil, err
	}

	busObject := conn.Object(UPOWER_BUS_NAME, UPOWER_BASE_OBJECT_PATH)

	var devices []dbus.ObjectPath
	err = busObject.Call(UPOWER_BUS_NAME+".EnumerateDevices", 0).Store(&devices)
	if err != nil {
		logger.Warningf("Failed to enumerate devices: '%v'", err)
		return nil, err
	}

	result := make(map[string]interface{})
	for _, d := range devices {
		device := getDeviceProperties(conn, d)
		result["upower.state."+device.NativePath+".state"] = device.State
		if device.DeviceType == DEVICE_TYPE_BATTERY {
			result["upower.energy."+device.NativePath+".current"] = device.Energy
			result["upower.energy."+device.NativePath+".full"] = device.EnergyFull
			result["upower.energy."+device.NativePath+".full_design"] = device.EnergyFullDesign
			result["upower.energy."+device.NativePath+".rate"] = device.EnergyRate
			result["upower.voltage."+device.NativePath+".voltage"] = device.Voltage
		}
	}

	return result, nil
}

func doMain(c *cli.Context) error {
	upower := UPowerPlugin{}
	helper := mp.NewMackerelPlugin(upower)
	helper.Run()

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "mackerel-plugin-upower"
	app.Version = version
	app.Usage = "Get metrics from UPower."
	app.Author = "KOJIMA Kazunori"
	app.Email = "kjm.kznr@gmail.com"
	app.Action = doMain

	app.Run(os.Args)
}
