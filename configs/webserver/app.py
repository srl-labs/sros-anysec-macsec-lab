from flask import Flask, render_template
import subprocess

app = Flask(__name__)

@app.route('/')
def index():
    return render_template('index.html')


@app.route('/execute_python')
def execute_python_script():
    try:
        result = subprocess.check_output(['python', '/flask_app/hello.py'], stderr=subprocess.STDOUT, text=True)
        return f"Python Script Execution Result: {result}"
    except subprocess.CalledProcessError as e:
        return f"Error executing Python script: {e.output}"

@app.route('/execute_shell')
def execute_shell_script():
    try:
        result = subprocess.check_output(['bash', '/flask_app/hello.sh'], stderr=subprocess.STDOUT, text=True)
        return f"Shell Script Execution Result: {result}"
    except subprocess.CalledProcessError as e:
        return f"Error executing Shell script: {e.output}"


@app.route('/execute_gnmic/link_enable_top')
def execute_gnmic_link_enable_top():
    try:
        # Execute gnmic command directly for enabling links
        result1 = subprocess.run('gnmic -a pe1:57400 -u admin -p admin --insecure set --update-path /configure/port[port-id=1/1/c1/1]/admin-state --update-value enable', shell=True, capture_output=True, text=True)
        
        result2 = subprocess.run('gnmic -a p3:57400 -u admin -p admin --insecure set --update-path /configure/port[port-id=1/1/c2/1]/admin-state --update-value enable', shell=True, capture_output=True, text=True)

        return f"gNMIc Commands Execution Result (Enable Links): {result1.stdout}\n{result2.stdout}"
    except subprocess.CalledProcessError as e:
        return f"Error executing gNMIc commands (Enable Links): {e.stderr}"


@app.route('/execute_gnmic/link_disable_top')
def execute_gnmic_link_disable_top():
    try:
        # Execute gnmic command directly for disabling links
        result1 = subprocess.run('gnmic -a pe1:57400 -u admin -p admin --insecure set --update-path /configure/port[port-id=1/1/c1/1]/admin-state --update-value disable', shell=True, capture_output=True, text=True)
        
        result2 = subprocess.run('gnmic -a p3:57400 -u admin -p admin --insecure set --update-path /configure/port[port-id=1/1/c2/1]/admin-state --update-value disable', shell=True, capture_output=True, text=True)

        return f"gNMIc Commands Execution Result (Disable Links): {result1.stdout}\n{result2.stdout}"
    except subprocess.CalledProcessError as e:
        return f"Error executing gNMIc commands (Disable Links): {e.stderr}"


if __name__ == '__main__':
    app.run(host='0.0.0.0', debug=True)