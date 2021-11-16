package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lyswifter/ic-auth/types"
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

func buildOutLogs(input string) ([]byte, error) {
	inputs := strings.Split(input, "\n")

	if len(input) == 0 {
		return nil, errors.New("input length must not be zero")
	}

	var logs []string
	for _, v := range inputs {
		if v == "" {
			continue
		}

		////

		fat := fmt.Sprintf("[%s]	%s", time.Now().Format("2006-01-02 15:04:05.999"), v)
		logs = append(logs, fat)
	}

	var out = types.CmdOutput{
		TaskID: "xxx",
		Logs:   logs,
	}

	outbyte, err := json.Marshal(out)
	if err != nil {
		return nil, errors.New("marshal error")
	}

	return outbyte, nil
}
