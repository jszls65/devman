#!/bin/bash
bin_path=$(dirname `readlink -f "$0"`)
source /etc/profile
cd ${bin_path}

nohup ./dev-utils >/dev/null &