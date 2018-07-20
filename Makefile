all: clean build
build: 
		go build -o desertship
clean: 
		go clean
		rm -f desertship