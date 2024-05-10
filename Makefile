compose-up:
	docker-compose up --build -d && docker-compose logs -f
.PHONY: compose-up

compose-down:
	docker-compose down --remove-orphans
.PHONY: compose-down

run:
	go run ./cmd/order-service/main.go
.PHONY: run

send:
	go run ./send_data.go
.PHONY: send

