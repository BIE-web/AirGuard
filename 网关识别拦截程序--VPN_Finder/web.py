import os
import subprocess
from flask import Flask, render_template
import json
from signal import SIGKILL
from multiprocessing import Process
import atexit

import vpn_finder

app = Flask(__name__) 
websocketd_process = 0
gotranal_process = 0
vpn_finder_process = 0
with open('config.json', 'r') as f:
    config = json.loads(f.read())

@app.route("/") 
def index(): 
    ws = "ws://" + config["ip"] + ":" + config["websocketd_port"]
    return  render_template("html/log.html",server = ws)

@app.route("/start", methods=['GET'])
def start():
    gotranal_process = subprocess.Popen(("ssh root@" +config['router_ip'] + " /root/web/gotranal on -I " + config['iface'] + ' -s '+ config['ip']+':' + config['vpn_finder_port']).split(" "))
    vpn_finder_process = Process(target=vpn_finder.run, args= [config['ip'],int(config['vpn_finder_port'])])
    vpn_finder_process.start()
    return str({gotranal_process.pid,vpn_finder_process.pid})

@app.route("/stop", methods=['GET'])
def stop():
    try:
        os.system("echo "" >  log")
        gotranal_process.kill()
        vpn_finder_process.terminate()
    except:
        pass
    return 1

@atexit.register
def exit_handler():
    try:
        os.system("echo "" >  log")
        websocketd_process.kill()
        gotranal_process.kill()
        vpn_finder_process.terminate()
    except:
        pass


if __name__ == "__main__":
    websocketd_process = subprocess.Popen(("./websocketd --port " + config['websocketd_port'] + " tail -f ./log").split(" "))
    app.run(host='0.0.0.0',port=int(config['flask_port'])) #运行app
