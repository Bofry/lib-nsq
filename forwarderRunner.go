package nsq

type ForwarderRunner struct {
	handle *Forwarder
}

func (r *ForwarderRunner) Start() {
	r.handle.logger.Println("Started")
}

func (r *ForwarderRunner) Stop() {
	r.handle.logger.Println("Stopping")
	r.handle.Close()
	r.handle.logger.Println("Stopped")
}
