package driver

//通道常量设置
const dataChDepth = int64(100)
const socketChDepth = int64(10)
const helperChDepth = int64(100000)
const socketDataLength = int64(500)
const socketDownSampleRate = int64(20)
const heaterDownSampleRate = int64(10)

//AXI地址映射
const ADCFilterAddress int64 = 0xA0000000
const ADCSPIAddress int64 = 0xB0020000
const ADCRawAddress int64 = 0xA0010000
const DACSPIAddress int64 = 0xB0010000

type DAQInfo struct {
	vos [4]float64
}

type HeaterInfo struct {
	voltage     [2]float64
	power       [2]float64
	temperature [2]float64
}

type CoreControllerInterface interface {
	Initialize()
	ConnectSocket() bool
	DisconnectSocket()
	ConnectADC()
	StartFetchData()
	StopFetchData()
	AvgVoltage(int64) avgVoltageData
	SetADCVos(int64, float64)
	SetDACVoltage(string, float64) float64
	GetDACVoltage(string) float64
	SetDACOffset(string, float64)
	StartProgramHeater(float64, float64, float64)
	StopProgramHeater()
	StartStaticHeater(float64, float64)
	StopStaticHeater()
	SetupStaticHeater(float64, float64, float64, float64, float64)
	SetupTemperature(float64)
}

type CoreController struct {
	WORKFLAG   int32
	datach     chan DAQDataCH
	socketch   chan socketCH
	dchhelper1 chan DAQDataCH
	dchhelper2 chan DAQDataCH
	ctlch      chan string
	DAQSetting DAQInfo

	socket socketController
	adc    ADCController
	dac    DACController
	tmp    TemperatureController

	heater HeaterInfo
}

var runningDAQ = bool(true)
var socketCHSign = uint8(0)
var helperCHSign = uint8(0)

func (corectl *CoreController) Initialize() {

	//初始化全局标志位

	runningDAQ = true
	socketCHSign = 0
	helperCHSign = 0

	//初始化缓存

	corectl.heater = HeaterInfo{[2]float64{0, 0}, [2]float64{0, 0}, [2]float64{0, 0}}

	//初始化信号通道

	corectl.datach = make(chan DAQDataCH, dataChDepth)
	corectl.socketch = make(chan socketCH, socketChDepth)
	corectl.dchhelper1 = make(chan DAQDataCH, helperChDepth)
	corectl.dchhelper2 = make(chan DAQDataCH, helperChDepth)
	corectl.ctlch = make(chan string, 10)

	//初始化Socket网络设备

	corectl.socket = *newSocketController()
	corectl.socket.setCH(corectl.socketch)

	//初始化ADC设置

	ADCvos := [4]float64{2300.9, 7023, 0, 0}
	corectl.DAQSetting = DAQInfo{vos: ADCvos}

	//初始化ADC

	corectl.adc = *newADCController(ADCSPIAddress, ADCFilterAddress)
	corectl.adc.setVos(0, corectl.DAQSetting.vos[0])
	corectl.adc.setVos(1, corectl.DAQSetting.vos[1])

	//初始化DAC
	corectl.dac = *newDACController(DACSPIAddress)
	corectl.dac.regDACPort("TP1", "HVDAC", 0)
	corectl.dac.regDACPort("TP2", "HVDAC", 1)

	//初始化温度控制器
	corectl.tmp = *newTemperatureController(corectl.dac, corectl.dchhelper2, &corectl.heater)

	//初始化数据交换器

	go interconnectHub(corectl.datach, &corectl.heater, corectl.socketch, corectl.dchhelper1, corectl.dchhelper2)
}

func (corectl *CoreController) ConnectSocket() bool {
	resultChan := make(chan bool, 1)
	go corectl.socket.start(resultChan)
	result := <-resultChan
	close(resultChan)
	return result
}

func (corectl *CoreController) DisconnectSocket() {
	corectl.socket.stop()
}

func (corectl *CoreController) ConnectADC() {
	corectl.adc.initialize()
}

func (corectl *CoreController) StartFetchData() {
	runningDAQ = true
	go corectl.fetchData(corectl.datach)
}

func (corectl *CoreController) StopFetchData() {
	runningDAQ = false
}
