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
	"path/filepath"

	"github.com/spf13/cobra"

	zlog "github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig string
var runtime string

func init() {
	deployCmd.Flags().StringVarP(&deployCmdFlags.flow, "flow", "f", "", "flow to run")
	deployCmd.Flags().StringVarP(&deployCmdFlags.kubeconfig, "config", "c", "", "k8s configuration")
	deployCmd.Flags().StringVarP(&deployCmdFlags.runtime, "runtime", "r", "", "the runtime to use")

	rootCmd.AddCommand(deployCmd)
}

// ==========================
//
// Run
//
// ==========================

type deployCmdFlagsType struct {
	kubeconfig string
	runtime    string
	flow       string
}

var deployCmdFlags deployCmdFlagsType

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy",
	Long:  `deploy`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := kubeconfig

		if conf == "" {
			conf = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		}

		config, err := clientcmd.BuildConfigFromFlags("", conf)
		if err != nil {
			zlog.Fatal().Msg(err.Error())
		}

		client, err := kubernetes.NewForConfig(config)
		if err != nil {
			zlog.Fatal().Msg(err.Error())
		}

		if client != nil {
		}
	},
}

func crdClient() {
}
