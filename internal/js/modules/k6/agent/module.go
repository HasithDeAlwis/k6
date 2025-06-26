package agent

import (
	"errors"
	"fmt"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/sirupsen/logrus"
	"go.k6.io/k6/output"
	"time"

	"github.com/grafana/sobek"

	"go.k6.io/k6/internal/js/modules/k6/browser/browser"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/js/modules/k6/http"
	"go.k6.io/k6/js/promises"
)

// New returns a pointer to a new RootModule instance.
var New func() *RootModule

func init() {
	// Initialize the global RootModule instance accessor.
	root := &RootModule{
		Metrics: NewMetrics(),
	}

	New = func() *RootModule { return root }
}

type (
	// RootModule is the global module instance that will create instances of our
	// module for each VU.
	RootModule struct {
		// Metrics
		*Metrics
		// Output stuff
		output.SampleBuffer
		start           time.Time
		params          output.Params
		periodicFlusher *output.PeriodicFlusher
		logger          logrus.FieldLogger
	}

	// ModuleInstance represents an instance of the agent module for a single VU.
	ModuleInstance struct {
		vu    modules.VU
		agent *Agent
	}
)

var (
	_ interface {
		modules.Module
		output.WithStopWithTestError
	} = &RootModule{}
	_ modules.Instance = &ModuleInstance{}
)

// NewModuleInstance implements the modules.Module interface and returns a new
// instance of our module for the given VU.
func (rm *RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	// Check whether the Anthropic API Key is provided
	_, apiKeyDefined := vu.InitEnv().LookupEnv("ANTHROPIC_API_KEY")
	if !apiKeyDefined {
		common.Throw(vu.Runtime(), errors.New("ANTHROPIC_API_KEY must be provided to use the k6/agent"))
	}

	// Initialize the Anthropic client
	client := anthropic.NewClient()
	agent := NewAgent(&client)

	// Set up the HTTP and Browser modules
	if err := vu.Runtime().Set("this", vu.Runtime().GlobalObject()); err != nil {
		common.Throw(vu.Runtime(), err)
	}

	if err := vu.Runtime().Set("http", http.New().NewModuleInstance(vu).Exports().Default); err != nil {
		common.Throw(vu.Runtime(), err)
	}

	if err := vu.Runtime().Set("browser", browser.New().NewModuleInstance(vu).Exports().Default.(*browser.JSModule).Browser); err != nil {
		common.Throw(vu.Runtime(), err)
	}

	return &ModuleInstance{vu: vu, agent: agent}
}

// Exports implements the modules.Module interface and returns the exports of
// our module.
func (mi *ModuleInstance) Exports() modules.Exports {
	return modules.Exports{
		Named: map[string]any{
			"explore": mi.Explore,
		},
	}
}

func (mi *ModuleInstance) Explore(url sobek.Value) *sobek.Promise {
	promise, resolve, reject := promises.New(mi.vu)

	// TODO: Handle the case when it's called during the Init context.

	if common.IsNullish(url) {
		reject(mi.vu.Runtime().NewTypeError("explore() failed; reason: url cannot be null or undefined"))
		return promise
	}

	// Obtain the underlying path string from the JS value.
	pathStr := url.String()
	if pathStr == "" {
		reject(mi.vu.Runtime().NewTypeError("explore() failed; reason: url cannot be empty"))
		return promise
	}

	go func() {
		result, err := mi.exploreImpl(pathStr)
		if err != nil {
			reject(err)
			return
		}

		resolve(result)
	}()

	return promise
}

func (mi *ModuleInstance) exploreImpl(url string) (string, error) {
	// This should be refactored, it's here just for the sake of the example.
	if _, err := mi.vu.Runtime().RunString(`(async function() { this.page = await browser.newPage(); }());`); err != nil {
		common.Throw(mi.vu.Runtime(), fmt.Errorf("could not call browser.newPage(): %w", err))
	}
	time.Sleep(500 * time.Millisecond)

	// From here, we could start our interaction with the LLM.
	script := `(async function() { return await this.page.goto("%s"); }());`

	// Navigate to the given url
	_, err := mi.vu.Runtime().RunString(fmt.Sprintf(script, url))
	if err != nil {
		return err.Error(), err
	}

	// Navigate to QuickPizza, to generate other metrics
	_, err = mi.vu.Runtime().RunString(fmt.Sprintf(script, "https://quickpizza.grafana.com/"))
	if err != nil {
		return err.Error(), err
	}

	return "Hello, I'm an agent that will explore: " + url, nil
}
