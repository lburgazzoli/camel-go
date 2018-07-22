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
	"github.com/spf13/cobra"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	zlog "github.com/rs/zerolog/log"
)

var kubeconfig string

func init() {
	deployCmd.Flags().StringVarP(&kubeconfig, "config", "c", "", "k8s configuration")

	rootCmd.AddCommand(runCmd)
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy",
	Long:  `deploy`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
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
