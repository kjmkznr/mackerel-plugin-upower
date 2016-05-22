mackerel-plugin-upower
======================

## Description

Get battery metrics from UPower for Mackerel.

## Usage

### Build this program

```
go get github.com/kjmkznr/mackerel-plugin-upower
cd $GO_HOME/src/github.com/kjmkznr/mackerel-plugin-upower
go build
```

### Execute this plugin

```
./mackerel-plugin-upower
upower.voltage.BAT0.voltage     12.486000       1463929148
upower.voltage.BAT1.voltage     12.536000       1463929148
upower.state.BAT0.state 4       1463929148
upower.state.AC.state   0       1463929148
upower.state.BAT1.state 4       1463929148
upower.energy.BAT1.current      77.420000       1463929148
upower.energy.BAT0.current      23.720000       1463929148
upower.energy.BAT0.full 23.740000       1463929148
upower.energy.BAT1.full 77.560000       1463929148
upower.energy.BAT1.full_design  71.280000       1463929148
upower.energy.BAT0.full_design  23.200000       1463929148
upower.energy.BAT1.full_design  71.280000       1463929148
upower.energy.BAT0.full_design  23.200000       1463929148
upower.energy.BAT1.rate 3.944000        1463929148
upower.energy.BAT0.rate 6.248000        1463929148
```

### Add mackerel-agent.conf

For example is below.

```
[plugin.metrics.upower]
command = "/path/to/mackerel-plugin-upower"
type = "metric"
```

## Author

[KOJIMA Kazunori](https://github.com/kjmkznr)
