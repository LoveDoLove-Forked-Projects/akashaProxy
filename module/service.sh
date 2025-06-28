#!/system/bin/sh
. /data/clash/clash.config



if [ ! -d ${Clash_run_path} ]; then
    mkdir -p ${Clash_run_path}
fi
crond -c ${Clash_run_path}
chmod -R 770 ${Clash_data_dir}

if [ ${self_start} == "true" ] ; then
    nohup /data/clash/scripts/clash.service -s && /data/clash/scripts/clash.iptables -s
fi
