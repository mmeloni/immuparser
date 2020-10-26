module github.com/mmeloni/immuparser

go 1.13

require (
	github.com/codenotary/immudb v0.8.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/schollz/progressbar/v3 v3.6.2
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.0
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897 // indirect
	golang.org/x/sys v0.0.0-20201026163216-e075d9370641 // indirect
	google.golang.org/grpc v1.29.1
)

replace github.com/codenotary/immudb v0.8.0 => github.com/codenotary/immudb v0.0.0-20201015224644-7387902c29a9
