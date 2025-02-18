COMMIT := $(shell git log -1 --format='%H')
GO_VERSION := $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f1,2)

# VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
VERSION := v2.0.0

BUILD_TIME := 202502100800.00 # format [[CC]YY]MMDDhhmm[.ss]

TEAM_ALLOCATION := 165000000000000
ifeq ($(ENV),kaon)
$(info 📑 Using Kaon environment...)
DENOM := tkyve
TEAM_TGE := 2023-02-07T14:00:00
TEAM_FOUNDATION_ADDRESS := kyve1vut528et85755xsncjwl6dx8xakuv26hxgyv0n
TEAM_BCP_ADDRESS := kyve1vut528et85755xsncjwl6dx8xakuv26hxgyv0n
else ifeq ($(ENV),mainnet)
$(info 📑 Using mainnet environment...)
DENOM := ukyve
TEAM_TGE := 2023-03-14T14:03:14
TEAM_FOUNDATION_ADDRESS := kyve1xjpl57p7f49y5gueu7rlfytaw9ramcn5zhjy2g
TEAM_BCP_ADDRESS := kyve1fnh4kghr25tppskap50zk5j385pt65tyyjaraa
endif

ldflags := $(LDFLAGS)
ldflags += -X github.com/cosmos/cosmos-sdk/version.Name=kyve \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=kyved \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X github.com/KYVENetwork/chain/x/global/types.Denom=$(DENOM) \
		  -X github.com/KYVENetwork/chain/x/team/types.TEAM_FOUNDATION_STRING=$(TEAM_FOUNDATION_ADDRESS) \
		  -X github.com/KYVENetwork/chain/x/team/types.TEAM_BCP_STRING=$(TEAM_BCP_ADDRESS) \
		  -X github.com/KYVENetwork/chain/x/team/types.TEAM_ALLOCATION_STRING=$(TEAM_ALLOCATION) \
		  -X github.com/KYVENetwork/chain/x/team/types.TGE_STRING=$(TEAM_TGE) \
		  -w -s
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -ldflags '$(ldflags)' -tags 'ledger' -trimpath -buildvcs=false

.PHONY: proto-setup proto-format proto-lint proto-gen \
	format lint vet test build release dev interchaintest
all: ensure_environment ensure_version proto-all format lint test interchaintest build

###############################################################################
###                                  Build                                  ###
###############################################################################

build: ensure_environment ensure_version
	@echo "🤖 Building kyved..."
	@go build $(BUILD_FLAGS) -o "$(PWD)/build/" ./cmd/kyved
	@echo "✅ Completed build!"

install: ensure_environment ensure_version
	@echo "🤖 Installing kyved..."
	@go install -mod=readonly $(BUILD_FLAGS) ./cmd/kyved
	@echo "✅ Completed installation!"

release: ensure_environment ensure_version
	@echo "🤖 Creating kyved releases (using timestamp $(BUILD_TIME))..."
	@rm -rf release
	@mkdir -p release

	@for b in darwin:amd64 darwin:arm64 linux:amd64 linux:arm64; do \
		os=$$(echo $$b | cut -d':' -f1); \
		arch=$$(echo $$b | cut -d':' -f2); \
		echo "➡️ "$$os" "$$arch""; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build $(BUILD_FLAGS) -o release/kyved_$(ENV)_"$$os"_"$$arch" ./cmd/kyved; \
		touch -a -m -t $(BUILD_TIME) release/kyved_$(ENV)_"$$os"_"$$arch"; \
		sha256sum release/kyved_$(ENV)_"$$os"_"$$arch" >> release/release_$(ENV)_checksum; \
		mv release/kyved_$(ENV)_"$$os"_"$$arch" release/kyved; \
		tar -C release -cf release/kyved_$(ENV)_"$$os"_"$$arch".tar kyved; \
		mv release/kyved  release/kyved_$(ENV)_"$$os"_"$$arch"; \
		touch -a -m -t $(BUILD_TIME) release/kyved_$(ENV)_"$$os"_"$$arch".tar; \
		gzip release/kyved_$(ENV)_"$$os"_"$$arch".tar; \
		sha256sum release/kyved_$(ENV)_"$$os"_"$$arch".tar.gz >> release/release_$(ENV)_checksum; \
	done

	@echo "✅  Completed release creation!"

###############################################################################
###                               Docker Build                              ###
###############################################################################

# Build a release image
docker-image:
	@DOCKER_BUILDKIT=1 docker build -t kyve-network/kyve:${VERSION} .
	@echo "✅ Completed docker image build!"

# Build a release nonroot image
docker-image-nonroot:
	@DOCKER_BUILDKIT=1 docker build \
		--build-arg IMG_TAG="nonroot" \
		-t kyve-network/kyve:${VERSION}-nonroot .
	@echo "✅ Completed docker image build! (nonroot)"

