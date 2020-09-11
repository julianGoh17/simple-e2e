package models

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConvertingToInteger(t *testing.T) {
	tables := []struct {
		variable string
		expected int
		err      error
	}{
		{
			"No work",
			0,
			fmt.Errorf("Could not convert 'No work' to type 'int'"),
		},
		{
			"12",
			12,
			nil,
		},
		{
			"12.1",
			0,
			fmt.Errorf("Could not convert '12.1' to type 'int'"),
		},
	}

	converter := TypeConverter{}

	for _, table := range tables {
		val, err := converter.GetInteger(table.variable)
		if table.err == nil {
			assert.NoError(t, err)
			assert.Equal(t, val, table.expected)
		} else {
			assert.Error(t, err)
			assert.Equal(t, table.err.Error(), err.Error())
		}
	}
}

func TestConvertingToFloat32(t *testing.T) {
	tables := []struct {
		variable string
		expected float32
		err      error
	}{
		{
			"No work",
			0,
			fmt.Errorf("Could not convert 'No work' to type 'float32'"),
		},
		{
			"12",
			float32(12),
			nil,
		},
		{
			"12.1",
			float32(12.1),
			nil,
		},
	}

	converter := TypeConverter{}

	for _, table := range tables {
		val, err := converter.GetFloat32(table.variable)
		if table.err == nil {
			assert.NoError(t, err)
			assert.Equal(t, val, table.expected)
		} else {
			assert.Error(t, err)
			assert.Equal(t, table.err.Error(), err.Error())
		}
	}
}

func TestConvertingToFloat64(t *testing.T) {
	tables := []struct {
		variable string
		expected float64
		err      error
	}{
		{
			"No work",
			0,
			fmt.Errorf("Could not convert 'No work' to type 'float64'"),
		},
		{
			"12",
			float64(12),
			nil,
		},
		{
			"12.1",
			float64(12.1),
			nil,
		},
	}

	converter := TypeConverter{}

	for _, table := range tables {
		val, err := converter.GetFloat64(table.variable)
		if table.err == nil {
			assert.NoError(t, err)
			assert.Equal(t, val, table.expected)
		} else {
			assert.Error(t, err)
			assert.Equal(t, table.err.Error(), err.Error())
		}
	}
}

func TestConvertingToBoolean(t *testing.T) {
	tables := []struct {
		variable string
		expected bool
		err      error
	}{
		{
			"No work",
			false,
			fmt.Errorf("Could not convert 'No work' to type 'bool'"),
		},
		{
			"true",
			true,
			nil,
		},
		{
			"True",
			true,
			nil,
		},
	}

	converter := TypeConverter{}

	for _, table := range tables {
		val, err := converter.GetBoolean(table.variable)
		if table.err == nil {
			assert.NoError(t, err)
			assert.Equal(t, val, table.expected)
		} else {
			assert.Error(t, err)
			assert.Equal(t, table.err.Error(), err.Error())
		}
	}
}

func TestConvertingToTimeDuartion(t *testing.T) {
	tables := []struct {
		variable string
		expected time.Duration
		err      error
	}{
		{
			"No work",
			time.Duration(0),
			fmt.Errorf("Could not parse unit 'No work'. Must be one of: '%s', '%s', '%s'", millisecondsUnit, secondsUnit, minuteUnit),
		},
		{
			"1m",
			time.Duration(1) * time.Minute,
			nil,
		},
		{
			"12",
			time.Duration(0),
			fmt.Errorf("Could not parse unit ''. Must be one of: '%s', '%s', '%s'", millisecondsUnit, secondsUnit, minuteUnit),
		},
	}

	converter := TypeConverter{}

	for _, table := range tables {
		val, err := converter.GetTimeDuration(table.variable)
		if table.err == nil {
			assert.NoError(t, err)
			assert.Equal(t, val, table.expected)
		} else {
			assert.Error(t, err)
			assert.Equal(t, table.err.Error(), err.Error())
		}
	}
}

func TestConvertUnitToTestDuration(t *testing.T) {
	testCases := []struct {
		variable string
		time     time.Duration
	}{
		{
			"ms",
			time.Millisecond,
		},
		{
			"s",
			time.Second,
		},
		{
			"m",
			time.Minute,
		},
		{
			"won't work",
			time.Hour,
		},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.time, convertUnitToDuration(testCase.variable))
	}
}

func TestConvertingToIntegerArray(t *testing.T) {
	tables := []struct {
		variable string
		expected []int
		err      error
	}{
		{
			"No work,1",
			[]int{},
			fmt.Errorf("Could not convert 'No work,1' to type '[]int'"),
		},
		{
			"12,",
			[]int{},
			fmt.Errorf("Could not convert '12,' to type '[]int'"),
		},
		{
			"12,1,3",
			[]int{12, 1, 3},
			nil,
		},
	}

	converter := TypeConverter{}

	for _, table := range tables {
		val, err := converter.GetIntegerArray(table.variable)
		if table.err == nil {
			assert.NoError(t, err)
			assert.Equal(t, val, table.expected)
		} else {
			assert.Error(t, err)
			assert.Equal(t, table.err.Error(), err.Error())
		}
	}
}

func TestConvertingToFloat32Array(t *testing.T) {
	tables := []struct {
		variable string
		expected []float32
		err      error
	}{
		{
			"No work,1",
			[]float32{},
			fmt.Errorf("Could not convert 'No work,1' to type '[]float32'"),
		},
		{
			"12,",
			[]float32{},
			fmt.Errorf("Could not convert '12,' to type '[]float32'"),
		},
		{
			"12,1,3",
			[]float32{12, 1, 3},
			nil,
		},
	}

	converter := TypeConverter{}

	for _, table := range tables {
		val, err := converter.GetFloat32Array(table.variable)
		if table.err == nil {
			assert.NoError(t, err)
			assert.Equal(t, val, table.expected)
		} else {
			assert.Error(t, err)
			assert.Equal(t, table.err.Error(), err.Error())
		}
	}
}

func TestConvertingToFloat64Array(t *testing.T) {
	tables := []struct {
		variable string
		expected []float64
		err      error
	}{
		{
			"No work,1",
			[]float64{},
			fmt.Errorf("Could not convert 'No work,1' to type '[]float64'"),
		},
		{
			"12,",
			[]float64{},
			fmt.Errorf("Could not convert '12,' to type '[]float64'"),
		},
		{
			"12,1,3",
			[]float64{12, 1, 3},
			nil,
		},
	}

	converter := TypeConverter{}

	for _, table := range tables {
		val, err := converter.GetFloat64Array(table.variable)
		if table.err == nil {
			assert.NoError(t, err)
			assert.Equal(t, val, table.expected)
		} else {
			assert.Error(t, err)
			assert.Equal(t, table.err.Error(), err.Error())
		}
	}
}

func TestConvertingToBooleanArray(t *testing.T) {
	tables := []struct {
		variable string
		expected []bool
		err      error
	}{
		{
			"No work,1",
			[]bool{},
			fmt.Errorf("Could not convert 'No work,1' to type '[]bool'"),
		},
		{
			"12,",
			[]bool{},
			fmt.Errorf("Could not convert '12,' to type '[]bool'"),
		},
		{
			"True,False,true",
			[]bool{true, false, true},
			nil,
		},
	}

	converter := TypeConverter{}

	for _, table := range tables {
		val, err := converter.GetBooleanArray(table.variable)
		if table.err == nil {
			assert.NoError(t, err)
			assert.Equal(t, val, table.expected)
		} else {
			assert.Error(t, err)
			assert.Equal(t, table.err.Error(), err.Error())
		}
	}
}
