module sigs.k8s.io/kind

// NOTE: This is the go language version, NOT the compiler version.
//
// This controls the *minimum* required go version and therefore available Go
// language features.
//
// See ./.go-version for the go compiler version used when building binaries
//
// https://go.dev/doc/modules/gomod-ref#go
go 1.22

require (
	al.essio.dev/pkg/shellescape v1.6.0
	github.com/BurntSushi/toml v1.5.0
	github.com/evanphx/json-patch/v5 v5.9.11
	github.com/mattn/go-isatty v0.0.20
	github.com/pelletier/go-toml v1.9.5
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.9.1
	github.com/spf13/pflag v1.0.7
	go.yaml.in/yaml/v3 v3.0.4
	sigs.k8s.io/yaml v1.6.0
)

// test-only transitive deps, these are used by sigs.k8s.io/yaml's tests
require (
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	golang.org/x/sys v0.21.0 // indirect
)
