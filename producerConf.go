package nsq

type ProducerConf struct {
	Address           []string
	ReplicationFactor int32
	Config            *Config
}

func (opt *ProducerConf) init() {
	if len(opt.Address) == 0 {
		opt.Address = []string{
			"localhost:4150",
		}
	}

	if opt.ReplicationFactor < 0 {
		opt.ReplicationFactor = 1
	}
}
