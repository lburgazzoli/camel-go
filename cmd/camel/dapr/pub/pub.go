package pub

import (
	"strings"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/spf13/cobra"
)

func NewPubCmd() *cobra.Command {
	type opts struct {
		PubsubName string
		Topic      string
	}

	var o opts

	cmd := cobra.Command{
		Use:   "pub",
		Short: "pub",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := dapr.NewClient()
			if err != nil {
				panic(err)
			}
			defer client.Close()

			data := make(map[string]interface{})
			for i := range args {
				items := strings.Split(args[i], "=")
				if len(items) == 2 {
					data[items[0]] = items[1]
				}
			}

			// Publish a single event
			if err := client.PublishEvent(cmd.Context(), o.PubsubName, o.Topic, data); err != nil {
				panic(err)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&o.PubsubName, "pubsub-name", o.PubsubName, "pubsub-name")
	cmd.Flags().StringVar(&o.Topic, "topic", o.Topic, "topic")

	_ = cmd.MarkFlagRequired("pubsub-name")
	_ = cmd.MarkFlagRequired("topic")

	return &cmd
}
