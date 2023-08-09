package logger

import (
	"fmt"
	"log"
)

const (
	warningLevel = "[ WARNING ] "
	errorLevel   = "[ ERROR ] "
	infoLevel    = "[ INFO ] "
	fatalLevel   = "[ FATAL ERROR ] "
	successLevel = "[ SUCCESS EXECUTED ] "
)

func Warning(message string, args ...interface{}) {
	printLog(warningLevel, message, args...)
}

func Info(message string, args ...interface{}) {
	printLog(infoLevel, message, args...)
}

func Error(message string, args ...interface{}) {
	printLog(errorLevel, message, args...)
}

func Fatal(message string, args ...interface{}) {
	printLog(fatalLevel, message, args...)
}

func Success(message string, args ...interface{}) {
	printLog(successLevel, message, args...)
}

func printLog(level, message string, args ...interface{}) {
	log.Println(
		fmt.Sprintf(level+message, args...),
	)
}
