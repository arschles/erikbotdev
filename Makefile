.PHONY: run
run:
	go build -o erikbotdev . && ./erikbotdev run -s

.PHONY: runserver
runserver:
	go build -o erikbotserver ./cmd/server && ./erikbotserver serve