# Build a release debug image
docker-image-debug:
	@DOCKER_BUILDKIT=1 docker build \
		--build-arg IMG_TAG="debug" \
		-t kyve-network/kyve:${VERSION}-debug .
	@echo "✅ Completed docker image build! (debug)"

# Build a release debug-nonroot image
docker-image-debug-nonroot:
	@DOCKER_BUILDKIT=1 docker build \
		--build-arg IMG_TAG="debug-nonroot" \
		-t kyve-network/kyve:${VERSION}-debug-nonroot .
	@echo "✅ Completed docker image build! (debug-nonroot)"

###############################################################################
###                                 Checks                                  ###
###############################################################################

ensure_environment:
ifndef ENV
	$(error ❌  Please specify a build environment..)
endif

ensure_version:
ifneq ($(GO_VERSION),1.22)
	$(error ❌  Please run Go v1.22.x..)
endif

###############################################################################
###                               Development                               ###
###############################################################################

dev:
	@ignite chain serve --reset-once --skip-proto --verbose

dev-continue:
	@ignite chain serve --skip-proto --verbose

###############################################################################
###                          Formatting & Linting                           ###
###############################################################################

gofumpt_cmd=mvdan.cc/gofumpt
golangci_lint_cmd=github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2

format:
	@echo "🤖 Running formatter..."
	@go run $(gofumpt_cmd) -l -w .
	@echo "✅ Completed formatting!"

lint:
	@echo "🤖 Running linter..."
	# TODO temporarily disabled govet
	@go run $(golangci_lint_cmd) run --exclude-dirs scripts --timeout=10m -D govet
	@echo "✅ Completed linting!"

###############################################################################
###                                Protobuf                                 ###
###############################################################################

BUF_VERSION=1.20.0

proto-all: proto-format proto-lint proto-gen

proto-format:
	@echo "🤖 Running protobuf formatter..."
	@docker run --rm --volume "$(PWD)":/workspace --workdir /workspace \
		bufbuild/buf:$(BUF_VERSION) format --diff --write
	@echo "✅ Completed protobuf formatting!"

proto-gen:
	@echo "🤖 Generating code from protobuf..."
	@docker run --rm \
		--volume "$(PWD)":/workspace 														`# Mount this repo as workspace` \
		--volume "$$(go list -m -f '{{.Dir}}' github.com/cosmos/cosmos-sdk)":/cosmos-sdk:ro `# Mount cosmos-sdk from installed go modules` \
		--volume "$$(go list -m -f '{{.Dir}}' github.com/cosmos/ibc-go/v8)":/ibc:ro 		`# Mount ibc-go from installed go modules` \
		--workdir /workspace \
		kyve-proto sh ./proto/generate.sh
	@echo "✅ Completed code generation!"

proto-lint:
	@echo "🤖 Running protobuf linter..."
	@docker run --rm --volume "$(PWD)":/workspace --workdir /workspace \
		bufbuild/buf:$(BUF_VERSION) lint
	@echo "✅ Completed protobuf linting!"

proto-setup:
	@echo "🤖 Setting up protobuf environment..."
	@docker build --rm --tag kyve-proto:latest --file proto/Dockerfile \
		--build-arg USER_ID=$$(id -u) --build-arg USER_GID=$$(id -g) .
	@echo "✅ Setup protobuf environment!"

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

ensure_heighliner:
	@which heighliner > /dev/null || (echo "❌ Heighliner not found. Please install it by running 'make heighliner-setup'." && exit 1)
	@docker inspect kaon:local > /dev/null || (echo "❌ Kaon image not found. Please build it by running 'make heighliner'." && exit 1)
	@docker inspect kyve:local > /dev/null || (echo "❌ Kyve image not found. Please build it by running 'make heighliner'." && exit 1)

heighliner:
	@echo "🤖 Building Kaon image..."
	@heighliner build --chain kaon --local 1> /dev/null
	@echo "✅ Completed build!"

	@echo "🤖 Building Kyve image..."
	@heighliner build --chain kyve --local 1> /dev/null
	@echo "✅ Completed build!"

heighliner-setup:
	@echo "🤖 Installing Heighliner..."
	@git clone https://github.com/strangelove-ventures/heighliner.git /tmp/heighliner
	@cd /tmp/heighliner && go install && cd ..
	@rm -rf heighliner
	@echo "✅ Completed installation!"

test:
	@echo "🤖 Running tests..."
	@go test -cover -mod=readonly ./x/...
	@echo "✅ Completed tests!"

interchaintest: # ensure_heighliner
#   TODO Heighliner is currently broken
#	@echo "🤖 Running interchain tests..."
#	@go test -mod=readonly ./interchaintest/...
#	@echo "✅ Completed interchain tests!"
