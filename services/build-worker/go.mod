module helios/build-worker

go 1.23.0

toolchain go1.24.3

require (
	github.com/nats-io/nats.go v1.46.1
	helios v0.0.0-00010101000000-000000000000
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/nats-io/nkeys v0.4.11 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/rs/zerolog v1.34.0
	github.com/stretchr/testify v1.9.0
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
)

replace helios => ../
