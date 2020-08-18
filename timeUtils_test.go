package main

import (
	"errors"
	"testing"
	"time"
)

func TestAvg_Zero(t *testing.T) {
	input := []time.Duration{time.Duration(0), time.Duration(0), time.Duration(0), time.Duration(0)}
	output := avg(input)
	if output != time.Duration(0) {
		t.Error("Average of multiple zeros should still be zero")
	}
}

func TestAvg_Constant(t *testing.T) {
	input := []time.Duration{time.Hour, time.Hour, time.Hour, time.Hour}
	output := avg(input)
	if output != time.Hour {
		t.Error("Average of 4 hour-long durations should be an hour")
	}
}

func TestAvg_Multiple(t *testing.T) {
	input := []time.Duration{2 * time.Hour, 2 * time.Hour, 4 * time.Hour, 4 * time.Hour}
	output := avg(input)
	if output != 3*time.Hour {
		t.Error("Average of 2 2-hour-long and 2 4-hour-long durations should be 3 hours")
	}
}

func TestTimeFunctionDuration_verifyError(t *testing.T) {

	_, err := timeFunctionDuration(func() error {
		return errors.New("fail")
	})
	if err != nil && err.Error() == "fail" {
		return
	}
	t.Error("timing function should return inner error")

}
func TestTimeFunctionDuration_verifyFunctionCall(t *testing.T) {
	hasFunctionBeenCalled := false

	_, _ = timeFunctionDuration(func() error {
		hasFunctionBeenCalled = true
		return nil
	})
	if !hasFunctionBeenCalled {
		t.Error("Average of multiple zeros should still be zero")
	}
}

func TestMap2StringArray_nil(t *testing.T) {
	results := map2StringArray(nil)
	if results != nil {
		t.Error("map2StringArray with nil input should return [nil]")
	}
}

func TestMap2StringArray_singleZeroDuration(t *testing.T) {
	results := map2StringArray([]time.Duration{time.Duration(0)})
	if len(results) != 1 {
		t.Error("map2StringArray should return 1 result with 1 input")
	}
	if results[0] != "0s" {
		t.Error("map2StringArray should return [0s] when given zero duration")
	}
}

func TestMap2StringArray_multipleNonZeroDurations(t *testing.T) {
	results := map2StringArray([]time.Duration{time.Hour, time.Minute})
	if len(results) != 2 {
		t.Error("map2StringArray should return 2 results with 2 inputs")
	}
	if results[0] != "1h0m0s" {
		t.Errorf("map2StringArray  returned [%s] instead of  [1h0m0s] when given [time.Hour] duration", results[0])
	}
	if results[1] != "1m0s" {
		t.Errorf("map2StringArray  returned [%s] instead of  [1m0s] when given [time.Minute] duration", results[1])
	}
}
