run-locally:
	go run ./cmd/main.go
run-docker:
	docker-compose up --build
generate-mocks-for-link-pgStorage:
	mockgen -source internal/storage/storage/link.go -destination internal/storage/pg/mocks/mocks.go