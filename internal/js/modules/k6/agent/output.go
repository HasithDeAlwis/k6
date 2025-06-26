package agent

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"go.k6.io/k6/metrics"
	"go.k6.io/k6/output"
)

// NewOutput is a wrapper on top of New, that uses the given output.Params
// and returns (the same) output.Output instance.
func NewOutput(params output.Params) (output.Output, error) {
	root := New()
	root.params = params
	root.logger = params.Logger
	return root, nil
}

// Description implements the output.Output interface, by returning the module's description.
func (*RootModule) Description() string {
	return "xk6-custosummary"
}

// Start implements the output.Output interface, exposing a method to initialize the output.
func (rm *RootModule) Start() error {
	rm.logger.Debug("Starting output...")

	pf, err := output.NewPeriodicFlusher(time.Second, rm.flushMetrics)
	if err != nil {
		return err
	}

	rm.logger.Debug("Started!")
	rm.start = time.Now()
	rm.periodicFlusher = pf

	return nil
}

// StopWithTestError flushes all remaining metrics and finalizes the test run
func (rm *RootModule) StopWithTestError(err error) error {
	logger := rm.loggerWithError(err)
	logger.Debug("Stopping...")
	defer rm.logger.Debug("Stopped!")

	rm.periodicFlusher.Stop()

	_, _ = fmt.Fprintln(os.Stdout)
	_, _ = fmt.Fprintln(os.Stdout)
	_, _ = fmt.Fprintln(os.Stdout)

	for url, vitals := range rm.metrics {
		_, _ = fmt.Fprintln(os.Stdout, url)
		for wvName, wvSink := range vitals {
			_, _ = fmt.Fprintln(os.Stdout, wvName, wvSink.Format(time.Since(rm.start)))
		}
	}

	return nil
}

// Stop implements the output.Output interface, exposing a method to stop the output.
func (rm *RootModule) Stop() error {
	return rm.StopWithTestError(nil)
}

func (rm *RootModule) loggerWithError(err error) logrus.FieldLogger {
	logger := rm.logger
	if err != nil {
		logger = logger.WithError(err)
	}
	return logger
}

func (rm *RootModule) flushMetrics() {
	samples := rm.GetBufferedSamples()
	for _, sc := range samples {
		samples := sc.GetSamples()
		for _, sample := range samples {
			rm.flushSample(sample)
		}
	}
}

func (rm *RootModule) flushSample(s metrics.Sample) {
	// Name is the tag used for the site URL
	nameTag, hasName := s.Tags.Get("name")
	if !hasName {
		return
	}

	// We register the metric and its sub-metrics,
	// and we add the sample value to their sinks.
	rm.Metrics.addSmple(nameTag, s)
}
