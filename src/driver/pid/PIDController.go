package pid

import (
	"math"
)

type PIDController struct {
	Kp               float64
	Ki               float64
	Kd               float64
	Interval         float64
	tolerance        float64
	errorTolerance   float64
	tau              float64
	Target           float64
	errLinr          float64
	errDerv          float64
	errIntg          float64
	lastError        float64
	lastMeasurement  float64
	lastOutPut       float64
	bufferCompension float64
	limitMin         float64
	limitMax         float64
	limitMinIntg     float64
	limitMaxIntg     float64
}

type PIDSetting struct {
	Kp             float64
	Ki             float64
	Kd             float64
	Interval       float64
	Tau            float64
	Tolerance      float64
	ErrorTolerance float64
	LimitMin       float64
	LimitMax       float64
	LimitMinIntg   float64
	LimitMaxIntg   float64
}

func (pidctl *PIDController) initialize() {

	pidctl.Kp = 0
	pidctl.Ki = 0
	pidctl.Kd = 0
	pidctl.Interval = 20000
	pidctl.tau = 0
	pidctl.tolerance = 0
	pidctl.errorTolerance = 0
	pidctl.Target = 0
	pidctl.errLinr = 0
	pidctl.errDerv = 0
	pidctl.errIntg = 0
	pidctl.lastError = 0
	pidctl.lastOutPut = 0
	pidctl.bufferCompension = 0
	pidctl.lastMeasurement = 0
	pidctl.limitMin = 0
	pidctl.limitMax = 0
	pidctl.limitMinIntg = 0
	pidctl.limitMaxIntg = 0

}

func (pidctl *PIDController) Linear(measurement float64) float64 {

	//fmt.Println(pidctl.Target)

	errorNow := pidctl.Target - measurement

	if errorNow < pidctl.errorTolerance && errorNow > -pidctl.errorTolerance {
		errorNow = 0
	}

	compension := pidctl.Kp * errorNow

	pidctl.bufferCompension = pidctl.bufferCompension + compension
	output := pidctl.lastOutPut + pidctl.bufferCompension

	if output > pidctl.limitMax {
		output = pidctl.limitMax
	} else if output < pidctl.limitMin {
		output = pidctl.limitMin
	}

	if math.Abs(pidctl.bufferCompension) < pidctl.tolerance {
		return 0
	} else {
		pidctl.lastOutPut = output
		pidctl.bufferCompension = 0
		return output
	}

}

func (pidctl *PIDController) Simple(measurement float64) float64 {

	errorNow := pidctl.Target - measurement

	pidctl.errLinr = pidctl.Kp * errorNow
	pidctl.errIntg = pidctl.errIntg + 0.5*pidctl.Ki*pidctl.Interval*(errorNow+pidctl.lastError)
	pidctl.errDerv = -1 * (2*pidctl.Kd*(measurement-pidctl.lastMeasurement) + 2*(pidctl.tau-pidctl.Interval)*pidctl.errDerv) / (2 * (pidctl.tau + pidctl.Interval))

	if pidctl.errIntg > pidctl.limitMaxIntg {
		pidctl.errIntg = pidctl.limitMaxIntg
	} else if pidctl.errIntg < pidctl.limitMinIntg {
		pidctl.errIntg = pidctl.limitMinIntg
	}

	pidctl.lastError = errorNow

	compension := pidctl.errLinr + pidctl.errDerv + pidctl.errIntg

	if compension > pidctl.limitMax {
		compension = pidctl.limitMax
	} else if compension < pidctl.limitMin {
		compension = pidctl.limitMin
	}

	if math.Abs(compension-pidctl.lastOutPut) < pidctl.tolerance {
		return 0
	} else {
		pidctl.lastOutPut = compension
		return compension
	}

}

func (pidctl *PIDController) Reset() {

	pidctl.Target = 0
	pidctl.errLinr = 0
	pidctl.errDerv = 0
	pidctl.errIntg = 0
	pidctl.lastError = 0
	pidctl.lastOutPut = 0
	pidctl.lastMeasurement = 0

}

func (pidctl *PIDController) Setup(setting PIDSetting) {

	pidctl.Kp = setting.Kp
	pidctl.Ki = setting.Ki
	pidctl.Kd = setting.Kd
	pidctl.Interval = setting.Interval
	pidctl.tau = setting.Tau
	pidctl.tolerance = setting.Tolerance
	pidctl.errorTolerance = setting.ErrorTolerance
	pidctl.limitMin = setting.LimitMin
	pidctl.limitMax = setting.LimitMax
	pidctl.limitMinIntg = setting.LimitMinIntg
	pidctl.limitMaxIntg = setting.LimitMaxIntg

}

func NewPIDController() *PIDController {

	var PIDCTL = &PIDController{}
	PIDCTL.initialize()

	return PIDCTL
}
