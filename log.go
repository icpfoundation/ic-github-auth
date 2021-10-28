package main

import (
	"fmt"
	"time"
)

// infof infof
func Infof(format string, a ...interface{}) {
	fat := fmt.Sprintf("[%s]		%s\n", time.Now().Format("2006-01-02 15:04:05.999"), format)
	fmt.Printf(fat, a...)
}

// Errorf Errorf
func Errorf(format string, a ...interface{}) {
	fat := fmt.Sprintf("[%s]		%s\n", time.Now().Format("2006-01-02 15:04:05.999"), format)
	fmt.Printf(fat, a...)
}
