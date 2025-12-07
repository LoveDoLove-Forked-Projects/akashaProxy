[ -z "$version" ] && . /data/clash/clash.config

pid=$(busybox pidof /data/clash/module/ruleconverter/bin/ruleconverter)
if [ -n "${pid}" ]; then
    kill -15 ${pid}
fi