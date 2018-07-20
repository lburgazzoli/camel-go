all: clean build
build: 
		go build -o camel-go
clean: 
		go clean
		rm -f camel-go