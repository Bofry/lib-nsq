package nsq

type ForwarderRunner struct {
	handle *Forwarder
}

func (r *ForwarderRunner) Start() {
	r.handle.Logger.Println("Started")
}

func (r *ForwarderRunner) Stop() {
	r.handle.Logger.Println("Stopping")
	r.handle.Close()
	r.handle.Logger.Println("Stopped")
}
