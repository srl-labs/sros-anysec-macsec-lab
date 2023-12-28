#!/bin/bash

### Install gnmic
echo ""
echo "Start gNMIc Install script execution!"
bash -c "$(curl -sL https://get-gnmic.openconfig.net)"
#wait
echo "END gNMIc Install script execution!"


#echo ""
#echo "Start set permissions !"
### set Exec permissions for the html dir
#chmod +x /usr/share/nginx/html/
## Set full permissions for script files
## Convert all scripts to unix EOL
#cd /tmp
#chmod -R 777 *
#apk add dos2unix 
#wait
#dos2unix /tmp/*
#wait
#echo "END set permissions !"


echo ""
echo "Install Sw!"
### Install lsof to view and kill Flask PID
#apk add lsof
apk add python3
#apk add python3-venv
wait
echo "END Install Sw !"




echo ""
echo "Start Flask and Python installation Install script execution!"
cd /
mkdir flask_app && cd flask_app
python3 -m venv venv
source venv/bin/activate
export FLASK_APP=/flask_app/app.py
#apk add pip
apk add gcc
apk add python3-dev
apk add musl-dev
apk add libffi-dev
/venv/bin/python3 -m pip install --upgrade pip
pip install pygnmi
pip install flask
#pip install pygnmi
#pip install gunicorn 
#python -m flask --version
#apt-get install gunicorn
#flask run &
flask run --host=0.0.0.0 &
sleep 5
echo "END Flask and Python installation Install script execution!"


echo ""
echo "Flask Web Server is running at port 5000."
#lsof -n | grep flask | grep ESTABLISHED
#wait
echo "Try http://<SERVER-IP>:5000/"
echo "Script execution Done!"
echo ""



