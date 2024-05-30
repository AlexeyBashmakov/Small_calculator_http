package environ_vars

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	"calculator/internal/constants"
)

const time_default = "10000"

func CheckEnvironmentVariables() bool {
	env_vars := []string{
		constants.TimeAdd,
		constants.TimeSub,
		constants.TimeMult,
		constants.TimeDiv,
		constants.CompPow,
	}
	exist := true
	for _, k := range env_vars {
		_, exist = os.LookupEnv(k)
	}
	return exist
}

func SetEnvironmentVariables() error {
	env_vars := map[string]string{
		constants.TimeAdd:  time_default,
		constants.TimeSub:  time_default,
		constants.TimeMult: time_default,
		constants.TimeDiv:  time_default,
		constants.CompPow:  fmt.Sprintf("%d", runtime.NumCPU()/2),
	}
	var err error
	for k, v := range env_vars {
		err = os.Setenv(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetValue(name string) string {
	if (name == constants.TimeAdd) || (name == constants.TimeSub) ||
		(name == constants.TimeMult) || (name == constants.TimeDiv) ||
		(name == constants.CompPow) {
		val, exist := os.LookupEnv(name)
		if !exist {
			switch name {
			case constants.TimeAdd, constants.TimeSub, constants.TimeMult, constants.TimeDiv:
				err := os.Setenv(name, time_default)
				if err != nil {
					return ""
				}
				return time_default
			case constants.CompPow:
				err := os.Setenv(name, fmt.Sprintf("%d", runtime.NumCPU()/2))
				if err != nil {
					return ""
				}
				return fmt.Sprintf("%d", runtime.NumCPU()/2)
			}
		} else {
			return val
		}
	}

	return ""
}

func GetValueInt(name string) (int, bool) {
	r := GetValue(name)
	if r == "" {
		return 0, false
	} else {
		result, err := strconv.Atoi(r)
		if err != nil {
			return 0, false
		} else {
			return result, true
		}
	}
}
