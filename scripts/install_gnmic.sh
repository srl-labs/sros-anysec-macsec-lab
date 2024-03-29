#!/bin/bash

### Install gnmic
echo "..."
echo "Start gNMIc Install script execution!"
bash -c "$(curl -sL https://get-gnmic.openconfig.net)"
echo "END gNMIc Install script execution!"


### Install other SW
echo "..."
echo "Install Sw!"
apk add python3
#apk add python3-venv
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
pip install -r /config/requirements.txt
#pip install pygnmi
#pip install gunicorn 
#python -m flask --version
#apt-get install gunicorn
#flask run &
HTTPS_PROXY= HTTP_PROXY= http_proxy= https_proxy= flask run --host=0.0.0.0 > /dev/pts/0 2>&1 &
sleep 5
echo "END Flask and Python installation script execution!"


### Finish installation script 
echo "..."
echo "Flask Web Server is running at port 35000."
#ps -ef
#wait
echo "Try http://<SERVER-IP>:35000/"
echo "Script execution Done!"
echo "..."



