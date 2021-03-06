#!/bin/sh

## 自动定时拉取仓库更新以及题库更新
# crontab -e
# 每天0点执行updatestart.sh：`0 0 * * * /bin/sh /path/to/updatestart.sh`

# answers_download_url="http://81.68.160.189:35247/download"
answers_file_name="answer.json"

current_path=$(cd `dirname $0`; pwd)
screen_name="lgb"

cd $current_path
echo $(date)
echo "running path: $current_path"

# pull latest git code
echo "pull latest git code..."
git reset --hard
git pull

# make updatestart.sh and release/HttpMonitor executable
echo "make updatestart.sh and release/HttpMonitor executable..."
chmod +x $current_path/updatestart.sh
chmod +x $current_path/release/HttpMonitor_linux

# update answer.json
# echo "update answer.json..."
# curl -s "$answers_download_url" > "$current_path/$answers_file_name"

# kill lgb proxy
echo "kill lgb proxy..."
screen -S "$screen_name" -X quit

# copy answers.json to release
# cp "$current_path/$answers_file_name" "$current_path/release/$answers_file_name"

# start lgb proxy
screen -dmS "$screen_name"
screen -x -S "$screen_name" -p 0 -X stuff "chmod +x $current_path/release/HttpMonitor_linux && $current_path/release/HttpMonitor_linux \n"
echo "lgb proxy started"