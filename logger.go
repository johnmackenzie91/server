package server

// deadendLogger is the default logger, that yields no output
type deadendLogger struct{}

func (d deadendLogger) Error(...interface{}) {}
func (d deadendLogger) Info(...interface{})  {}
func (d deadendLogger) Debug(...interface{}) {}
