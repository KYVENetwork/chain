COMMIT := $(shell git log -1 --format='%H')
DENOM := tkyve
VERSION := 1.0.0-rc0 # $(shell echo $(shell git describe --tags) | sed 's/^v//')

# Team Module
TEAM_AUTHORITY_ADDRESS := kyve1vut528et85755xsncjwl6dx8xakuv26hxgyv0n
TEAM_ALLOCATION := 165000000000000
TEAM_TGE := 2023-02-07T14:00:00

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=kyve \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=kyved \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X github.com/KYVENetwork/chain/x/global/types.Denom=$(DENOM) \
		  -X github.com/KYVENetwork/chain/x/team/types.TEAM_AUTHORITY_STRING=$(TEAM_AUTHORITY_ADDRESS) \
		  -X github.com/KYVENetwork/chain/x/team/types.TEAM_ALLOCATION_STRING=$(TEAM_ALLOCATION) \
		  -X github.com/KYVENetwork/chain/x/team/types.TGE_STRING=$(TEAM_TGE)

# TODO(@john): Are the missing flags needed?
BUILD_FLAGS := -ldflags '$(ldflags)' -tags 'ledger' -trimpath

.PHONY: proto-setup proto-format proto-lint proto-gen \
	format lint vet test build release dev
all: proto-all format lint test build

###############################################################################
###                                  Build                                  ###
###############################################################################

build:
	@echo "🤖 Building kyved..."
	@go build $(BUILD_FLAGS) -o "$(PWD)/build/" ./cmd/kyved
	@echo "✅ Completed build!"

install:
	@echo "🤖 Installing kyved..."
	@go install -mod=readonly $(BUILD_FLAGS) ./cmd/kyved
	@echo "✅ Completed installation!"

release:
	@echo "🤖 Creating kyved releases..."
	@rm -rf release
	@mkdir -p release

	@GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) ./cmd/kyved
	@tar -czf release/kyved_darwin_amd64.tar.gz kyved
	@sha256sum release/kyved_darwin_amd64.tar.gz >> release/release_checksum

	@GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) ./cmd/kyved
	@tar -czf release/kyved_darwin_arm64.tar.gz kyved
	@sha256sum release/kyved_darwin_arm64.tar.gz >> release/release_checksum

	@GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) ./cmd/kyved
	@tar -czf release/kyved_linux_amd64.tar.gz kyved
	@sha256sum release/kyved_linux_amd64.tar.gz >> release/release_checksum

	@GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) ./cmd/kyved
	@tar -czf release/kyved_linux_arm64.tar.gz kyved
	@sha256sum release/kyved_linux_arm64.tar.gz >> release/release_checksum

	@rm kyved
	@echo "✅ Completed release creation!"

###############################################################################
###                               Development                               ###
###############################################################################

# TODO(@john): Switch to the Docker image?
dev:
	@ignite chain serve --reset-once --skip-proto --verbose

###############################################################################
###                          Formatting & Linting                           ###
###############################################################################

gofumpt_cmd=mvdan.cc/gofumpt
golangci_lint_cmd=github.com/golangci/golangci-lint/cmd/golangci-lint

format:
	@echo "🤖 Running formatter..."
	@go run $(gofumpt_cmd) -l -w .
	@echo "✅ Completed formatting!"

lint:
	@echo "🤖 Running linter..."
	@go run $(golangci_lint_cmd) run --skip-dirs scripts --timeout=10m
	@echo "✅ Completed linting!"

# TODO(@john): Can we remove this since we use GolangCI?
vet:
	@echo "🤖 Running vet..."
	@go vet ./...
	@echo "✅ Completed vet!"

###############################################################################
###                                Protobuf                                 ###
###############################################################################

BUF_VERSION=1.13.1

proto-all: proto-format proto-lint proto-gen

proto-format:
	@echo "🤖 Running protobuf formatter..."
	@docker run --volume "$(PWD)":/workspace --workdir /workspace \
		bufbuild/buf:$(BUF_VERSION) format --diff --write
	@echo "✅ Completed protobuf formatting!"

proto-gen:
	@echo "🤖 Generating code from protobuf..."
	@docker run --rm --volume "$(PWD)":/workspace --workdir /workspace \
		kyve-proto sh ./proto/generate.sh
	@echo "✅ Completed code generation!"

proto-lint:
	@echo "🤖 Running protobuf linter..."
	@docker run --volume "$(PWD)":/workspace --workdir /workspace \
		bufbuild/buf:$(BUF_VERSION) lint
	@echo "✅ Completed protobuf linting!"

proto-setup:
	@echo "🤖 Setting up protobuf environment..."
	@docker build --rm --tag kyve-proto:latest --file proto/Dockerfile .
	@echo "✅ Setup protobuf environment!"

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

test:
	@echo "🤖 Running tests..."
	@go test -cover -mod=readonly ./x/...
	@echo "✅ Completed tests!"
