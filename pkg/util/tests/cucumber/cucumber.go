package cucumber

import (
	"os"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

const (
	DefaultConcurrency = 10
)

func NewTestSuite() *TestSuite {
	return &TestSuite{}
}

func DefaultOptions() godog.Options {
	opts := godog.Options{
		Output:      colors.Colored(os.Stdout),
		Format:      "progress",
		Paths:       []string{"features"},
		Randomize:   time.Now().UTC().UnixNano(), // randomize TestScenario execution order
		Concurrency: DefaultConcurrency,
	}

	return opts
}

// TestSuite holds the state global to all the test scenarios.
// It is accessed concurrently from all test scenarios.
type TestSuite struct {
}

// TestScenario holds that state of single scenario.
// It is not accessed concurrently.
type TestScenario struct {
	Suite     *TestSuite
	Variables map[string]interface{}
}

// StepModules is the list of functions used to add steps to a godog.ScenarioContext, you can
// add more to this list if you need test TestSuite specific steps.
var StepModules []func(ctx *godog.ScenarioContext, s *TestScenario)

func (suite *TestSuite) InitializeScenario(ctx *godog.ScenarioContext) {
	s := &TestScenario{
		Suite:     suite,
		Variables: map[string]interface{}{},
	}

	for _, module := range StepModules {
		module(ctx, s)
	}
}
