.PHONY: help init down

help:
	@echo "Available commands:"
	@echo "  init    Initialize and start the database using docker-compose"
	@echo "  down    Stop and remove the database container"

init:
	@echo "Starting database..."
	docker-compose up -d
	@echo "Database started."

down:
	@echo "Stopping database..."
	docker-compose down
	@echo "Database stopped."
