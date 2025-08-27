package otelwrap

import (
	"fmt"
	"os"
)

const (
	otelTracesExportersEnvKey = "OTEL_TRACES_EXPORTER"
)

func init() {
	checkExporterEnv()
}

func checkExporterEnv() {
	expType := os.Getenv(otelTracesExportersEnvKey)
	if expType != "" {
		return
	}

	err := os.Setenv(otelTracesExportersEnvKey, "none")
	if err != nil {
		panic(fmt.Errorf("os.Setenv[OTEL_TRACES_EXPORTER=none] err[%s]", err))
	}
}
