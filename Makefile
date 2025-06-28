NAME=akashaProxy
ghproxy ?= ""

.PHONY: all pack download-mihomo build-webui clean clean-cache-build clean-cache-webui

all: default

default: clean \
	build-webui \
	download-mihomo \
	pack

pack:
	echo "id=Clash_For_Magisk\nname=akashaProxy\nversion="$(shell git rev-parse --short HEAD)"\nversionCode="$(date +%s)"\nauthor=heinu\ndescription=akasha terminal transparent proxy module that supports tproxy and tun and adds many easy-to-use features. Compatible with Magisk/KernelSU">module/module.prop
	cd module && zip -r ../$(NAME).zip *

download-mihomo:
	@[ ! -f module/bin ] && mkdir -p module/bin
	remote_clash_ver=$$(curl --connect-timeout 5 -Ls "$(ghproxy)https://github.com/MetaCubeX/mihomo/releases/latest/download/version.txt") && \
	curl --connect-timeout 5 -Ls -o module/bin/mihomo-android-arm64-v8.gz \
	"$(ghproxy)https://github.com/MetaCubeX/mihomo/releases/latest/download/mihomo-android-arm64-v8-$${remote_clash_ver}.gz"
	@echo "done"

build-webui:
	cd webui && yarn --frozen-lockfile && yarn build
	mv -f ./webui/out ./module/webroot

clean-cache-build:
	rm -rf ./module/module.prop
	rm -rf $(NAME).zip

clean-cache-webui:
	rm -rf ./webui/node_modules
	rm -rf ./module/webroot

clean: clean-cache-build clean-cache-webui
