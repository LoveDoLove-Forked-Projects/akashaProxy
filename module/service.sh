until [ "$(getprop sys.boot_completed)" = "1" ]; do
    sleep 2
done

until [ -d "/sdcard/Android" ]; do
    sleep 2
done

chmod -R 770 /data/clash
chown -R root:root /data/clash

. /data/clash/clash.config

if [ ! -d /data/clash/run ]; then
    mkdir -p /data/clash/run
fi
crond -c /data/clash/run


if [ "${self_start}" = "true" ] ; then
    nohup /data/clash/scripts/clash.service -s && /data/clash/scripts/clash.iptables -s & > /data/clash/run/run.logs 2>&1 &
fi
