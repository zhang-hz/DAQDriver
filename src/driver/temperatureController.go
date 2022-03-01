package driver

import (
	pidctl "DAQDriver/src/driver/pid"
	"fmt"
	"time"
)

type TemperatureController interface {
	initialize()
	startProgHeater(float64, float64, float64)
	stopProgHeater()
	startStaticHeater(float64, float64)
	stopStaticHeater()
	setupStaticHeater(pidctl.PIDSetting)
	setupTemperature(float64)
}

type TemperatureControllerInstance struct {
	vtRelation      []vtCoefficient
	datain          chan DAQDataCH
	dac             DACController
	heater          pidctl.PIDController
	info            *HeaterInfo
	baseVoltage     float64
	progHeatingInfo progHeatingSetting
}

type vtCoefficient struct {
	start float64
	end   float64
	a     float64
	b     float64
}

type progHeatingSetting struct {
	startTime       int64
	speed           float64
	baseTemperature float64
	maxTemperature  float64
}

func (tmpctl *TemperatureControllerInstance) initialize() {

	tmpctl.vtRelation = make([]vtCoefficient, 1)
	tmpctl.vtRelation[0] = vtCoefficient{
		start: 0,
		end:   600,
		a:     27.99268e6,
		b:     -611.16518e6,
	}
	tmpctl.baseVoltage = 0

	tmpctl.progHeatingInfo = progHeatingSetting{startTime: 0, speed: 0, baseTemperature: 0}

}

func (tmpctl *TemperatureControllerInstance) vtmap(temperature float64) float64 {

	var v = float64(0)

	for i := 0; i < len(tmpctl.vtRelation); i++ {
		if temperature > tmpctl.vtRelation[i].start && temperature < tmpctl.vtRelation[i].end {
			v = tmpctl.vtRelation[i].a*temperature + tmpctl.vtRelation[i].b
			break
		}
	}

	//fmt.Println("Core API: V-T mapping: ", temperature, "->", v)

	return v

}

func (tmpctl *TemperatureControllerInstance) progVTMap() float64 {

	timeNow := time.Now().UnixMicro()
	progTemp := tmpctl.progHeatingInfo.speed*float64(timeNow-tmpctl.progHeatingInfo.startTime) + tmpctl.progHeatingInfo.baseTemperature
	if progTemp > tmpctl.progHeatingInfo.maxTemperature {
		progTemp = tmpctl.progHeatingInfo.maxTemperature
	}
	return tmpctl.vtmap(progTemp)

}

func (tmpctl *TemperatureControllerInstance) startProgHeater(basevoltage float64, heatingSpeed float64, baseTemperature float64) {

	fmt.Println("Core API: Starting program heating temperature controller: ")
	tmpctl.dac.setDACVoltage("TP1", 0)
	time.Sleep(time.Duration(1) * time.Millisecond)
	tmpctl.heater.Reset()
	tmpctl.baseVoltage = basevoltage
	tmpctl.progHeatingInfo.speed = heatingSpeed
	tmpctl.progHeatingInfo.baseTemperature = baseTemperature
	tmpctl.progHeatingInfo.maxTemperature = 150
	tmpctl.progHeatingInfo.startTime = time.Now().UnixMicro()
	go tmpctl.progHeating(tmpctl.datain, tmpctl.info)
	fmt.Println("Core API: Started program heating temperature controller: ")

}

func (tmpctl *TemperatureControllerInstance) stopProgHeater() {

	fmt.Println("Core API: Stopping program heating temperature controller: ")
	helperCHSign = helperCHSign & 0xFD
	time.Sleep(time.Duration(1) * time.Millisecond)
	tmpctl.progHeatingInfo.speed = 0
	tmpctl.dac.setDACVoltage("TP2", 0)
	time.Sleep(time.Duration(1) * time.Millisecond)
	tmpctl.dac.setDACVoltage("TP1", 0)
	fmt.Println("Core API: Stopped program heating temperature controller: ")

}

func (tmpctl *TemperatureControllerInstance) startStaticHeater(basevoltage float64, targetTemperature float64) {

	fmt.Println("Core API: Starting heater temperature controller: ")
	tmpctl.dac.setDACVoltage("TP1", 0)
	time.Sleep(time.Duration(1) * time.Millisecond)
	tmpctl.heater.Reset()
	tmpctl.baseVoltage = basevoltage
	tmpctl.heater.Target = tmpctl.vtmap(targetTemperature) + tmpctl.baseVoltage
	fmt.Println("Core API: V-T mapping: ", targetTemperature, "->", tmpctl.heater.Target)
	fmt.Println("Core API: Base Voltage: ", tmpctl.baseVoltage)
	fmt.Println("Core API: Target Voltage: ", tmpctl.heater.Target)

	go tmpctl.heating(tmpctl.datain, tmpctl.info)
	fmt.Println("Core API: Started heater temperature controller: ")
}

func (tmpctl *TemperatureControllerInstance) stopStaticHeater() {

	helperCHSign = helperCHSign & 0xFD
	tmpctl.dac.setDACVoltage("TP2", 0)
	time.Sleep(time.Duration(1) * time.Millisecond)
	tmpctl.dac.setDACVoltage("TP1", 0)

}

func (tmpctl *TemperatureControllerInstance) setupStaticHeater(pidsetting pidctl.PIDSetting) {

	tmpctl.heater.Setup(pidsetting)

}

func (tmpctl *TemperatureControllerInstance) setupTemperature(temperature float64) {

	fmt.Println("Core API: Setup temperature: ", temperature)
	tmpctl.heater.Target = tmpctl.vtmap(temperature) + tmpctl.baseVoltage
	fmt.Println("Core API: Base Voltage: ", tmpctl.baseVoltage)
	fmt.Println("Core API: Target Voltage: ", tmpctl.heater.Target)

}

func newTemperatureController(dacctl DACController, dchinput chan DAQDataCH, info *HeaterInfo) *TemperatureController {

	var PIDCTL pidctl.PIDController = *pidctl.NewPIDController()
	var TMPCTL TemperatureController = &TemperatureControllerInstance{dac: dacctl, heater: PIDCTL, datain: dchinput, info: info}

	TMPCTL.initialize()

	return &TMPCTL

}
