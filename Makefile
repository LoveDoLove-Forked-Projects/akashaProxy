NAME=akashaProxy
CC ?= clang
.PHONY: all pack download-dashboard download-mihomo build-webui clean check-deps build-tools default


all: default

default: check-deps \
	clean \
	download-mihomo \
	download-dashboard \
	build-webui \
	build-tools \
	pack

check-deps:
	@command -v curl >/dev/null 2>&1 || { echo >&2 "[ERROR] curl is not installed. Please install curl."; exit 1; }
	@command -v unzip >/dev/null 2>&1 || { echo >&2 "[ERROR] unzip is not installed. Please install unzip."; exit 1; }
	@command -v yarn >/dev/null 2>&1 || { echo >&2 "[ERROR] yarn is not installed. Please install yarn."; exit 1; }
	@command -v pnpm >/dev/null 2>&1 || { echo >&2 "[ERROR] pnpm is not installed. Please install pnpm."; exit 1; }
	@command -v go >/dev/null 2>&1 || { echo >&2 "[ERROR] go is not installed. Please install go."; exit 1; }
	@command -v upx >/dev/null 2>&1 || { echo >&2 "[ERROR] upx is not installed. Please install upx."; exit 1; }


pack:
	echo "id=Clash_For_Magisk\nname=akashaProxy\nversion="$(shell git rev-parse --short HEAD)"\nversionCode="$(shell git log -1 --format=%ct)"\nauthor=heinu\ndescription=akasha terminal transparent proxy module that supports tproxy and tun and adds many easy-to-use features. Compatible with Magisk/KernelSU">module/module.prop
	cd module && zip -r ../$(NAME).zip *

download-mihomo:
	@[ ! -f module/bin ] && mkdir -p module/bin
	remote_clash_ver=$$(curl --connect-timeout 5 -Ls "https://github.com/MetaCubeX/mihomo/releases/latest/download/version.txt") && \
	curl --connect-timeout 5 -Ls -o module/bin/mihomo-android-arm64-v8.gz \
	"https://github.com/MetaCubeX/mihomo/releases/latest/download/mihomo-android-arm64-v8-$${remote_clash_ver}.gz"
	@echo "mihomo download successfully"

download-dashboard:
	@[ ! -f module/clash/bin ] && mkdir -p module/clash/zashboard
	curl --connect-timeout 5 -Ls -o module/clash/zashboard/dist-cdn-fonts.zip \
	"https://github.com/Zephyruso/zashboard/releases/latest/download/dist-cdn-fonts.zip"
	unzip -o module/clash/zashboard/dist-cdn-fonts.zip -d module/clash/zashboard/
	mv -f module/clash/zashboard/dist/* module/clash/zashboard/
	rm -rf module/clash/zashboard/dist
	rm -rf module/clash/zashboard/dist-cdn-fonts.zip
	@echo "dashboard download successfully"

build-webui:
	cd webui && pnpm i --frozen-lockfile && pnpm build
	mv -f ./webui/out ./module/webroot
	@echo "webui built successfully"

build-tools:
	cd yamlcli && CGO_ENABLED=0 GOOS=android GOARCH=arm64 go build -trimpath -ldflags="-s -w" -buildvcs=false -o ../module/clash/bin/yamlcli
	upx ../module/clash/bin/yamlcli
	@echo "yamlcli built successfully"

clean: 
	rm -rf ./module/module.prop
	rm -rf $(NAME).zip
	rm -rf ./module/webroot
	rm -rf ./module/clash/zashboard
	rm -rf ./module/bin/*
	rm -rf ./module/clash/bin/yamlcli