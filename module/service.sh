#!/system/bin/sh
. /data/clash/clash.config

crond -c ${Clash_run_path}

if [ ${self_start} == "true" ] ; then
    nohup /data/clash/scripts/clash.service -s && /data/clash/scripts/clash.iptables -s
fi
