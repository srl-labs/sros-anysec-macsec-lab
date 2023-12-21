#!/bin/bash

### Install gnmic
bash -c "$(curl -sL https://get-gnmic.openconfig.net)"

### set Exec permissions for the html dir
chmod +x /usr/share/nginx/html/


## Set full permissions for script files
cd /opt/lampp/htdocs/html
chmod -R 777 *

## Convert all scripts to unix EOL
apt install dos2unix  
dos2unix /tmp/*


echo "END gNMIc Install script execution!"





apt update
apt install python3
Y
python3 -V
apt-get install python3-venv
Y
mkdir flask_app && cd flask_app
python3 -m venv venv
source venv/bin/activate
pip install flask
#pip install gunicorn
python -m flask --version
#apt-get install gunicorn
export FLASK_APP=/usr/share/nginx/html/hello.py
flask run &


echo "END Flask and Python installation Install script execution!"