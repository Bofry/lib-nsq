package nsq

import "log"

type ProducerConfig struct {
	Address           []string
	ReplicationFactor int32
	Config            *Config
	Logger            *log.Logger
}

func (c *ProducerConfig) init() {
	if len(c.Address) == 0 {
		c.Address = []string{
			"localhost:4150",
		}
	}

	if c.ReplicationFactor < 0 {
		c.ReplicationFactor = 1
	}

	if c.Config == nil {
		c.Config = NewConfig()
	}
}

func (c *ProducerConfig) Validate() error {
	if c.Config != nil {
		return c.Config.Validate()
	}
	return nil
}
