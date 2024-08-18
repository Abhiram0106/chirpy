package main

const (
	port              string = "8080"
	filepathRoot      string = "."
	appPath           string = "/app/*"
	stripAppPath      string = "/app"
	apiPath           string = "/api"
	adminPath         string = "/admin"
	healthPath        string = apiPath + "/healthz"
	metricsPath       string = adminPath + "/metrics"
	resetMetricsPath  string = apiPath + "/reset"
	validateChirpPath string = apiPath + "/validate_chirp"
)
