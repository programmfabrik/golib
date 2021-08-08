package lib

import "net/http"

type LogInfo map[string]interface{}

func LogInfoFromRequest(req *http.Request) LogInfo {
	logInfo := req.Context().Value("LogInfo")
	// this can be <nil> if the middleware is not loaded
	if logInfo == nil {
		return LogInfo{}
	}
	return logInfo.(LogInfo)
}
