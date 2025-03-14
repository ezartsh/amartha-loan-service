package logger

var AppLog LogHandler

func InitAppLog() {
	AppLog = newAppLogger("app.log")
}

func newAppLogger(fileName string) LogHandler {
	return NewFileConsoleStructured(fileName)
}
