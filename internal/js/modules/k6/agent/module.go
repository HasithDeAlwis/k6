package agent

import (
	"errors"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/grafana/sobek"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/js/promises"
)

type (
	// RootModule is the global module instance that will create instances of our
	// module for each VU.
	RootModule struct{}

	// ModuleInstance represents an instance of the agent module for a single VU.
	ModuleInstance struct {
		vu    modules.VU
		agent *Agent
	}
)

var (
	_ modules.Module   = &RootModule{}
	_ modules.Instance = &ModuleInstance{}
)

// New returns a pointer to a new [RootModule] instance.
func New() *RootModule {
	return &RootModule{}
}

// NewModuleInstance implements the modules.Module interface and returns a new
// instance of our module for the given VU.
func (rm *RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	_, apiKeyDefined := vu.InitEnv().LookupEnv("ANTHROPIC_API_KEY")
	if !apiKeyDefined {
		common.Throw(vu.Runtime(), errors.New("ANTHROPIC_API_KEY must be provided to use the k6/agent"))
	}

	client := anthropic.NewClient()
	agent := NewAgent(&client)
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
	return "Hello, I'm an agent that will explore: " + url, nil
}
