package main

import (
	"math"

	"github.com/distatus/battery"
)

type batterySimple struct {
	props properties
	env   environmentInfo

	battery.Battery
	Pct  int
	Icon string
}

const (
	//NewProp enables something
	NewProp Property = "newprop"
)

func (bs *batterySimple) enabled() bool {
	batteries, err := bs.env.getBatteryInfo()

	if err != nil || len(batteries) == 0 {
		return false
	}
	for _, bt := range batteries {
		bs.Battery.Full += bt.Full
		bs.Battery.Current += bt.Current
		bs.Battery.State = bs.newTotalState(bt.State)
	}
	bs.Pct = int(math.Min(100, bs.Battery.Current/bs.Battery.Full*100))

	switch bs.Battery.State {
	case battery.Discharging, battery.NotCharging:
		bs.Icon = bs.props.getString(DischargingIcon, "")
	case battery.Charging:
		bs.Icon = bs.props.getString(ChargingIcon, "")
	case battery.Full:
		bs.Icon = bs.props.getString(ChargedIcon, "")
		/* case battery.Unknown, battery.Empty:
		return true */
	}
	return true
}

func (bs *batterySimple) newTotalState(newState battery.State) battery.State {
	switch bs.State {
	case battery.NotCharging, battery.Discharging:
		return battery.Discharging
	case battery.Empty, battery.Unknown, battery.Full:
		return newState
	case battery.Charging:
		if newState == battery.Discharging {
			return battery.Discharging
		}
		return battery.Charging
	}
	return newState
}

func (bs *batterySimple) string() string {
	segmentTemplate := bs.props.getString(SegmentTemplate, "{hahah {.Icon}{.Pct}}")
	template := &textTemplate{
		Template: segmentTemplate,
		Context:  bs,
		Env:      bs.env,
	}
	text, err := template.render()
	if err != nil {
		return err.Error()
	}
	return text
}

func (bs *batterySimple) init(props properties, env environmentInfo) {
	bs.props = props
	bs.env = env
}
