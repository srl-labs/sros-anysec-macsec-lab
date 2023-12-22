from flask import Flask, render_template, jsonify
import test
import os
import sys

# Modules
from pygnmi.client import gNMIclient

app = Flask(__name__)

    return "<p>Hello, World From Flask!</p>"
    
@app.route('/', methods=['POST', 'GET'])
def index():
    if request.method == "POST":
        os.system("sh test.sh")

    return render_template('index.html')

    return "<p>Done!</p>"

  
# Example - https://stackoverflow.com/questions/3777301/how-to-call-a-shell-script-from-python-code
#import subprocess
#import shlex
#subprocess.call(shlex.split('./test.sh param1 param2'))
#subprocess.call(["notepad"])
#subprocess.call(['sh', './test.sh']) 


#import os
#import sys
#os.system("sh test.sh")
#os.system(command)