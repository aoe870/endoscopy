package logs

import (
	"endoscopy/internal/version"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/pkg/errors"
)

type logLevel int

const (
	levelDebug logLevel = iota
	levelInfo
	levelWarning
	levelError
)

var (
	prefixs = []string{
		"DEBUG",
		"INFO",
		"WARNING",
		"ERROR",
	}
	logFile *os.File
	logger  *log.Logger

	logFileSwitch = true
)

var logo = "                                                                        \n                                                                        \n  ,--,                                                                  \n,--.'|                                 ,---,                            \n|  | :                     ,---,     ,---.'|                      ,--,  \n:  : '                 ,-+-. /  |    |   | :                    ,'_ /|  \n|  ' |     ,--.--.    ,--.'|'   |    |   | |   ,--.--.     .--. |  | :  \n'  | |    /       \\  |   |  ,\"' |  ,--.__| |  /       \\  ,'_ /| :  . |  \n|  | :   .--.  .-. | |   | /  | | /   ,'   | .--.  .-. | |  ' | |  . .  \n'  : |__  \\__\\/: . . |   | |  | |.   '  /  |  \\__\\/: . . |  | ' |  | |  \n|  | '.'| ,\" .--.; | |   | |  |/ '   ; |:  |  ,\" .--.; | :  | : ;  ; |  \n;  :    ;/  /  ,.  | |   | |--'  |   | '/  ' /  /  ,.  | '  :  `--'   \\ \n|  ,   /;  :   .'   \\|   |/      |   :    :|;  :   .'   \\:  ,      .-./ \n ---`-' |  ,     .-./'---'        \\   \\  /  |  ,     .-./ `--`----'     \n         `--`---'                  `----'    `--`---'                   \n                                                                        "

func New(fileName string, have bool) {

	logFileSwitch = have
	if logFileSwitch {
		// 创建日志文件
		var err error
		filepath := path.Join(fileName)
		logFile, err = os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
		if err != nil {
			log.Println("log file create fail!")
		} else {
			// 创建日志
			logger = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
		}
	}
	Info(logo)
	Info("Engine start ..... Version: " + version.Version)
}

func out(level logLevel, v interface{}) {
	logger.SetPrefix(fmt.Sprintf("[%s] ", prefixs[level]))
	logger.Output(3, fmt.Sprint(v))
}

func Info(v interface{}) {
	if logFileSwitch {
		out(levelInfo, v)
	} else {
		fmt.Println("[Info] ", time.Now().Format(time.RFC3339), v)
	}

}

func Warn(v interface{}) {
	if logFileSwitch {
		out(levelWarning, v)
	} else {
		fmt.Println("[Warn] ", time.Now().Format(time.RFC3339), v)
	}

}

func Error(err error) {
	if logFileSwitch {
		out(levelError, errors.WithStack(err).Error())
	} else {
		fmt.Println("[Error] ", time.Now().Format(time.RFC3339), errors.WithStack(err).Error())
	}

}
