 
all: clean build
build: 
		go build -buildmode=plugin -o filter4.so  function_filter4.go
		go build -buildmode=plugin -o set_header.so function_set_header.go
		go build -buildmode=plugin -o trace.so function_trace.go
clean: 
		go clean
		rm -f *.so