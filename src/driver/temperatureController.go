package driver

import (
	pidctl "DAQDriver/src/driver/pid"
	"fmt"
	"time"
)

type TemperatureController interface {
	initialize()
	startStaticHeater(float64, float64)
	stopStaticHeater()
	setupStaticHeater(pidctl.PIDSetting)
	setupTemperature(float64)
}

type TemperatureControllerInstance struct {
	vtRelation  []vtCoefficient
	datain      chan DAQDataCH
	powerout    chan DAQPowerCH
	dac         DACController
	heater      pidctl.PIDController
	baseVoltage float64
}

type vtCoefficient struct {
	start float64
	end   float64
	a     float64
	b     float64
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

}

func (tmpctl *TemperatureControllerInstance) vtmap(temperature float64) float64 {

	var v = float64(0)

	for i := 0; i < len(tmpctl.vtRelation); i++ {
		if temperature > tmpctl.vtRelation[i].start && temperature < tmpctl.vtRelation[i].end {
			v = tmpctl.vtRelation[i].a*temperature + tmpctl.vtRelation[i].b
			break
		}
	}

	fmt.Println("Core API: V-T mapping: ", temperature, "->", v)

	return v

}

func (tmpctl *TemperatureControllerInstance) startStaticHeater(basevoltage float64, targetTemperature float64) {

	fmt.Println("Core API: Starting heater temperature controller: ")
	tmpctl.heater.Reset()
	tmpctl.baseVoltage = basevoltage
	tmpctl.heater.Target = tmpctl.vtmap(targetTemperature) + tmpctl.baseVoltage
	fmt.Println("Core API: Base Voltage: ", tmpctl.baseVoltage)
	fmt.Println("Core API: Target Voltage: ", tmpctl.heater.Target)

	go tmpctl.heating(tmpctl.datain, tmpctl.powerout)
	fmt.Println("Core API: Started heater temperature controller: ")
}

func (tmpctl *TemperatureControllerInstance) stopStaticHeater() {

	helperCHSign = helperCHSign & 0xFD
	powerCHSign = 0
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

func newTemperatureController(dacctl DACController, dchinput chan DAQDataCH, pchoutput chan DAQPowerCH) *TemperatureController {

	var PIDCTL pidctl.PIDController = *pidctl.NewPIDController()
	var TMPCTL TemperatureController = &TemperatureControllerInstance{dac: dacctl, heater: PIDCTL, datain: dchinput, powerout: pchoutput}

	TMPCTL.initialize()

	return &TMPCTL

}
