integration_test:
	docker-compose up -d && go test && docker-compose down
