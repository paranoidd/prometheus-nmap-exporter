build:
	go fmt .
	go build .

docker_build:
	docker build -t prometheus-nmap-exporter:latest .

docker_run_example:
	docker run --rm -it \
		-v ./config.json.example:/app/config.json \
		prometheus-nmap-exporter:latest
