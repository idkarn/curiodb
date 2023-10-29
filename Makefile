run:
	go run .

watch:
	air

build:
	go build -o bin/curiodb

docker:
	docker build -t curiodb:0.1 .
	docker run -p 3141:3141 -d curiodb:0.1