module github.com/julianGoh17/simple-e2e/framework

go 1.14

require (
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/containerd/containerd v1.4.0 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	// https://github.com/moby/moby/issues/40185#issuecomment-550443447 Need to use git has because docker stopped using semantic versioning
	github.com/docker/docker v17.12.0-ce-rc1.0.20200821074627-7ae5222c72cc+incompatible
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/rs/zerolog v1.19.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.1
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.1
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/testify v1.6.1
	golang.org/x/net v0.0.0-20200813134508-3edf25e44fcc // indirect
	google.golang.org/grpc v1.31.0 // indirect
	gopkg.in/yaml.v2 v2.3.0
)
