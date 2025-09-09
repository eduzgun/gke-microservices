# Read Terraform outputs (assumes `terraform apply` already run)
TF_PROJECT := $(shell cd terraform && terraform output -raw project_id 2>/dev/null)
TF_REGION := $(shell cd terraform && terraform output -raw region 2>/dev/null)
TF_REPO := $(shell cd terraform && terraform output -raw artifact_registry_repo 2>/dev/null)
AR := $(TF_REGION)-docker.pkg.dev/$(TF_PROJECT)/$(TF_REPO)

# Fallback if terraform not applied yet
ifeq ($(TF_PROJECT),)
$(error Run 'make terraform-apply' first to set up Artifact Registry)
endif

# Build and tag all 4 images for GKE, then push them
build:
	@echo "Building and pushing all images to $(AR)"
	docker build -t $(AR)/gke-backend:latest -f Dockerfile.backend .
	docker build -t $(AR)/gke-session:latest -f Dockerfile.session .
	docker build -t $(AR)/gke-migrate:latest -f Dockerfile.migrate .
	cd frontend && docker build -t $(AR)/gke-frontend:latest -f Dockerfile .
	@echo "Pushing all images..."
	docker push $(AR)/gke-backend:latest
	docker push $(AR)/gke-session:latest
	docker push $(AR)/gke-migrate:latest
	docker push $(AR)/gke-frontend:latest
	@echo "All images built and pushed."

# Generate protobuf code
proto-gen:
	mkdir -p ./internal/proto/gen/go/session
	protoc \
		--go_out=./internal/proto/gen/go/session \
		--go_opt=paths=source_relative \
		--go-grpc_out=./internal/proto/gen/go/session \
		--go-grpc_opt=paths=source_relative \
		--proto_path=./internal/proto \
		./internal/proto/session.proto

# Clean generated protobuf files
clean-proto:
	rm -rf ./internal/proto/gen/

# Generate and run session service
run-session: proto-gen
	go run cmd/session-service/main.go

# Generate and run main app
run-backend: go run cmd/app/main.go

run-frontend:
	cd frontend && npm run dev

# Run all services for local dev
dev: proto-gen
	@echo "Starting backend, session, and frontend..."
	cd frontend && npm run dev & go run cmd/backend/main.go & go run cmd/session/main.go & wait

deploy-infra:
	./deploy.template.sh

# Destroy EVERYTHING: K8s, images, Terraform
destroy:
	@echo "Destroying Kubernetes resources..."
	kubectl delete all --all


	@echo "Deleting Docker images from Artifact Registry..."
	gcloud artifacts docker images delete $(AR)/gke-backend:latest --quiet || true
	gcloud artifacts docker images delete $(AR)/gke-session:latest --quiet || true
	gcloud artifacts docker images delete $(AR)/gke-migrate:latest --quiet || true
	gcloud artifacts docker images delete $(AR)/gke-frontend:latest --quiet || true

	@echo "Destroying Terraform infrastructure..."
	cd terraform && terraform destroy -auto-approve

	@echo "Full destroy complete."

# Declare all phony targets
.PHONY: build proto-gen clean-proto run-session run-backend run-frontend dev deploy-infra terraform-apply destroy