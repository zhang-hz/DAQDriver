package driver

import (
	pidctl "DAQDriver/src/driver/pid"
	"fmt"
	"time"
)

/**
*
* 获取ADC平均电压
*
**/

type avgVoltageData struct {
	Voltage [4]float64 //平均电压
	Time    float64    //耗时
}

func (corectl *CoreController) AvgVoltage(SampleNumber int64) avgVoltageData {

	var result avgVoltageData
	rchan := make(chan avgVoltageData, 1)
	go avgVoltageExec(rchan, SampleNumber, corectl.dchhelper1)
	result = <-rchan
	close(rchan)

	return result
}

func avgVoltageExec(dout chan<- avgVoltageData, SampleNumber int64, din <-chan DAQDataCH) {

	fmt.Print("Start avg Data\n")
	var avg = [4]float64{0, 0, 0, 0}

	var i int64
	var n int64
	var data DAQDataCH

	for len(din) > 0 {
		<-din
	}

	helperCHSign = helperCHSign | 0x1

	for i = 0; i < dataChDepth+100; i++ {
		data = <-din
	}

	startTime := time.Now().UnixNano()

	for i = 0; i < SampleNumber; i++ {
		data = <-din
		for n = 0; n < 4; n++ {
			avg[n] = avg[n] + data.directv[n]
		}
	}
	helperCHSign = helperCHSign & 0xFE
	endTime := time.Now().UnixNano()

	for i = 0; i < 4; i++ {
		avg[i] = avg[i] / float64(SampleNumber)
	}

	result := avgVoltageData{
		Voltage: avg,
		Time:    float64(endTime - startTime),
	}

	dout <- result

	//wg.Done()
}

/**
*
*  设定ADC电压偏置
*
**/

func (corectl *CoreController) SetADCVos(ADCCH int64, VosNumber float64) {

	corectl.adc.setVos(ADCCH, VosNumber)

}

/**
*
*  设定DAC电压
*
**/

func (corectl *CoreController) SetDACVoltage(DACPort string, Voltage float64) float64 {

	return corectl.dac.setDACVoltage(DACPort, Voltage)

}
func (corectl *CoreController) GetDACVoltage(DACPort string) float64 {

	return corectl.dac.getDACVoltage(DACPort)

}

/**
*
*  设定DAC电压补偿
*
**/
func (corectl *CoreController) SetDACOffset(DACPort string, offsetVoltage float64) {

	corectl.dac.setDACOffset(DACPort, offsetVoltage)
}

/**
*
*  启动温度控制器
*
**/

func (corectl *CoreController) StartStaticHeater(basevoltage float64, targetTemperature float64) {

	corectl.tmp.startStaticHeater(basevoltage, targetTemperature)

}

/**
*
*  停止温度控制器
*
**/

func (corectl *CoreController) StopStaticHeater() {
	corectl.tmp.stopStaticHeater()
}

func (corectl *CoreController) SetupTemperature(temperature float64) {

	corectl.tmp.setupTemperature(temperature)
}

/**
*
*  停止温度控制器参数
*
**/

func (corectl *CoreController) SetupStaticHeater(kp float64, ki float64, kd float64, tolerance float64, errorTolerance float64) {

	pidsetting := pidctl.PIDSetting{
		Kp:             kp,
		Ki:             ki,
		Kd:             kd,
		Interval:       20000 * float64(heaterDownSampleRate),
		Tau:            20000 * float64(heaterDownSampleRate),
		Tolerance:      tolerance,
		ErrorTolerance: errorTolerance,
		LimitMin:       0,
		LimitMax:       3e9,
		LimitMinIntg:   0,
		LimitMaxIntg:   0,
	}
	corectl.tmp.setupStaticHeater(pidsetting)
}