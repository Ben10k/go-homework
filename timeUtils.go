package main

import "time"

func avg(array []time.Duration) time.Duration {
	if len(array) == 0 {
		return 0
	}
	result := time.Duration(0)
	for _, v := range array {
		result += v
	}
	return time.Duration(int64(result) / int64(len(array)))
}
func map2StringArray(array []time.Duration) []string {
	var output []string
	for _, v := range array {
		output = append(output, v.String())
	}
	return output
}
func timeFunctionDuration(function func() error) (time.Duration, error) {
	startTime := time.Now()
	err := function()
	if err != nil {
		return 0, err
	}
	return time.Now().Sub(startTime), nil
}
