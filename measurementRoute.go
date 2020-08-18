package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/eapache/go-resiliency.v1/deadline"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func measureHandler(w http.ResponseWriter, r *http.Request) {

	samples, err := strconv.Atoi(r.URL.Query().Get("samples"))
	if err != nil {
		w.WriteHeader(400)
		result, _ := json.Marshal(struct {
			Message string
		}{
			Message: "Invalid sample count",
		})

		w.Write(result)
		return
	}

	protocol := r.URL.Query().Get("protocol")
	host := r.URL.Query().Get("host")
	if protocol == "" || host == "" {
		w.WriteHeader(400)
		result, _ := json.Marshal(struct {
			Message  string
			Host     string
			Protocol string
		}{
			Message:  "Both [host] and [protocol] parameters must be present",
			Host:     host,
			Protocol: protocol,
		})

		w.Write(result)
		return
	}
	urlString := buildUrlString(protocol, host)

	work, err := Work(urlString, samples, 5000*time.Millisecond)
	if err != nil {
		w.WriteHeader(400)
		result, _ := json.Marshal(struct {
			Message string
		}{
			Message: err.Error(),
		})

		w.Write(result)
		return
	}
	result, err := json.Marshal(MeasurementResult{
		Host:     host,
		Protocol: protocol,
		Results:  work,
	})
	if err != nil {
		w.WriteHeader(400)
		result, _ := json.Marshal(struct {
			Message string
		}{
			Message: err.Error(),
		})

		w.Write(result)
		return
	}
	w.Write(result)
}

func Work(urlString string, samples int, timeout time.Duration) (Results, error) {
	var latencies []time.Duration
	dl := deadline.New(timeout)

	for i := 0; i < samples; i++ {
		err := dl.Run(func(stopper <-chan struct{}) error {
			duration, err := timeFunctionDuration(func() error {
				err := openPage(urlString)
				if err != nil {
					return err
				}
				return nil
			})

			if err != nil {
				return err
			}
			latencies = append(latencies, duration)

			return nil
		})

		if err != nil {
			return Results{}, err
		}
	}

	return Results{
		Measurements:   map2StringArray(latencies),
		AverageLatency: avg(latencies).String(),
	}, nil
}

func openPage(url string) error {
	res, err := http.Get(url)
	if err != nil || res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("unable to open url [%s]", url))
	}
	return nil
}

func buildUrlString(scheme string, host string) string {
	u := url.URL{
		Scheme: scheme,
		Host:   host,
	}
	return u.String()
}
