package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	secondsUnit      = "s"
	millisecondsUnit = "ms"
	minuteUnit       = "m"
)

// TypeConverter aims to convert the string variable in the step.variables and converts it to the appropriate type wanted by the user
type TypeConverter struct {
}

// GetInteger converts the string'd variable and converts it to an integer if possible otherwise will return an error
func (converter *TypeConverter) GetInteger(variable string) (int, error) {
	value, err := strconv.Atoi(variable)
	if err != nil {
		return 0, fmt.Errorf("Could not convert '%s' to type 'int'", variable)
	}
	return value, nil
}

// GetFloat32 converts the string'd variable and converts it to an float32 if possible otherwise will return an error
func (converter *TypeConverter) GetFloat32(variable string) (float32, error) {
	value, err := strconv.ParseFloat(variable, 32)
	if err != nil {
		return float32(0), fmt.Errorf("Could not convert '%s' to type 'float32'", variable)
	}
	return float32(value), nil
}

// GetFloat64 converts the string'd variable and converts it to an float64 if possible otherwise will return an error
func (converter *TypeConverter) GetFloat64(variable string) (float64, error) {
	value, err := strconv.ParseFloat(variable, 64)
	if err != nil {
		return float64(0), fmt.Errorf("Could not convert '%s' to type 'float64'", variable)
	}
	return float64(value), nil
}

// GetBoolean converts the string'd variable and converts it to an boolean if possible otherwise will return an error
func (converter *TypeConverter) GetBoolean(variable string) (bool, error) {
	value, err := strconv.ParseBool(variable)
	if err != nil {
		return false, fmt.Errorf("Could not convert '%s' to type 'bool'", variable)
	}
	return value, nil
}

// GetTimeDuration converts the string'd variable and converts it to an time duration if possible otherwise will return an error
// Conversions:
// - s: seconds
// - ms: milliseconds
// - m: minutes
// I.E: 5m = 5 minutes
func (converter *TypeConverter) GetTimeDuration(variable string) (time.Duration, error) {
	stringNum := ""
	unit := ""
	for _, char := range variable {
		if unicode.IsDigit(char) {
			stringNum += string(char)
		} else {
			unit += string(char)
		}
	}

	duration := convertUnitToDuration(unit)
	if duration == time.Hour {
		return time.Duration(0), fmt.Errorf("Could not parse unit '%s'. Must be one of: '%s', '%s', '%s'", unit, millisecondsUnit, secondsUnit, minuteUnit)
	}

	num, err := strconv.Atoi(stringNum)
	if err != nil {
		return time.Duration(0), err
	}

	return time.Duration(num) * duration, nil
}

func convertUnitToDuration(unit string) time.Duration {
	switch unit {
	case secondsUnit:
		return time.Second
	case millisecondsUnit:
		return time.Millisecond
	case minuteUnit:
		return time.Minute
	}
	return time.Hour
}

// GetIntegerArray converts the string'd variable, splits it by ",", and converts it to an array of integers
// if possible otherwise will return an error
func (converter *TypeConverter) GetIntegerArray(variables string) ([]int, error) {
	components := strings.Split(variables, ",")
	intArray := []int{}
	for _, component := range components {
		parsedInt, err := converter.GetInteger(component)
		if err != nil {
			return []int{}, fmt.Errorf("Could not convert '%s' to type '[]int'", variables)
		}
		intArray = append(intArray, parsedInt)
	}
	return intArray, nil
}

// GetFloat32Array converts the string'd variable, splits it by ",", and converts it to an array of float32s
// if possible otherwise will return an error
func (converter *TypeConverter) GetFloat32Array(variables string) ([]float32, error) {
	components := strings.Split(variables, ",")
	intArray := []float32{}
	for _, component := range components {
		parsedInt, err := converter.GetFloat32(component)
		if err != nil {
			return []float32{}, fmt.Errorf("Could not convert '%s' to type '[]float32'", variables)
		}
		intArray = append(intArray, parsedInt)
	}
	return intArray, nil
}

// GetFloat64Array converts the string'd variable, splits it by ",", and converts it to an array of float64s
// if possible otherwise will return an error
func (converter *TypeConverter) GetFloat64Array(variables string) ([]float64, error) {
	components := strings.Split(variables, ",")
	intArray := []float64{}
	for _, component := range components {
		parsedInt, err := converter.GetFloat64(component)
		if err != nil {
			return []float64{}, fmt.Errorf("Could not convert '%s' to type '[]float64'", variables)
		}
		intArray = append(intArray, parsedInt)
	}
	return intArray, nil
}

// GetBooleanArray converts the string'd variable, splits it by ",", and converts it to an array of booleans
// if possible otherwise will return an error
func (converter *TypeConverter) GetBooleanArray(variables string) ([]bool, error) {
	components := strings.Split(variables, ",")
	intArray := []bool{}
	for _, component := range components {
		parsedInt, err := converter.GetBoolean(component)
		if err != nil {
			return []bool{}, fmt.Errorf("Could not convert '%s' to type '[]bool'", variables)
		}
		intArray = append(intArray, parsedInt)
	}
	return intArray, nil
}
