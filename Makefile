.PHONY: all build clean darwin-amd64 darwin-arm64 linux-arm64 linux-amd64 android

NAME := sub2sing
VERSION := 1.0.0
LDFLAGS := -s -w

all: build

build:
	go build -ldflags "$(LDFLAGS)" -o $(NAME) .

darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(NAME)-darwin-amd64 .

darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(NAME)-darwin-arm64 .

linux-arm64:
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(NAME)-linux-arm64 .

linux-amd64:
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(NAME)-linux-amd64 .

android: linux-arm64
	@echo "推送到 Android: adb push $(NAME)-linux-arm64 /sdcard/"
	@echo "然后在 Termux 中: cp /sdcard/$(NAME)-linux-arm64 ~/ && chmod +x ~/$(NAME)-linux-arm64"

clean:
	rm -f $(NAME) $(NAME)-darwin-* $(NAME)-linux-*
