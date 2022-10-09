package main

type benchArgs struct {
	URL      string `long:"url" required:"true" description:"PostgreSQL connection URL"`
	Save     bool   `long:"save" description:"Set to save benchmark result"`
	SavePath string `long:"save-path" description:"Path to save benchmark result" default:"./bench"`
}
