//go:build integration

package cucumber

import (
	"flag"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/onsi/gomega"
	"github.com/spf13/pflag"
)

var opts = DefaultOptions()

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestMain(m *testing.M) {
	for _, arg := range os.Args[1:] {
		if arg == "-test.v=true" || arg == "-test.v" || arg == "-v" {
			opts.Format = "pretty"
		}
	}

	flag.Parse()
	pflag.Parse()

	if len(pflag.Args()) != 0 {
		opts.Paths = pflag.Args()
	}

	os.Exit(m.Run())
}

func TestFeatures(t *testing.T) {
	g := gomega.NewWithT(t)

	for i := range opts.Paths {
		root := opts.Paths[i]

		err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
			g.Expect(err).To(gomega.BeNil())

			if info.IsDir() {
				return nil
			}

			name := filepath.Base(info.Name())
			ext := filepath.Ext(info.Name())

			if ext != ".feature" {
				return nil
			}

			testName := strings.TrimSuffix(name, ext)
			testName = strings.ReplaceAll(testName, "-", "_")

			t.Run(testName, func(t *testing.T) {
				// To preserve the current behavior, the test are market to be "safely" run in parallel, however
				// we may think to introduce a new naming convention i.e. files that ends with _parallel would
				// cause t.Parallel() to be invoked, other tests won't, so they won't be executed concurrently.
				//
				// This could help reducing/removing the need of explicit lock
				t.Parallel()

				o := opts
				o.TestingT = t
				o.Paths = []string{path.Join(root, info.Name())}

				s := NewTestSuite()

				status := godog.TestSuite{
					Name:                "cos",
					Options:             &o,
					ScenarioInitializer: s.InitializeScenario,
				}.Run()

				g.Expect(status).To(gomega.Equal(0))
			})

			return nil
		})

		g.Expect(err).To(gomega.BeNil())
	}
}
