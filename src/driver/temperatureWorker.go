package driver

import (
	"fmt"
	"time"
)

func (tmpctl *TemperatureControllerInstance) heating(din <-chan DAQDataCH, heaterInfo *HeaterInfo) {

	for len(din) > 0 {
		<-din
	}
	helperCHSign = helperCHSign | 0x2

	for i := int64(0); i < dataChDepth+100; i++ {
		<-din
	}

	var dtmp DAQDataCH

	//sleepTime := time.Duration(tmpctl.heater.Interval/2) * time.Nanosecond
	dacInterval := time.Duration(20) * time.Microsecond
	for {

		if helperCHSign&0x2 == 0 {
			helperCHSign = helperCHSign & 0xFD
			heaterInfo.voltage[0] = 0
			heaterInfo.voltage[1] = 0
			heaterInfo.power[0] = 0
			heaterInfo.power[1] = 0
			fmt.Print("Core API: Stop heater temperature controller \n")
			return
		}

		dtmp = <-din
		heater := tmpctl.heater.Linear(dtmp.directv[0] - dtmp.directv[1])
		//fmt.Println(heater)

		if heater != 0 {
			tmpctl.dac.setDACVoltage("TP2", heater)
			time.Sleep(dacInterval)
			//tmpctl.dac.setDACVoltage("TP1", heater)
			//time.Sleep(dacInterval)

			heaterInfo.voltage[0] = heater
			heaterInfo.voltage[1] = heater
		}
		//time.Sleep(sleepTime)

		heaterInfo.power[0] = 0
		heaterInfo.power[1] = 0

	}

}
