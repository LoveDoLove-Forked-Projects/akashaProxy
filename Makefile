NAME=akashaProxy
CC ?= clang
.PHONY: all pack download-dashboard download-mihomo build-webui clean check-deps build-tools default download
all: default

default: check-deps clean download pack

check-deps:
	@command -v curl >/dev/null 2>&1 || { echo >&2 "[ERROR] curl is not installed. Please install curl."; exit 1; }
	@command -v unzip >/dev/null 2>&1 || { echo >&2 "[ERROR] unzip is not installed. Please install unzip."; exit 1; }
	@command -v pnpm >/dev/null 2>&1 || { echo >&2 "[ERROR] pnpm is not installed. Please install pnpm."; exit 1; }
	@command -v go >/dev/null 2>&1 || { echo >&2 "[ERROR] go is not installed. Please install go."; exit 1; }
	@command -v upx >/dev/null 2>&1 || { echo >&2 "[ERROR] upx is not installed. Please install upx."; exit 1; }

download: download-mihomo download-dashboard

pack: build-tools build-webui
	echo "id=akashaProxy\nname=akashaProxy\nversion="$(shell git rev-parse --short HEAD)"\nversionCode="$(shell git log -1 --format=%ct)"\nauthor=akashaProxy developer\ndescription=akasha terminal transparent proxy module that supports tproxy and tun and adds many easy-to-use features. Compatible with Magisk/KernelSU">module/module.prop
	cd module && zip -r ../$(NAME).zip *
	@echo "module pack successfully"

download-mihomo:
	@[ ! -f module/bin ] && mkdir -p module/bin
	remote_mihomo_ver=$$(curl --connect-timeout 5 -Ls "https://github.com/MetaCubeX/mihomo/releases/latest/download/version.txt") && \
	curl --connect-timeout 5 -Ls -o module/bin/mihomo-android-arm64-v8.gz \
	"https://github.com/MetaCubeX/mihomo/releases/latest/download/mihomo-android-arm64-v8-$${remote_mihomo_ver}.gz"
	@echo "mihomo download successfully"

download-dashboard:
	@[ ! -f module/src/bin ] && mkdir -p module/src/zashboard
	curl --connect-timeout 5 -Ls -o module/src/zashboard/dist-cdn-fonts.zip \
	"https://github.com/Zephyruso/zashboard/releases/latest/download/dist-cdn-fonts.zip"
	unzip -o module/src/zashboard/dist-cdn-fonts.zip -d module/src/zashboard/
	mv -f module/src/zashboard/dist/* module/src/zashboard/
	rm -rf module/src/zashboard/dist
	rm -rf module/src/zashboard/dist-cdn-fonts.zip
	@echo "dashboard download successfully"

build-webui:
	cd webui && pnpm i --frozen-lockfile
	cd webui && pnpm build
	mv -f ./webui/out ./module/webroot
	@echo "webui build successfully"

build-tools:
	cd yamlcli && go mod tidy
	cd yamlcli && CGO_ENABLED=0 GOOS=android GOARCH=arm64 go build -trimpath -ldflags="-s -w" -buildvcs=false -o ../module/src/bin/yamlcli
	upx module/src/bin/yamlcli
	@echo "yamlcli build successfully"

clean:
	rm -rf ./module/module.prop
	rm -rf $(NAME).zip
	rm -rf ./module/webroot
	rm -rf ./module/src/zashboard
	rm -rf ./module/bin/*
	rm -rf ./module/src/bin/yamlcli