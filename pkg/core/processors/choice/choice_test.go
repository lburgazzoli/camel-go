package choice

import (
	"context"
	"testing"

	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"
)

func TestChoice(t *testing.T) {
	support.Run(t, "simple", func(t *testing.T, ctx context.Context) {
		t.Helper()

		/*
			c := camel.GetContext(ctx)

			content1 := `{ "foo": "bar" }`
			wg1 := make(chan camel.Message)

			content2 := `{ "bar": "baz" }`
			wg2 := make(chan camel.Message)

			choice := Choice{
				DefaultVerticle: processors.NewDefaultVerticle(),
				When: []*When{
					{
						DefaultVerticle: processors.NewDefaultVerticle(),
						Language: language.Language{
							Jq: &jq.Jq{Expression: `.foo == "bar"`},
						},
						Steps: []processors.Step{
							processors.NewStep(support.NewChannelVerticle(wg1)),
						},
					},
					{
						DefaultVerticle: processors.NewDefaultVerticle(),
						Language: language.Language{
							Jq: &jq.Jq{Expression: `.bar == "baz"`},
						},
						Steps: []processors.Step{
							processors.NewStep(support.NewChannelVerticle(wg2)),
						},
					},
				},
			}

			id, err := choice.Reify(ctx)

			assert.Nil(t, err)
			assert.NotNil(t, id)

			msg1, err := message.New()
			msg1.SetContent(content1)
			assert.Nil(t, err)

			assert.Nil(t, c.Send(id, msg1))

			msg2, err := message.New()
			msg2.SetContent(content2)
			assert.Nil(t, err)

			assert.Nil(t, c.Send(id, msg2))

			RegisterTestingT(t)

			recv1, err := message.New()
			assert.Nil(t, err)

			recv2, err := message.New()
			assert.Nil(t, err)

			Eventually(wg1).Should(Receive(&recv1))
			Eventually(wg2).Should(Receive(&recv2))

			Expect(recv1.Content()).To(Equal(content1))
			Expect(recv2.Content()).To(Equal(content2))
		*/
	})

	support.Run(t, "otherwise", func(t *testing.T, ctx context.Context) {
		t.Helper()

		/*
			c := camel.GetContext(ctx)

			content1 := `{ "foo": "bar" }`
			wg1 := make(chan camel.Message)

			content2 := `{ "bar": "baz" }`
			wg2 := make(chan camel.Message)

			content3 := `{ "foo": "baz" }`
			wg3 := make(chan camel.Message)

			choice := Choice{
				DefaultVerticle: processors.NewDefaultVerticle(),
				When: []*When{
					{
						DefaultVerticle: processors.NewDefaultVerticle(),
						Language: language.Language{
							Jq: &jq.Jq{Expression: `.foo == "bar"`},
						},
						Steps: []processors.Step{
							processors.NewStep(support.NewChannelVerticle(wg1)),
						},
					},
					{
						DefaultVerticle: processors.NewDefaultVerticle(),
						Language: language.Language{
							Jq: &jq.Jq{Expression: `.bar == "baz"`},
						},
						Steps: []processors.Step{
							processors.NewStep(support.NewChannelVerticle(wg2)),
						},
					},
				},
				Otherwise: &Otherwise{
					DefaultVerticle: processors.NewDefaultVerticle(),
					Steps: []processors.Step{
						processors.NewStep(support.NewChannelVerticle(wg3)),
					},
				},
			}

			id, err := choice.Reify(ctx)

			assert.Nil(t, err)
			assert.NotNil(t, id)

			msg3, err := message.New()
			msg3.SetContent(content3)
			assert.Nil(t, err)

			assert.Nil(t, c.Send(id, msg3))

			msg1, err := message.New()
			msg1.SetContent(content1)
			assert.Nil(t, err)

			assert.Nil(t, c.Send(id, msg1))

			msg2, err := message.New()
			msg2.SetContent(content2)
			assert.Nil(t, err)

			assert.Nil(t, c.Send(id, msg2))

			RegisterTestingT(t)

			recv1, err := message.New()
			assert.Nil(t, err)

			recv2, err := message.New()
			assert.Nil(t, err)

			recv3, err := message.New()
			assert.Nil(t, err)

			Eventually(wg3).Should(Receive(&recv3))
			Eventually(wg1).Should(Receive(&recv1))
			Eventually(wg2).Should(Receive(&recv2))

			Expect(recv3.Content()).To(Equal(content3))
			Expect(recv1.Content()).To(Equal(content1))
			Expect(recv2.Content()).To(Equal(content2))
		*/
	})
}
