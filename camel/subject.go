package camel

// ==========================
//
//
//
// ==========================

// NewSubject --
func NewSubject() *Subject {
	p := Subject{}
	p.in = make(chan *Exchange)

	return &p
}

// Subject --
type Subject struct {
	in chan *Exchange
}

// Subscribe --
func (subject *Subject) Subscribe(consumer func(*Exchange)) *Subject {
	go func() {
		for exchange := range subject.in {
			consumer(exchange)
		}
	}()

	return subject
}

// SubscribeTo --
func (subject *Subject) SubscribeTo(source *Subject) *Subject {
	source.Subscribe(func(e *Exchange) {
		subject.Publish(e)
	})

	return subject
}

// SubscribeWithTransformer --
func (subject *Subject) SubscribeWithProcessor(source *Subject, processor Processor, processors ...Processor) *Subject {
	source.Subscribe(func(e *Exchange) {
		processor.Process(e)

		for _, proc := range processors {
			proc.Process(e)
		}

		subject.Publish(e)
	})

	return subject
}

// SubscribeTo --
func (subject *Subject) SubscribeWithTransformer(source *Subject, processor Transformer, processors ...Transformer) *Subject {
	source.Subscribe(func(e *Exchange) {
		e = processor.Transform(e)

		for _, proc := range processors {
			e = proc.Transform(e)
		}

		subject.Publish(e)
	})

	return subject
}

// SubscribeTo --
func (subject *Subject) SubscribeWithPredicate(source *Subject, processor Predicate, processors ...Predicate) *Subject {
	source.Subscribe(func(e *Exchange) {
		if ok := processor.Test(e); !ok {
			return
		}

		for _, proc := range processors {
			if ok := proc.Test(e); !ok {
				return
			}
		}

		subject.Publish(e)
	})

	return subject
}

// Publish --
func (subject *Subject) Publish(exchange *Exchange) *Subject {
	subject.in <- exchange

	return subject
}
