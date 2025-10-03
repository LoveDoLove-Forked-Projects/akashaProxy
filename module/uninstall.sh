data_dir="/data/clash"

rm_data() {
    rm -rf ${data_dir}
    rm -rf ${data_dir}.old
}

rm_data