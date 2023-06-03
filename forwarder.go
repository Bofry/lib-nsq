package nsq

type Forwarder struct {
	*Producer
}

func NewForwarder(conf *ProducerConf) (*Forwarder, error) {
	producer, err := NewProducer(conf)
	if err != nil {
		return nil, err
	}
	instance := &Forwarder{
		Producer: producer,
	}
	return instance, nil
}

func (f *Forwarder) Runner() *ForwarderRunner {
	return &ForwarderRunner{
		handle: f,
	}
}
