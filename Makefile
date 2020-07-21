APP=Go-Chat
APP_EXECUTABLE="./out/$(APP)"


fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	@for p in $(UNIT_TEST_PACKAGES); do \
		echo "==> Linting $$p"; \
		golint $$p | { grep -vwE "exported (var|function|method|type|const) \S+ should have comment" || true; } \
	done

install:
	go install ./...

compile:
	GO111MODULE=on go mod vendor
	mkdir -p out/
	go build -o $(APP_EXECUTABLE)

build: fmt vet lint compile