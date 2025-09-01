until [ "$(getprop sys.boot_completed)" = "1" ]; do
    sleep 2
done

until [ -d "/sdcard/Android" ]; do
    sleep 2
done

chmod -R 770 /data/clash
chown -R root:root /data/clash

. /data/clash/clash.config

if [ ! -d ${Clash_run_path} ]; then
    mkdir -p ${Clash_run_path}
fi
crond -c ${Clash_run_path}


if [ ${self_start} == "true" ] ; then
    nohup /data/clash/scripts/clash.service -s && /data/clash/scripts/clash.iptables -s
fi
