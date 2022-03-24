package dockerlib

type DockerError struct {
	msg       string
	baseError error
}

func (e DockerError) Error() string {
	return e.msg + ": " + e.baseError.Error()
}

type ContainerError struct {
	msg           string
	containerName string
	baseError     error
}

func (e ContainerError) Error() string {
	return e.msg + " " + e.containerName + ": " + e.baseError.Error()
}

type NetworkError struct {
	msg         string
	networkName string
	baseError   error
}

func (e NetworkError) Error() string {
	return e.msg + " " + e.networkName + ": " + e.baseError.Error()
}
