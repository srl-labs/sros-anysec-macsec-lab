#!/bin/bash

### Install gnmic
echo "..."
echo "Start gNMIc Install script execution!"
bash -c "$(curl -sL https://get-gnmic.openconfig.net)" -- -v 0.36.1
echo "END gNMIc Install script execution!"


### Install other SW
echo "..."
echo "Install Sw!"
apk add python3
apk add python3-venv
wait
echo "END Install Sw !"


### Start Flask and Python installation 
echo "..."
echo "Start Flask and Python installation script execution!" 
cd /flask_app
python3 -m venv venv
source venv/bin/activate
export FLASK_APP=/flask_app/app.py
#apk add pip
apk add gcc
apk add python3-dev
apk add musl-dev
apk add libffi-dev
./venv/bin/python3 -m pip install --upgrade pip
#pip install -r /config/requirements.txt
#pip install pygnmi
#pip install gunicorn 
python -m flask --version
#apt-get install gunicorn
#flask run &
#HTTPS_PROXY="http://10.158.101.16:8080" HTTP_PROXY="http://10.158.101.16:8080" http_proxy="http://10.158.101.16:8080" https_proxy="http://10.158.101.16:8080" 
flask run --host=0.0.0.0 > /dev/pts/0 2>&1 &
sleep 5
echo "END Flask and Python installation script execution!"


### Finish installation script 
echo "..."
ps -ef
sleep 2
echo "Flask Web Server is running at port 8443."
echo "Try http://<SERVER-IP>:8443/"
echo "Script execution Done!"
echo "..."



