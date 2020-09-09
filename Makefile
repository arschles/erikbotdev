.PHONY: run
run:
	go build -o erikbotdev . && ./erikbotdev run -s

.PHONY: runserver
runserver:
	go build -o erikbotserver ./cmd/server && ./erikbotserver serve

.PHONY: dockerbuildserver
dockerbuildserver:
	docker build -t arschles/erikbotserver .

.PHONY: dockerrunserver
dockerrunserver:
	docker run -e PORT=9090 -p 9090:9090 -e ERIKBOTDEV_CONFIG_FILE_NAME=/configs/aaronbot5000.json --rm arschles/erikbotserver