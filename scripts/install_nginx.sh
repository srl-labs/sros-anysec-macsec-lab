#!/bin/bash

### Install NGINX server - hellt/nginx-uwsgi-flask-alpine-docker:py3


echo ""
echo "Start gNMIc Install script execution!"
bash -c "$(curl -sL https://get-gnmic.openconfig.net)"
wait
echo "END gNMIc Install script execution!"


echo ""
echo "Start set permissions !"
### set Exec permissions for the html dir
## Set full permissions for script files
cd /app
chmod +x /app/
## Convert all scripts to unix EOL
apt install dos2unix 
wait
dos2unix /app/*
wait
### Install lsof to view and kill Flask PID
apt install lsof
wait 
# lsof -n -i:5000
# kill PID
echo "END set permissions !"


#echo ""
#echo "Start Flask and Python installation Install script execution!"
#apt update
#wait
#apt install python3
#wait
#Y
#python3 -V
#apt install pip
#wait
#apt-get install python3-venv
#wait
#Y
#mkdir flask_app && cd flask_app
#python3 -m venv venv
#wait
#source venv/bin/activate
#wait
#pip install flask
#wait
#pip install pygnmi
#wait
#pip install gunicorn
#python -m flask --version
#apt-get install gunicorn
#export FLASK_APP=/usr/share/nginx/html/index.py
#flask run &
#flask run --host=0.0.0.0 &
#sleep 5
#echo "END Flask and Python installation Install script execution!"


echo ""
echo "Flask Web Server is running at port:"
#lsof -n -i:5000
lsof -n | grep flask | grep ESTABLISHED
echo "Try http://<SERVER-IP>:<PORT>/"
echo "Script execution Done!"
echo ""



