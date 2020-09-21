build/uup: *.go cmd/*.go
	@mkdir -p build
	go build -o build/uup cmd/uup.go

install: *.go cmd/*.go
	go install cmd/uup.go
