build:
	go fmt .
	go build .

docker_build:
	docker build -t prometheus-nmap-exporter:latest .

docker_run_example: docker_build
	docker run --rm -it \
		-v $(shell pwd)/config.json.example:/app/config.json \
		-p 9777:9777 \
		prometheus-nmap-exporter:latest
