package bdd

import (
	"flag"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "pretty",
	Paths:  []string{"features"},
}

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestMain(m *testing.M) {
	// Skip BDD tests in coverage mode or when explicitly requested
	if os.Getenv("COVERAGE_MODE") == "true" || os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		os.Exit(0)
	}

	flag.Parse()
	opts.Paths = flag.Args()

	status := godog.TestSuite{
		Name:                "Customer Service BDD Tests",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func TestCustomerService(t *testing.T) {
	// Skip BDD tests in coverage mode or when explicitly requested
	if os.Getenv("COVERAGE_MODE") == "true" || os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping BDD tests in coverage mode")
	}

	suite := godog.TestSuite{
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			InitializeScenario(ctx)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func TestBDDFeatures(t *testing.T) {
	// Skip BDD tests in coverage mode or when explicitly requested
	if os.Getenv("COVERAGE_MODE") == "true" || os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping BDD tests in coverage mode")
	}

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:        "pretty",
			Paths:         []string{"features"},
			TestingT:      t,
			StopOnFailure: true,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("BDD feature tests failed")
	}
}
