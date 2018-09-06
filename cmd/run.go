// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	// import components
	"github.com/lburgazzoli/camel-go/app"
	"github.com/lburgazzoli/camel-go/logger"

	// load camel component
	_ "github.com/lburgazzoli/camel-go/components/http"
	_ "github.com/lburgazzoli/camel-go/components/log"
	_ "github.com/lburgazzoli/camel-go/components/timer"
)

func init() {
	runCmd.Flags().StringVarP(&runCmdFlags.flow, "flow", "f", "", "flow to run")
	runCmd.Flags().StringVarP(&runCmdFlags.route, "route", "r", "", "route to run")

	rootCmd.AddCommand(runCmd)
}

// ==========================
//
// Run
//
// ==========================

type runCmdFlagsType struct {
	flow  string
	route string
}

var runCmdFlags runCmdFlagsType

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run",
	Long:  `run`,
	Run: func(cmd *cobra.Command, args []string) {
		var ca *app.Application
		var err error

		if runCmdFlags.flow != "" {
			ca, err = app.New(runCmdFlags.flow)
		}
		if runCmdFlags.route != "" {
			b, err := ioutil.ReadFile(runCmdFlags.route)
			if err != nil {
				fmt.Print(err)
			}

			ca, err = app.NewJS(string(b))
		}

		if ca == nil {
			logger.Log(zerolog.FatalLevel, "Unable to build Application")
		}
		if err != nil {
			logger.Log(zerolog.FatalLevel, err.Error())
		}

		ca.Start()

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c

		ca.Stop()
	},
}
