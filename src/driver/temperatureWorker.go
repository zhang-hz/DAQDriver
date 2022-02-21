package driver

import (
	"fmt"
	"time"
)

func (tmpctl *TemperatureControllerInstance) heating(din <-chan DAQDataCH, pout chan<- DAQPowerCH) {

	for len(din) > 0 {
		<-din
	}
	helperCHSign = helperCHSign | 0x2

	for i := int64(0); i < dataChDepth+100; i++ {
		<-din
	}
	powerCHSign = 1

	var dtmp DAQDataCH
	var ptmp DAQPowerCH

	sleepTime := time.Duration(tmpctl.heater.Interval/2) * time.Nanosecond
	dacInterval := time.Duration(50) * time.Microsecond
	for {

		if helperCHSign&0x2 == 0 || powerCHSign == 0 {
			helperCHSign = helperCHSign & 0xFD
			powerCHSign = 0
			fmt.Print("Core API: Stop heater temperature controller \n")
			return
		}

		dtmp = <-din
		heater := tmpctl.heater.Linear(dtmp.directv[1])
		//fmt.Println(heater)

		if heater != 0 {
			tmpctl.dac.setDACVoltage("TP2", heater)
			time.Sleep(dacInterval)
			tmpctl.dac.setDACVoltage("TP1", heater)
		}
		time.Sleep(sleepTime)
		ptmp.heaterv[0] = heater
		ptmp.heaterv[1] = heater
		ptmp.heaterv[2] = heater

		ptmp.heaterp = [3]float64{0, 0, 0}

		if int64(len(pout)) < dataChDepth-1 {
			pout <- ptmp
		}

	}

}
