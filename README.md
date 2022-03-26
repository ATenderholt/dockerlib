# dockerlib
This library provides a helper struct around the Docker GO SDK.

```go
import "github.com/ATenderholt/dockerlib"
```

## Example

```go
package example

ctx = context.TODO()
controller, err := dockerlib.NewDockerController()

container := dockerlib.Container{
    Name:  "example",
    Image: "alpine",
    Mounts: []mount.Mount{
        {
            Source:      "/abs/path/to/src",
            Target:      "target",
            Type:        mount.TypeBind,
            ReadOnly:    ...,
            Consistency: ...,
        },
    },
    Ports: map[int]int{
        from: to,
	},
    Command:     []string{...},
    Environment: []string{"KEY=VALUE"},
    Network:     []string{"example"},
}

// make sure network 'example' exists
err := controller.EnsureNetwork(ctx, "example")

// start the container, return channel receives value when given substring is
// found
ready, err := controller.Start(ctx, container, "Container is ready")

<-ready

controller.ShutdownAll(ctx)
controller.CleanupNetworks(ctx)
```

# Documentation

- [func SetLogger(newLogger *zap.Logger)](<#func-setlogger>)
- [type Container](<#type-container>)
  - [func (c Container) PortBindings() (map[nat.Port]struct{}, map[nat.Port][]nat.PortBinding, error)](<#func-container-portbindings>)
  - [func (c Container) String() string](<#func-container-string>)
- [type ContainerError](<#type-containererror>)
  - [func (e ContainerError) Error() string](<#func-containererror-error>)
- [type DockerController](<#type-dockercontroller>)
  - [func NewDockerController() (*DockerController, error)](<#func-newdockercontroller>)
  - [func (controller *DockerController) CleanupNetworks(ctx context.Context) error](<#func-dockercontroller-cleanupnetworks>)
  - [func (controller *DockerController) EnsureImage(ctx context.Context, image string) error](<#func-dockercontroller-ensureimage>)
  - [func (controller *DockerController) EnsureNetwork(ctx context.Context, name string) error](<#func-dockercontroller-ensurenetwork>)
  - [func (controller *DockerController) Remove(ctx context.Context, c Container) error](<#func-dockercontroller-remove>)
  - [func (controller *DockerController) Shutdown(ctx context.Context, c Container) error](<#func-dockercontroller-shutdown>)
  - [func (controller *DockerController) ShutdownAll(ctx context.Context) error](<#func-dockercontroller-shutdownall>)
  - [func (controller *DockerController) Start(ctx context.Context, c Container, ready string) (chan bool, error)](<#func-dockercontroller-start>)
- [type DockerError](<#type-dockererror>)
  - [func (e DockerError) Error() string](<#func-dockererror-error>)
- [type EnsureImageProgress](<#type-ensureimageprogress>)
  - [func (p EnsureImageProgress) String() string](<#func-ensureimageprogress-string>)
- [type EnsureImageProgressDetail](<#type-ensureimageprogressdetail>)
- [type NetworkError](<#type-networkerror>)
  - [func (e NetworkError) Error() string](<#func-networkerror-error>)


## func SetLogger

```go
func SetLogger(newLogger *zap.Logger)
```

## type Container

Container represents a simplified interface for starting a Docker container

```go
type Container struct {
    Name        string
    Image       string
    ID          string
    Mounts      []mount.Mount
    Ports       map[int]int
    Command     []string
    Environment []string
    Network     []string
}
```

### func \(Container\) PortBindings

```go
func (c Container) PortBindings() (map[nat.Port]struct{}, map[nat.Port][]nat.PortBinding, error)
```

PortBindings Helper method to return the structs required to start a Docker container\, or any error

### func \(Container\) String

```go
func (c Container) String() string
```

Returns a simplified string representation

## type ContainerError

```go
type ContainerError struct {
    // contains filtered or unexported fields
}
```

### func \(ContainerError\) Error

```go
func (e ContainerError) Error() string
```

## type DockerController

DockerController is a concrete type that can be used to control Docker containers using its SDK\.

```go
type DockerController struct {
    // contains filtered or unexported fields
}
```

### func NewDockerController

```go
func NewDockerController() (*DockerController, error)
```

NewDockerController is a helper method to create a new instance of a DockerController\.

### func \(\*DockerController\) CleanupNetworks

```go
func (controller *DockerController) CleanupNetworks(ctx context.Context) error
```

### func \(\*DockerController\) EnsureImage

```go
func (controller *DockerController) EnsureImage(ctx context.Context, image string) error
```

EnsureImage is a helper method to pull the specified image to the local machine running Docker\.

### func \(\*DockerController\) EnsureNetwork

```go
func (controller *DockerController) EnsureNetwork(ctx context.Context, name string) error
```

EnsureNetwork Creates a bridge network for the given name if it doesn't already exist\.

### func \(\*DockerController\) Remove

```go
func (controller *DockerController) Remove(ctx context.Context, c Container) error
```

Remove removes the specified \(stopped\) container based on its ID\.

### func \(\*DockerController\) Shutdown

```go
func (controller *DockerController) Shutdown(ctx context.Context, c Container) error
```

Shutdown terminates the specified running Container based on its ID\.

### func \(\*DockerController\) ShutdownAll

```go
func (controller *DockerController) ShutdownAll(ctx context.Context) error
```

ShutdownAll terminates and removes all running containers

### func \(\*DockerController\) Start

```go
func (controller *DockerController) Start(ctx context.Context, c Container, ready string) (chan bool, error)
```

Start is the method used to Start a Docker container using the specified Container c\. It also automatically follows logs and creates a channel that is used to indicate when a running container is ready according to the provided ready string\.

## type DockerError

```go
type DockerError struct {
    // contains filtered or unexported fields
}
```

### func \(DockerError\) Error

```go
func (e DockerError) Error() string
```

## type EnsureImageProgress

EnsureImageProgress is an object to unmarshall JSON returned from Docker during a pull\.

```go
type EnsureImageProgress struct {
    Status         string
    ProgressDetail EnsureImageProgressDetail
    Progress       string
    ID             string
}
```

### func \(EnsureImageProgress\) String

```go
func (p EnsureImageProgress) String() string
```

## type EnsureImageProgressDetail

EnsureImageProgressDetail is an object to help unmarshall JSON returned from Docker during a pull\.

```go
type EnsureImageProgressDetail struct {
    Current int
    Total   int
}
```

## type NetworkError

```go
type NetworkError struct {
    // contains filtered or unexported fields
}
```

### func \(NetworkError\) Error

```go
func (e NetworkError) Error() string
```



Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
