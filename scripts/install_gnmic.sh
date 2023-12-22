#!/bin/bash

### Install gnmic
echo ""
echo "Start gNMIc Install script execution!"
bash -c "$(curl -sL https://get-gnmic.openconfig.net)"
wait
echo "END gNMIc Install script execution!"


#echo ""
#echo "Start set permissions !"
### set Exec permissions for the html dir
#chmod +x /usr/share/nginx/html/
## Set full permissions for script files
## Convert all scripts to unix EOL
#cd /tmp
#chmod -R 777 *
#apt install dos2unix 
#wait
#dos2unix /tmp/*
#wait
#echo "END set permissions !"


echo ""
echo "Install Sw!"
### Install lsof to view and kill Flask PID
apt update -y
apt install lsof -y
apt install python3 -y
apt install python3-venv -y
wait
echo "END Install Sw !"




echo ""
echo "Start Flask and Python installation Install script execution!"
cd /
mkdir flask_app && cd flask_app
python3 -m venv venv
source venv/bin/activate
export FLASK_APP=/flask_app/app.py
#apt install pip
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
wait
echo "Try http://<SERVER-IP>:5000/"
echo "Script execution Done!"
echo ""



