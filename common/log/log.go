// log
// use like log.Debug("debug %s", Password("secret")
// 		log.Info("info")
//    	log.Notice("notice")
//	    log.Warning("warning")
//		log.Error("err")
//		log.Critical("crit")
package log

import (
	"github.com/op/go-logging"
	"io"
	"os"
)

var (
	logger = logging.MustGetLogger("")
)

type LogConfig struct {
	Module   string `yaml:"module"`
	FileName string `yaml:"filename"`
	Level    int    `yaml:"level"`
	Format   string `yaml:"format"`
}

func InitLog(logConfig *LogConfig) error {
	logger = logging.MustGetLogger(logConfig.Module)
	var output io.Writer
	if len(logConfig.FileName) == 0 {
		output = os.Stderr
	} else {
		f, err := os.Create(logConfig.FileName)
		if err != nil {
			return err
		}
		output = f
	}
	backend := logging.NewLogBackend(output, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend,
		logging.MustStringFormatter(logConfig.Format))

	backendLeveled := logging.AddModuleLevel(backend)
	backendLeveled.SetLevel(logging.Level(logConfig.Level), "")
	logging.SetBackend(backendLeveled, backendFormatter)
	logger.Debug("init log end")
	return nil
}

func Debug(format string, args ...interface{}) {
	logger.Debug(format, args)
}

func Info(format string, args ...interface{}) {
	logger.Info(format, args)
}

func Error(format string, args ...interface{}) {
	logger.Error(format, args)
}
func Println(format string, args ...interface{}) {
	logger.Info(format, args)
}
func Fatal(format string, args ...interface{}) {
	logger.Critical(format, args)
}
