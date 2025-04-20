SERVICE_NAME := template-grpc-go-service
GIT_REMOTE := origin
GIT_BRANCH := master
GIT_DIR := /home/omnia/go-service
GIT_DEPLOY_USER := gitlab+deploy-token-1
GIT_DEPLOY_TOKEN := gldt-akDKStuwmrxNtVx8nHmS
# Docker Configuration
DOCKER_COMPOSE := docker compose
DOCKER_PROJECT_NAME := go-service
DOCKER_SERVICE_NAME := goservice 
.PHONY: help start stop restart status enable disable logs logs-follow \
        git-fetch git-pull git-status git-update \
        docker-build docker-up docker-down docker-logs docker-restart \
        full-update
help:
	@echo "Systemd Service Commands:"
	@echo "  make start     - Start the service"
	@echo "  make stop      - Stop the service"
	@echo "  make restart   - Restart the service"
	@echo "  make status    - Show service status"
	@echo "  make enable    - Enable service autostart"
	@echo "  make disable   - Disable service autostart"
	@echo "  make logs      - Show service logs (last 50 lines)"
	@echo "  make logs-follow - Follow service logs in real-time"
	@echo ""
	@echo "Git Commands (working in $(GIT_DIR)):"
	@echo "  make git-fetch - Fetch latest changes from remote"
	@echo "  make git-pull  - Pull latest changes (fetch + merge)"
	@echo "  make git-status - Show git status"
	@echo "  make git-update - Full update (fetch, pull, restart service)"
	@echo ""
	@echo "Docker Commands:"
	@echo "  make docker-build - Build Docker containers"
	@echo "  make docker-up    - Start Docker containers"
	@echo "  make docker-down  - Stop Docker containers"
	@echo "  make docker-logs  - View Docker container logs"
	@echo "  make docker-restart - Restart Docker containers"
	@echo ""
	@echo "Combined Commands:"
	@echo "  make full-update - Git pull + Docker rebuild + restart"
### Systemd Service Management ###
start:
	systemctl --user start $(SERVICE_NAME)
stop:
	systemctl --user stop $(SERVICE_NAME)
restart:
	systemctl --user restart $(SERVICE_NAME)
status:
	systemctl --user status $(SERVICE_NAME)
enable:
	systemctl --user enable $(SERVICE_NAME)
disable:
	systemctl --user disable $(SERVICE_NAME)
logs:
	journalctl --user -u $(SERVICE_NAME) -n 50 --no-pager
logs-follow:
	journalctl --user -u $(SERVICE_NAME) -f
### Git Commands (executed in GIT_DIR) ###
git-fetch:
	@echo "Fetching updates using deploy token..."
	git -C $(GIT_DIR) fetch https://$(GIT_DEPLOY_USER):$(GIT_DEPLOY_TOKEN)@$(GIT_REMOTE) $(GIT_BRANCH)
git-pull:
	@echo "Pulling updates using deploy token..."
	git -C $(GIT_DIR) pull https://$(GIT_DEPLOY_USER):$(GIT_DEPLOY_TOKEN)@$(GIT_REMOTE) $(GIT_BRANCH)
git-status:
	git -C $(GIT_DIR) status
git-update: git-fetch git-pull restart
	@echo "System updated and service restarted"
### Docker Commands ###
docker-build:
	@echo "Building Docker containers..."
	$(DOCKER_COMPOSE) -p $(DOCKER_PROJECT_NAME) build $(DOCKER_SERVICE_NAME)
docker-up:
	@echo "Starting Docker containers..."
	$(DOCKER_COMPOSE) -p $(DOCKER_PROJECT_NAME) up -d
docker-down:
	@echo "Stopping Docker containers..."
	$(DOCKER_COMPOSE) -p $(DOCKER_PROJECT_NAME) down
docker-logs:
	@echo "Showing Docker container logs..."
	$(DOCKER_COMPOSE) -p $(DOCKER_PROJECT_NAME) logs -f $(DOCKER_SERVICE_NAME)
docker-restart: docker-down docker-up
	@echo "Docker containers restarted"
### Combined Commands ###
full-update: git-pull docker-build docker-restart
	@echo "Full update completed: Git pull + Docker rebuild + restart"
### Extended Systemd Commands ###
daemon-reload:
	systemctl --user daemon-reload
logs-all:
	journalctl --user -u $(SERVICE_NAME) --no-pager
logs-today:
	journalctl --user -u $(SERVICE_NAME) --since today --no-pager
logs-err:
	journalctl --user -u $(SERVICE_NAME) -p err --no-pager