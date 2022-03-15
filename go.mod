module go.xargs.dev/bindl

go 1.18

require (
	github.com/fatih/color v1.13.0
	github.com/rs/zerolog v1.26.1
	github.com/spf13/cobra v1.3.0
	sigs.k8s.io/yaml v1.3.0
)

require (
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.0.0-20211205182925-97ca703d548d // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

// Lock to v2.3.0 until v3 is out to prevent line wrap
// ref: https://github.com/go-yaml/yaml/pull/670#issuecomment-726666943
replace gopkg.in/yaml.v2 v2.4.0 => gopkg.in/yaml.v2 v2.3.0
