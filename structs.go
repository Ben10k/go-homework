package main

type Results struct {
	Measurements   []string `json:"measurements"`
	AverageLatency string   `json:"averageLatency"`
}

type MeasurementResult struct {
	Host     string  `json:"host"`
	Protocol string  `json:"protocol"`
	Results  Results `json:"results"`
}
