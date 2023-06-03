package nsq

import (
	_ "unsafe"
)

//go:linkname NewConfig github.com/nsqio/go-nsq.NewConfig
func NewConfig() *Config
