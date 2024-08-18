package main

const (
	port             string = "8080"
	filepathRoot     string = "."
	databasePath     string = filepathRoot + "/database.json"
	appPath          string = "/app/*"
	stripAppPath     string = "/app"
	apiPath          string = "/api"
	adminPath        string = "/admin"
	healthPath       string = apiPath + "/healthz"
	metricsPath      string = adminPath + "/metrics"
	resetMetricsPath string = apiPath + "/reset"
	chirpsPath       string = apiPath + "/chirps"
	chirpIDWildCard  string = "chirpID"
	chirpByIDPath    string = chirpsPath + "/{" + chirpIDWildCard + "}"
)
