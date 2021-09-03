module github.com/Tamplier2911/gorest/cmd/app

go 1.16

replace github.com/Tamplier2911/gorest/internal => ../../internal

replace github.com/Tamplier2911/gorest/pkg => ../../pkg

require (
	github.com/Tamplier2911/gorest/internal v0.0.0-20210825175559-050d77d5da16
	github.com/Tamplier2911/gorest/pkg v0.0.0-20210825175559-050d77d5da16 // indirect
)
