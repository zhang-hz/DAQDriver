package web

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {

	r := gin.New()

	apiv1 := r.Group("/api/v1")
	{
		apiv1.GET("/setfpga/:bitfilename", setPLBitFile)
		apiv1.GET("/connect", connectPC)
		apiv1.GET("/connectsocket", connectSocket)
		apiv1.GET("/disconnectsocket", disconnectSocket)
		apiv1.GET("/avgvoltage/:samplenumber", avgVoltage)
		apiv1.GET("/setadcvos/:adcch/:vosnumber", setADCVos)
		apiv1.GET("/setdacvoltage/:dacport/:voltage", setDACVoltage)
		apiv1.GET("/getdacvoltage/:dacport", getDACVoltage)
		apiv1.GET("/setdacoffset/:dacport/:offset", setDACOffset)
		apiv1.GET("/heater/static/start/:temperature/:basevoltage", startHeaterStaticPID)
		apiv1.GET("/heater/static/stop", stopHeaterStaticPID)
		apiv1.GET("/heater/pid/temperature/:temperature", setupTemperature)
		apiv1.GET("/heater/pid/parameters/:kp/:ki/:kd/:tolerance/:errorTolerance", setupHeaterPIDParameter)
		apiv1.GET("/heater/prog/start/:basevoltage/:heatingspeed/:basetemperature", startHeaterProgramPID)
		apiv1.GET("/heater/prog/stop", stopHeaterProgramPID)

	}

	return r

}
