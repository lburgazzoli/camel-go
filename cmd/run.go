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
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	zlog "github.com/rs/zerolog/log"

	// import components
	"github.com/lburgazzoli/camel-go/app"

	// load camel component
	_ "github.com/lburgazzoli/camel-go/components/log"
	_ "github.com/lburgazzoli/camel-go/components/timer"
)

var flow string

func init() {
	runCmd.Flags().StringVarP(&flow, "flow", "f", "", "flow to run")

	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run",
	Long:  `run`,
	Run: func(cmd *cobra.Command, args []string) {
		app, err := app.New(flow)

		if err != nil {
			zlog.Fatal().Msgf("%s", err)
		}

		app.Start()

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c

		app.Stop()
	},
}
