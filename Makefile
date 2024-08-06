LOCAL_BIN:=$(CURDIR)/bin

PROTOC := PATH="$$PATH:$(LOCAL_BIN)" protoc

# Путь до завендореных protobuf файлов
VENDOR_PROTO_PATH := $(CURDIR)/vendor.protobuf

# Путь до protobuf файлов
AUTH_PROTO_PATH := $(CURDIR)/auth/api/proto/auth
USER_PROTO_PATH := $(CURDIR)/user/api/proto/user
FRIENDSHIP_PROTO_PATH := $(CURDIR)/friendship/api/proto/friendship
CHAT_PROTO_PATH := $(CURDIR)/chat/api/proto/chat

CLIENTS_AUTH_PROTO_PATH := $(CURDIR)/clients/api/proto/auth

# Путь до сгенеренных .pb.go файлов
CHAT_PKG_PROTO_PATH := $(CURDIR)/chat/pkg
AUTH_PKG_PROTO_PATH := $(CURDIR)/auth/pkg
USER_PKG_PROTO_PATH := $(CURDIR)/user/pkg
FRIENDSHIP_PKG_PROTO_PATH := $(CURDIR)/friendship/pkg

# Установка всех необходимых зависимостей
.bin-deps: export GOBIN := $(LOCAL_BIN)
.bin-deps:
	$(info Installing binary dependencies...)

	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.4.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.20.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.20.0
	go install github.com/bufbuild/buf/cmd/buf@v1.32.2
	go install github.com/yoheimuta/protolint/cmd/protolint@latest

# Вендоринг внешних proto файлов
vendor:	.vendor-reset .vendor-googleapis .vendor-google-protobuf .vendor-protovalidate .vendor-protoc-gen-openapiv2 .vendor-tidy

.vendor-reset:
	rm -rf $(VENDOR_PROTO_PATH)
	mkdir -p $(VENDOR_PROTO_PATH)

.vendor-tidy:
	find $(VENDOR_PROTO_PATH) -type f ! -name "*.proto" -delete
	find $(VENDOR_PROTO_PATH) -empty -type d -delete

# Устанавливаем proto описания google/protobuf
.vendor-google-protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf $(VENDOR_PROTO_PATH)/protobuf &&\
	cd $(VENDOR_PROTO_PATH)/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p $(VENDOR_PROTO_PATH)/google
	mv $(VENDOR_PROTO_PATH)/protobuf/src/google/protobuf $(VENDOR_PROTO_PATH)/google
	rm -rf $(VENDOR_PROTO_PATH)/protobuf

# Устанавливаем proto описания validate
.vendor-protovalidate:
	git clone -b main --single-branch --depth=1 --filter=tree:0 \
		https://github.com/bufbuild/protovalidate $(VENDOR_PROTO_PATH)/protovalidate && \
	cd $(VENDOR_PROTO_PATH)/protovalidate
	git checkout
	mv $(VENDOR_PROTO_PATH)/protovalidate/proto/protovalidate/buf $(VENDOR_PROTO_PATH)
	rm -rf $(VENDOR_PROTO_PATH)/protovalidate

# Устанавливаем proto описания google/api
.vendor-googleapis:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/googleapis/googleapis $(VENDOR_PROTO_PATH)/googleapis &&\
	cd $(VENDOR_PROTO_PATH)/googleapis &&\
	git checkout
	mv $(VENDOR_PROTO_PATH)/googleapis/google $(VENDOR_PROTO_PATH)
	rm -rf $(VENDOR_PROTO_PATH)/googleapis

# Устанавливаем proto описания protoc-gen-openapiv2/options
.vendor-protoc-gen-openapiv2:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway $(VENDOR_PROTO_PATH)/grpc-gateway && \
 	cd $(VENDOR_PROTO_PATH)/grpc-gateway && \
	git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
	git checkout
	mkdir -p $(VENDOR_PROTO_PATH)/protoc-gen-openapiv2
	mv $(VENDOR_PROTO_PATH)/grpc-gateway/protoc-gen-openapiv2/options $(VENDOR_PROTO_PATH)/protoc-gen-openapiv2
	rm -rf $(VENDOR_PROTO_PATH)/grpc-gateway

# генерация .go файлов с помощью protoc
define protoc-generate
	mkdir -p $1
	$(PROTOC) -I $(VENDOR_PROTO_PATH) --proto_path=$(CURDIR) \
	--go_out=$1 --go_opt paths=source_relative \
	--go-grpc_out=$1 --go-grpc_opt paths=source_relative \
	--grpc-gateway_out=$1 --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true \
	--openapiv2_out=. --openapiv2_opt logtostderr=true \
	$2/messages.proto $2/service.proto
endef

.PHONY: generate-auth
generate-auth: .bin-deps
	$(call protoc-generate,$(AUTH_PKG_PROTO_PATH),$(AUTH_PROTO_PATH))

.PHONY: generate-user
generate-user: .bin-deps
	$(call protoc-generate,$(USER_PKG_PROTO_PATH),$(USER_PROTO_PATH))

.PHONY: generate-friendship
generate-friendship: .bin-deps
	$(call protoc-generate,$(FRIENDSHIP_PKG_PROTO_PATH),$(FRIENDSHIP_PROTO_PATH))

.PHONY: generate-chat
generate-chat: .bin-deps
	$(call protoc-generate,$(CHAT_PKG_PROTO_PATH),$(CHAT_PROTO_PATH))

.PHONY: generate-all
generate-all: generate-auth generate-user generate-friendship generate-chat

#.PHONY: generate-clients
#generate-clients: .bin-deps
#	$(call protoc-generate,$(CLIENTS_PROTO_PATH),auth,clients)
#	$(call generate_proto,$(CLIENTS_PROTO_PATH),user,clients)
#	$(call generate_proto,$(CLIENTS_PROTO_PATH),friendship,clients)
#	$(call generate_proto,$(CLIENTS_PROTO_PATH),chat,clients)

# Генерация .pb файлов с помощью buf
.buf-generate:
	$(info run buf generate...)
	PATH="$(LOCAL_BIN):$(PATH)" $(LOCAL_BIN)/buf generate

# Генерация кода из protobuf
generate: .bin-deps .buf-generate proto-format

# Форматирование protobuf файлов
proto-format:
	$(info run buf format...)
	$(LOCAL_BIN)/buf format -w

# Линтер
lint:
	$(call proto-lint,$(AUTH_PROTO_PATH))
	$(call proto-lint,$(USER_PROTO_PATH))
	$(call proto-lint,$(CHAT_PROTO_PATH))
	$(call proto-lint,$(FRIENDSHIP_PROTO_PATH))

# Линтер proto файлов
define proto-lint
	$(LOCAL_BIN)/protolint -config_path ./.protolint.yaml $1
endef


.PHONY: \
	.bin-deps \
	.protoc-generate \
	.buf-generate \
	.tidy \
	.vendor-protovalidate \
	.proto-lint \
	proto-format \
	vendor \
	lint

.PHONY: build up down

build:
	docker compose build

up:
	docker compose up -d

down:
	docker compose down
