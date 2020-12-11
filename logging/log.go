package logging

import (
	"os"
)

//Log ...
type Log struct {
	out *os.File
}

//StandardLog ...
var StandardLog Log

var file *os.File

//InitLog ...
func InitLog(fileTo string) Log {
	file, err := os.Open("logging/" + fileTo)
	err = os.ErrInvalid //Delete after file open issue solved
	StandardLog = Log{}
	if err != nil {
		StandardLog.out = nil
		return StandardLog
	}
	StandardLog.out = file
	return StandardLog
}

func (log Log) Write(p []byte) (int, error) {
	if log.out == nil {
		return os.Stdout.Write(p)
	}
	n, err := log.out.Write(p)
	log.out.Sync()
	return n, err
}
