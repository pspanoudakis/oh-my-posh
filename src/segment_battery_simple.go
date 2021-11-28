package main

import (
	"math"

	"github.com/distatus/battery"
)

type batterySimple struct {
	props properties
	env   environmentInfo

	battery.Battery
	BatteryIcon string
	Pct         int
	StatusIcon  string
}

const (
	UseLevelIcons   Property = "use_level_icons"
	Icon            Property = "icon"
	FullLevelIcon   Property = "full_lvl_icon"
	GoodLevelIcon   Property = "good_lvl_icon"
	MediumLevelIcon Property = "medium_lvl_icon"
	LowLevelIcon    Property = "low_lvl_icon"
	EmptyLevelIcon  Property = "empty_lvl_icon"
	// ChargingIcon to display when charging
	ChargingIcn     Property = "charging_icon"
	ChargingFGColor Property = "charging_fg_color"
	ChargingBGColor Property = "charging_bg_color"
	// DischargingIcon o display when discharging
	DischargingIcn     Property = "discharging_icon"
	DischargingFGColor Property = "discharging_fg_color"
	DischargingBGColor Property = "discharging_bg_color"
	// ChargedIcon to display when fully charged
	ChargedIcn     Property = "charged_icon"
	ChargedFGColor Property = "charged_fg_color"
	ChargedBGColor Property = "charged_bg_color"
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
		bs.StatusIcon = bs.props.getString(DischargingIcn, "")
	case battery.Charging:
		bs.StatusIcon = bs.props.getString(ChargingIcn, "")
	case battery.Full:
		bs.StatusIcon = bs.props.getString(ChargedIcn, "")
		/* case battery.Unknown, battery.Empty:
		return true */
	}
	bs.setColors()
	bs.setIcon()
	return true
}

func (bs *batterySimple) setIcon() {
	if bs.props.getBool(UseLevelIcons, false) {
		if bs.Pct == 100 {
			bs.BatteryIcon = bs.props.getString(FullLevelIcon, "")
		} else if bs.Pct >= 65 {
			bs.BatteryIcon = bs.props.getString(GoodLevelIcon, "")
		} else if bs.Pct >= 35 {
			bs.BatteryIcon = bs.props.getString(MediumLevelIcon, "")
		} else if bs.Pct >= 10 {
			bs.BatteryIcon = bs.props.getString(LowLevelIcon, "")
		} else {
			bs.BatteryIcon = bs.props.getString(EmptyLevelIcon, "")
		}
	} else {
		bs.BatteryIcon = bs.props.getString(Icon, "")
	}
}

func (bs *batterySimple) setColors() {
	if !bs.props.hasOneOf(ChargedFGColor, ChargedBGColor, ChargingFGColor, ChargingBGColor, DischargingFGColor, DischargingBGColor) {
		return
	}
	var fgColor Property
	var bgColor Property
	switch bs.Battery.State {
	case battery.Discharging, battery.NotCharging:
		fgColor = DischargingFGColor
		bgColor = DischargingBGColor
	case battery.Charging:
		fgColor = ChargingFGColor
		bgColor = ChargingBGColor
	case battery.Full:
		fgColor = ChargedFGColor
		bgColor = ChargedBGColor
	case battery.Empty, battery.Unknown:
		return
	}
	bs.props[BackgroundOverride] = bs.props.getColor(bgColor, bs.props.getColor(BackgroundOverride, ""))
	bs.props[ForegroundOverride] = bs.props.getColor(fgColor, bs.props.getColor(ForegroundOverride, ""))
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
	segmentTemplate := bs.props.getString(SegmentTemplate, "{{.StatusIcon}}{{.BatteryIcon}}{{.Pct}}%")
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
