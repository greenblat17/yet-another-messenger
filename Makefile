LOCAL_BIN:=$(CURDIR)/bin

PROTOC := PATH="$$PATH:$(LOCAL_BIN)" protoc
AUTH_PROTO_PATH := api/proto/auth/v1
USER_PROTO_PATH := api/proto/user/v1
FRIENDSHIP_PROTO_PATH := api/proto/friendship/v1
NOTIFICATION_PROTO_PATH := api/proto/notification/v1
CHAT_PROTO_PATH := api/proto/chat/v1
VENDOR_PROTO_DIR := vendor.proto

# Установка всех необходимых зависимостей
.PHONY: .bin-deps
.bin-deps:
	$(info Installing binary dependencies...)
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.4.0
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.20.0
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.20.0
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.0.4


# Вендоринг внешних proto файлов
.vendor-proto: $(VENDOR_PROTO_DIR)/google/protobuf $(VENDOR_PROTO_DIR)/google/api $(VENDOR_PROTO_DIR)/validate $(VENDOR_PROTO_DIR)/protoc-gen-openapiv2/options

vendor.proto/protoc-gen-openapiv2/options:
	rm -rf $(VENDOR_PROTO_DIR)/grpc-ecosystem
	git clone -b main --single-branch --depth=1 --filter=tree:0 \
		https://github.com/grpc-ecosystem/grpc-gateway $(VENDOR_PROTO_DIR)/grpc-ecosystem
	cd $(VENDOR_PROTO_DIR)/grpc-ecosystem && \
		git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
		git checkout main
	mkdir -p $(VENDOR_PROTO_DIR)/protoc-gen-openapiv2
	mv $(VENDOR_PROTO_DIR)/grpc-ecosystem/protoc-gen-openapiv2/options $(VENDOR_PROTO_DIR)/protoc-gen-openapiv2
	rm -rf $(VENDOR_PROTO_DIR)/grpc-ecosystem

vendor.proto/google/protobuf:
	rm -rf $(VENDOR_PROTO_DIR)/protobuf
	git clone -b main --single-branch --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf $(VENDOR_PROTO_DIR)/protobuf
	cd $(VENDOR_PROTO_DIR)/protobuf && \
		git sparse-checkout set --no-cone src/google/protobuf && \
		git checkout main
	mkdir -p $(VENDOR_PROTO_DIR)/google
	mv $(VENDOR_PROTO_DIR)/protobuf/src/google/protobuf $(VENDOR_PROTO_DIR)/google
	rm -rf $(VENDOR_PROTO_DIR)/protobuf

vendor.proto/google/api:
	rm -rf $(VENDOR_PROTO_DIR)/googleapis
	git clone -b master --single-branch --depth=1 --filter=tree:0 \
		https://github.com/googleapis/googleapis $(VENDOR_PROTO_DIR)/googleapis
	cd $(VENDOR_PROTO_DIR)/googleapis && \
		git sparse-checkout set --no-cone google/api && \
		git checkout master
	mkdir -p $(VENDOR_PROTO_DIR)/google
	mv $(VENDOR_PROTO_DIR)/googleapis/google/api $(VENDOR_PROTO_DIR)/google
	rm -rf $(VENDOR_PROTO_DIR)/googleapis

vendor.proto/validate:
	rm -rf $(VENDOR_PROTO_DIR)/tmp
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate $(VENDOR_PROTO_DIR)/tmp
	cd $(VENDOR_PROTO_DIR)/tmp && \
		git sparse-checkout set --no-cone validate && \
		git checkout main
	mkdir -p $(VENDOR_PROTO_DIR)/validate
	mv $(VENDOR_PROTO_DIR)/tmp/validate $(VENDOR_PROTO_DIR)/
	rm -rf $(VENDOR_PROTO_DIR)/tmp


# Генерация proto файлов
define generate_proto
	mkdir -p pkg/$1
	$(PROTOC) \
		-I api/proto \
		-I $(VENDOR_PROTO_DIR) \
		$1/$2.proto \
		--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go --go_out=./pkg/$1 --go_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc --go-grpc_out=./pkg/$1 --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway --grpc-gateway_out=./pkg/$1 --grpc-gateway_opt=paths=source_relative,generate_unbound_methods=true \
		--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 --openapiv2_out=./pkg/$1 \
		--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate --validate_out="lang=go,paths=source_relative:pkg/$1"
endef

.PHONY: generate-auth
generate-auth: .bin-deps .vendor-proto
	$(call generate_proto,$(AUTH_PROTO_PATH),auth)

.PHONY: generate-user
generate-user: .bin-deps .vendor-proto
	$(call generate_proto,$(USER_PROTO_PATH),user)

.PHONY: generate-friendship
generate-friendship: .bin-deps .vendor-proto
	$(call generate_proto,$(FRIENDSHIP_PROTO_PATH),friendship)

.PHONY: generate-chat
generate-chat: .bin-deps .vendor-proto
	$(call generate_proto,$(CHAT_PROTO_PATH),chat)

.PHONY: generate-all
generate-all: generate-auth generate-user generate-friendship generate-chat


.PHONY: build up down

build:
	docker compose build

up:
	docker compose up -d

down:
	docker compose down
