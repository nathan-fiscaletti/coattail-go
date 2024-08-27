package host

func Run(config HostConfig) error {
	return newHost(config).start()
}
