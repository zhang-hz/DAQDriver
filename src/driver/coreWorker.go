package driver

import (
	"fmt"
	"time"
)

type DAQDataCH struct {
	directv [4]float64
}

type DAQPowerCH struct {
	heaterv [3]float64
	heaterp [3]float64
}

func (corectl *CoreController) fetchData(dout chan<- DAQDataCH) {

	fmt.Print("Start Fetch Data\n")

	var tmp [4]float64

	for {
		if !runningDAQ {
			fmt.Print("Stop Fetch Data\n")
			return
		}

		tmp = corectl.adc.read()
		dout <- DAQDataCH{directv: [4]float64{-1 * tmp[0], -1 * tmp[1], tmp[2], tmp[3]}}

	}
}

func interconnectHub(din <-chan DAQDataCH, pin <-chan DAQPowerCH, socket chan<- socketCH, dout1 chan<- DAQDataCH, dout2 chan<- DAQDataCH) {

	var socketdata socketCH
	var num = int64(0)
	var downsample = int64(0)
	var heaterDownSample = int64(0)
	var data = DAQDataCH{}
	var power = DAQPowerCH{[3]float64{0, 0, 0}, [3]float64{0, 0, 0}}

	for {
		//fmt.Print(helperCHSign, "\n")
		data = <-din
		diffv := data.directv[0] - data.directv[1]

		if powerCHSign == 1 && len(pin) > 0 {
			power = <-pin
		}

		if downsample >= socketDownSampleRate-1 {

			downsample = 0

			if socketCHSign == 1 || (socketCHSign == 0 && (int64(len(socket)) < socketChDepth-1)) {

				if num >= socketDataLength {
					socket <- socketdata
					num = 0
					socketdata.time = time.UnixMicro(time.Now().UnixMicro())
					socketdata.interval = (1e9 * float64(socketDownSampleRate)) / 50e3
					socketdata.length = socketDataLength
				}

				socketdata.directv[0][num] = data.directv[0]
				socketdata.directv[1][num] = data.directv[1]
				socketdata.directv[2][num] = data.directv[2]
				socketdata.directv[3][num] = data.directv[3]
				socketdata.diffv[num] = diffv

				socketdata.heaterv[0][num] = power.heaterv[0]
				socketdata.heaterv[1][num] = power.heaterv[1]
				socketdata.heaterv[2][num] = power.heaterv[2]

				socketdata.heaterp[0][num] = power.heaterp[0]
				socketdata.heaterp[1][num] = power.heaterp[1]
				socketdata.heaterp[2][num] = power.heaterp[2]

				num++
			}
		} else {
			downsample++
		}

		if helperCHSign&0x1 != 0 && int64(len(dout1)) < helperChDepth-1 {
			dout1 <- data
		}

		if (helperCHSign&0x2 != 0 && int64(len(dout2)) < helperChDepth-1) || (helperCHSign&0x2 == 0 && (int64(len(dout2)) < helperChDepth-1)) {
			if heaterDownSample > heaterDownSampleRate-1 {
				dout2 <- data
				heaterDownSample = 0
			} else {
				heaterDownSample++
			}

		}

	}

}
