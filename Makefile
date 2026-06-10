# ============================================================
# task-mgr — Makefile
# ============================================================

REGISTRY   := homelab:30500
APP        := task-mgr
MIGRATE    := task-mgr-migrate
NAMESPACE  := task-mgr

# ── Versjonering ─────────────────────────────────────────────
VERSION    := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
MAJOR      := $(shell echo $(VERSION) | cut -d. -f1 | tr -d 'v')
MINOR      := $(shell echo $(VERSION) | cut -d. -f2)
PATCH      := $(shell echo $(VERSION) | cut -d. -f3)
BRANCH     := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT     := $(shell git rev-parse --short HEAD)

ifeq ($(BRANCH), main)
	IMAGE_TAG := $(VERSION)
else
	IMAGE_TAG := $(VERSION)-$(subst /,-,$(BRANCH))-$(COMMIT)
endif

APP_IMAGE     := $(REGISTRY)/$(APP):$(IMAGE_TAG)
MIGRATE_IMAGE := $(REGISTRY)/$(MIGRATE):$(IMAGE_TAG)

# ── Farger ───────────────────────────────────────────────────
CYAN  := $(shell printf '\033[0;36m')
GREEN := $(shell printf '\033[0;32m')
RESET := $(shell printf '\033[0m')

.PHONY: help version build build-migrate push push-migrate ship \
        up down migrate test vet clean \
        bump-patch bump-minor bump-major \
        branch-fix branch-feature branch-major \
        deploy rollout

# ── Help ─────────────────────────────────────────────────────
help:
	@echo ""
	@echo "$(CYAN)task-mgr — tilgjengelige kommandoer$(RESET)"
	@echo ""
	@echo "  $(GREEN)Versjon$(RESET)"
	@echo "    version          Vis gjeldende versjon og image-tag"
	@echo "    bump-patch       Bump patch-versjon (v0.0.x)"
	@echo "    bump-minor       Bump minor-versjon (v0.x.0)"
	@echo "    bump-major       Bump major-versjon (vx.0.0)"
	@echo ""
	@echo "  $(GREEN)Branching$(RESET)"
	@echo "    branch-fix       Bump patch og lag fix/NAME branch"
	@echo "    branch-feature   Bump minor og lag feature/NAME branch"
	@echo "    branch-major     Bump major og lag major/NAME branch"
	@echo "    NAME=<navn>      Bruk med branch-targets"
	@echo ""
	@echo "  $(GREEN)Bygg og deploy$(RESET)"
	@echo "    build            Bygg app Docker-image"
	@echo "    build-migrate    Bygg migrate Docker-image"
	@echo "    push             Push app-image til registry"
	@echo "    push-migrate     Push migrate-image til registry"
	@echo "    deploy           Apply alle k8s-manifester"
	@echo "    rollout          Restart deployment i k3s"
	@echo "    ship             Bygg, push og deploy i én kommando"
	@echo ""
	@echo "  $(GREEN)Utvikling$(RESET)"
	@echo "    up               Start docker-compose"
	@echo "    down             Stopp docker-compose"
	@echo "    migrate          Kjør migrasjoner lokalt"
	@echo "    test             Kjør alle tester"
	@echo "    vet              Kjør go vet"
	@echo "    clean            Slett kompilerte filer"
	@echo ""

# ── Versjon ───────────────────────────────────────────────────
version:
	@echo "$(CYAN)Versjon:$(RESET)   $(VERSION)"
	@echo "$(CYAN)Branch:$(RESET)    $(BRANCH)"
	@echo "$(CYAN)Commit:$(RESET)    $(COMMIT)"
	@echo "$(CYAN)Image-tag:$(RESET) $(IMAGE_TAG)"

bump-patch:
	@NEW=v$(MAJOR).$(MINOR).$(shell expr $(PATCH) + 1); \
	git tag $$NEW; \
	echo "$(GREEN)Bumped til $$NEW$(RESET)"

bump-minor:
	@NEW=v$(MAJOR).$(shell expr $(MINOR) + 1).0; \
	git tag $$NEW; \
	echo "$(GREEN)Bumped til $$NEW$(RESET)"

bump-major:
	@NEW=v$(shell expr $(MAJOR) + 1).0.0; \
	git tag $$NEW; \
	echo "$(GREEN)Bumped til $$NEW$(RESET)"

# ── Branching ─────────────────────────────────────────────────
branch-fix: bump-patch
	$(if $(NAME),,$(error NAME er påkrevd. Bruk: make branch-fix NAME=beskrivelse))
	git checkout -b fix/$(NAME)
	@echo "$(GREEN)Opprettet branch fix/$(NAME)$(RESET)"

branch-feature: bump-minor
	$(if $(NAME),,$(error NAME er påkrevd. Bruk: make branch-feature NAME=beskrivelse))
	git checkout -b feature/$(NAME)
	@echo "$(GREEN)Opprettet branch feature/$(NAME)$(RESET)"

branch-major: bump-major
	$(if $(NAME),,$(error NAME er påkrevd. Bruk: make branch-major NAME=beskrivelse))
	git checkout -b major/$(NAME)
	@echo "$(GREEN)Opprettet branch major/$(NAME)$(RESET)"

# ── Bygg ──────────────────────────────────────────────────────
build:
	@echo "$(CYAN)Bygger $(APP_IMAGE)...$(RESET)"
	docker build -t $(APP_IMAGE) .
	@echo "$(GREEN)Ferdig: $(APP_IMAGE)$(RESET)"

build-migrate:
	@echo "$(CYAN)Bygger $(MIGRATE_IMAGE)...$(RESET)"
	docker build -f Dockerfile.migrate -t $(MIGRATE_IMAGE) .
	@echo "$(GREEN)Ferdig: $(MIGRATE_IMAGE)$(RESET)"

# ── Push ──────────────────────────────────────────────────────
push: build
	@echo "$(CYAN)Pusher $(APP_IMAGE)...$(RESET)"
	docker push $(APP_IMAGE)
	@echo "$(GREEN)Ferdig$(RESET)"

push-migrate: build-migrate
	@echo "$(CYAN)Pusher $(MIGRATE_IMAGE)...$(RESET)"
	docker push $(MIGRATE_IMAGE)
	@echo "$(GREEN)Ferdig$(RESET)"

# ── Deploy ────────────────────────────────────────────────────
deploy:
	@echo "$(CYAN)Deployer til k3s...$(RESET)"
	kubectl apply -f k8s/
	@echo "$(GREEN)Ferdig$(RESET)"

rollout:
	@echo "$(CYAN)Restarter deployment...$(RESET)"
	kubectl rollout restart deployment/$(APP) -n $(NAMESPACE)
	kubectl rollout status deployment/$(APP) -n $(NAMESPACE)

ship: push push-migrate deploy rollout
	@echo "$(GREEN)Ship fullført — $(IMAGE_TAG) kjører i k3s$(RESET)"

# ── Utvikling ─────────────────────────────────────────────────
up:
	docker compose up -d

down:
	docker compose down

migrate:
	migrate -path migrations -database "$(DATABASE_URL)" up

test:
	go test -v ./...

vet:
	go vet ./...

clean:
	rm -f $(APP)
