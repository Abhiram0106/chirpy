package main

const (
	port             string = "8080"
	filepathRoot     string = "."
	appPath          string = "/app/*"
	stripAppPath     string = "/app"
	apiPath          string = "/api"
	healthPath       string = apiPath + "/healthz"
	metricsPath      string = apiPath + "/metrics"
	resetMetricsPath string = apiPath + "/reset"
)
