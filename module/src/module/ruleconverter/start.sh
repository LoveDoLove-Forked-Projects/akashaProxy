[ -z "$version" ] && . /data/clash/clash.config

if [ "$ruleconverter" != "true" ]; then
    exit 0
fi

nohup /data/clash/module/ruleconverter/bin/ruleconverter -port ${ruleconverter_port} 2>&1 &