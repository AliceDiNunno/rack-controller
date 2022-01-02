package rest

import "os"

type Status struct {
	Success bool
	Message string
	Data    interface{}
	Host    string
}

func success(data interface{}) Status {
	hostname, err := os.Hostname()

	if err != nil {
		hostname = ""
	}

	return Status{
		Success: true,
		Message: "success",
		Data:    data,
		Host:    hostname,
	}
}
