package web

import (
	"fmt"
	"os/exec"
)

func setBitfile(filename string) bool {

	command := exec.Command("python3", "./setuppl.py", filename)
	err := command.Run()

	return err == nil
}

type avgVoltageWeb struct {
	Voltage [4]float64
}

func getAvgVoltage(SampleNumber int64) avgVoltageWeb {

	av := device.AvgVoltage(SampleNumber)

	var result = &avgVoltageWeb{
		Voltage: av.Voltage,
	}

	fmt.Print("---------------CH0-----------------\n")
	fmt.Print(av.Voltage[0]/float64(1e6), " mV\n")
	fmt.Print(SampleNumber, "\n")
	fmt.Print("---------------CH1-----------------\n")
	fmt.Print(av.Voltage[1]/float64(1e6), " mV\n")
	fmt.Print(SampleNumber, "\n")
	fmt.Print("---------------CH2-----------------\n")
	fmt.Print(av.Voltage[2]/float64(1e6), " mV\n")
	fmt.Print(SampleNumber, "\n")
	fmt.Print("---------------CH3-----------------\n")
	fmt.Print(av.Voltage[3]/float64(1e6), " mV\n")
	fmt.Print(SampleNumber, "\n")
	fmt.Print("-------------Time------------------\n")
	fmt.Print(float64(av.Time/1e6), " ms\n")

	return *result

}

func setADCVosLocal(ADCCH int64, VosNumber float64) {
	device.SetADCVos(ADCCH, VosNumber)
}

func setDACVoltageLocal(DACport string, Voltage float64) float64 {
	return device.SetDACVoltage(DACport, Voltage)
}

func getDACVoltageLocal(DACport string) float64 {
	return device.GetDACVoltage(DACport)
}

func setDACOffsetLocal(DACport string, offsetVoltage float64) {
	device.SetDACOffset(DACport, offsetVoltage)
}

func startHeaterStaticPIDLocal(temperature float64, baseVoavgVoltage float64) {
	device.StartStaticHeater(baseVoavgVoltage, temperature)
}

func stopHeaterStaticPIDLocal() {
	device.StopStaticHeater()
}

func setupHeaterTemperaturePIDLocal(temperature float64) {
	device.SetupTemperature(temperature)
}

func setupHeaterPIDParameterLocal(kp float64, ki float64, kd float64, tolerance float64, errorTolerance float64) {
	device.SetupStaticHeater(kp, ki, kd, tolerance, errorTolerance)
}
