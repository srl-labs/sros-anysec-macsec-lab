import json
import re
import signal
import subprocess
import threading
import traceback 
import sys

from time import sleep
from flask import Flask, jsonify, render_template, request
from inventory import anysecs, hosts, icmp_types, links
from pygnmi.client import gNMIclient, gNMIException
import grpc, io
from pprint import pprint


def threaded(fn):
    def wrapper(*args, **kwargs):
        thread = threading.Thread(target=fn, args=args, kwargs=kwargs)
        thread.start()
        return thread

    return wrapper


class Telemetry:
    def __init__(self):
        threading.Thread.__init__(self)
        self.timeout = 10
        self.routers = {}
        self.links = {"top": "enabled", "bottom": "enabled"}
        self.anysecs = {"vll": "enabled", "vpls": "enabled", "vprn": "enabled"}
        self.icmp_request = {"vll": None, "vpls": None, "vprn": None}
        self.icmp_status = {"vll": "disabled", "vpls": "disabled", "vprn": "disabled"}

    @threaded
    def monitor_links(self):
        while True:
            for link, routers_data in links.items():
                link_status = False
                port_status = [False, False]
                index = 0
                for router, port in routers_data.items():
                    if router in self.routers:
                        if port in self.routers[router]["ports"]:
                            if (
                                self.routers[router]["ports"][port]["admin-state"]
                                == "enable"
                            ):
                                port_status[index] = True
                    index += 1
                link_status = bool(port_status[0] * port_status[1])
                if link_status:
                    self.links[link] = "enabled"
                else:
                    self.links[link] = "disabled"

    @threaded
    def monitor_anysecs(self):
        while True:
            for anysec, routers_data in anysecs.items():
                link_status = False
                anysec_status = [False, False]
                index = 0
                for router, anysec_data in routers_data.items():
                    anysec_group = anysec_data["group_name"]
                    if router in self.routers:
                        if "anysec_group" in self.routers[router]:
                            if anysec_group in self.routers[router]["anysec_group"]:
                                if (
                                    self.routers[router]["anysec_group"][anysec_group][
                                        "admin-state"
                                    ]
                                    == "enable"
                                ):
                                    anysec_status[index] = True
                    index += 1
                link_status = bool(anysec_status[0] * anysec_status[1])
                if link_status:
                    self.anysecs[anysec] = "enabled"
                else:
                    self.anysecs[anysec] = "disabled"

    @threaded
    def start_icmp_trafic_red(self, icmp_type, icmp_size, icmp_interval):
        self.icmp_request["vll"] = True
        if icmp_type in icmp_types:
            icmp_ip = icmp_types[icmp_type]
            p = subprocess.Popen(
                "ping -s "
                + icmp_size
                + " -i "
                + icmp_interval
                + " "
                + icmp_ip
                + " 1>/dev/null",
                stdout=subprocess.PIPE,
                shell=True,
            )
            while p.poll() is None:
                self.icmp_status[icmp_type] = "enabled"
                if self.icmp_request[icmp_type] is False:
                    p.send_signal(signal.SIGINT)
                    self.icmp_request[icmp_type] = None
        self.icmp_status[icmp_type] = "disabled"

    @threaded
    def sub_telemetry(self, host_entry):
        # print(host_entry["hostname"])
        subscribe = host_entry["subscribe"]
        # print(subscribe)
        gc = gNMIclient(
                    target=(host_entry["hostname"], host_entry["port"]),
                    username=host_entry["username"],
                    password=host_entry["password"],
                    insecure=True,
                )
        while True:
            try:
                gc.connect()
                telemetry_stream = gc.subscribe_stream(subscribe=subscribe)
                for telemetry_entry in telemetry_stream:
                    # telemetry_entry_json = json.dumps(telemetry_entry)
                    # if telemetry_entry["update"]["update"][0]["path"] == "admin-state":
                    if "update" in telemetry_entry:
                        if "prefix" in telemetry_entry["update"]:
                            if "port-id" in telemetry_entry["update"]["prefix"]:
                                try:
                                    port = re.findall(
                                        "port-id=(.*)",
                                        re.findall("\[(.*?)\]", telemetry_entry["update"]["prefix"])[0],
                                    )[0]
                                except IndexError as e:
                                    raise grpc.FutureTimeoutError()
                                if host_entry["hostname"] not in self.routers:
                                    self.routers[host_entry["hostname"]] = {}
                                if "ports" not in self.routers[host_entry["hostname"]]:
                                    self.routers[host_entry["hostname"]]["ports"] = {}
                                if port not in self.routers[host_entry["hostname"]]["ports"]:
                                    self.routers[host_entry["hostname"]]["ports"][port] = {}
                                self.routers[host_entry["hostname"]]["ports"][port][
                                    "admin-state"
                                ] = telemetry_entry["update"]["update"][0]["val"]
                            if "anysec" in telemetry_entry["update"]["prefix"]:
                                try:
                                    anysec_group = re.findall(
                                        "group-name=(.*)",
                                        re.findall("\[(.*?)\]", telemetry_entry["update"]["prefix"])[0],
                                    )[0]
                                except IndexError as e:
                                    raise grpc.FutureTimeoutError()
                                if host_entry["hostname"] not in self.routers:
                                    self.routers[host_entry["hostname"]] = {}
                                if "anysec_group" not in self.routers[host_entry["hostname"]]:
                                    self.routers[host_entry["hostname"]]["anysec_group"] = {}
                                if (
                                    anysec_group
                                    not in self.routers[host_entry["hostname"]]["anysec_group"]
                                ):
                                    self.routers[host_entry["hostname"]]["anysec_group"][
                                        anysec_group
                                    ] = {}
                                self.routers[host_entry["hostname"]]["anysec_group"][anysec_group][
                                    "admin-state"
                                ] = telemetry_entry["update"]["update"][0]["val"]
                        else:
                            raise grpc.FutureTimeoutError()
                    else:
                        raise grpc.FutureTimeoutError()
        
                    # print(self.routers)
                    # print(host_entry["hostname"]+' - '+telemetry_entry_json)
                print("[Flask] - [Debug] - sub_telemetry - before break")
                break
            except gNMIException as e:
                print(host_entry["hostname"])
                #traceback.print_exception(*sys.exc_info())
                print(e)
            except grpc.FutureTimeoutError as e:
                print("[Flask] - sub_telemetry - unable to connect to the host "+str(host_entry["hostname"])+" will try agrain in "+str(self.timeout)+" seconds")
                #traceback.print_exception(*sys.exc_info())
                #raise
                try:
                    gc.close()
                except grpc._channel._MultiThreadedRendezvous as e2:
                    print(e2)
                sleep(self.timeout)
                #self.sub_telemetry(host_entry)
                #print(e)

    @threaded
    def update_port_status(self, host_entry):
        paths = ["/configure/port/admin-state"]
        gc = gNMIclient(
                    target=(host_entry["hostname"], host_entry["port"]),
                    username=host_entry["username"],
                    password=host_entry["password"],
                    insecure=True,
                )
        while True:
            try:
                gc.connect()
                result = gc.get(path=paths, encoding="json")
                if "notification" in result:
                    if "update" in result["notification"][0]:
                        for port_result in result["notification"][0]["update"]:
                            try:
                                port = re.findall(
                                    "port-id=(.*)", re.findall("\[(.*?)\]", port_result["path"])[0]
                                )[0]
                            except IndexError as e:
                                raise grpc.FutureTimeoutError()
                            if host_entry["hostname"] not in self.routers:
                                self.routers[host_entry["hostname"]] = {}
                            if "ports" not in self.routers[host_entry["hostname"]]:
                                self.routers[host_entry["hostname"]]["ports"] = {}
                            if port not in self.routers[host_entry["hostname"]]["ports"]:
                                self.routers[host_entry["hostname"]]["ports"][port] = {}
                            self.routers[host_entry["hostname"]]["ports"][port][
                                "admin-state"
                            ] = port_result["val"]
                        print("[Flask] - [Debug] - update_port_status - before break")
                        break
                    else:
                        raise grpc.FutureTimeoutError()
                else:
                    raise grpc.FutureTimeoutError()
            except gNMIException as e:
                print(host_entry["hostname"])
                #traceback.print_exception(*sys.exc_info())
                print(e)
            except grpc.FutureTimeoutError as e:
                print("[Flask] - update_port_status - unable to connect to the host "+str(host_entry["hostname"])+" will try agrain in "+str(self.timeout)+" seconds")
                #traceback.print_exception(*sys.exc_info())
                try:
                    gc.close()
                except grpc._channel._MultiThreadedRendezvous as e2:
                    print(e2)
                sleep(self.timeout)
                #self.update_port_status(host_entry)
                #print(e)

    @threaded
    def update_anysec_status(self, host_entry):
        paths = [
            "/configure/anysec/tunnel-encryption/encryption-group/peer/admin-state"
        ]
        if host_entry["has_anysec"]:
            gc = gNMIclient(
                        target=(host_entry["hostname"], host_entry["port"]),
                        username=host_entry["username"],
                        password=host_entry["password"],
                        insecure=True,
                    )
            while True:
                try:
                    gc.connect()
                    result = gc.get(path=paths, encoding="json")
                    if "notification" in result:
                        if "update" in result["notification"][0]:
                            for anysec_result in result["notification"][0]["update"]:
                                try:
                                    anysec_group = re.findall(
                                        "group-name=(.*)",
                                        re.findall("\[(.*?)\]", anysec_result["path"])[0],
                                    )[0]
                                except IndexError as e:
                                    raise grpc.FutureTimeoutError()
                                if host_entry["hostname"] not in self.routers:
                                    self.routers[host_entry["hostname"]] = {}
                                if "anysec_group" not in self.routers[host_entry["hostname"]]:
                                    self.routers[host_entry["hostname"]]["anysec_group"] = {}
                                if (
                                    anysec_group
                                    not in self.routers[host_entry["hostname"]]["anysec_group"]
                                ):
                                    self.routers[host_entry["hostname"]]["anysec_group"][
                                        anysec_group
                                    ] = {}
                                self.routers[host_entry["hostname"]]["anysec_group"][
                                    anysec_group
                                ]["admin-state"] = anysec_result["val"]
                            print("[Flask] - [Debug] - update_anysec_status - before break")
                            break
                        else:
                            raise grpc.FutureTimeoutError()
                    else:
                        raise grpc.FutureTimeoutError()
                except gNMIException as e:
                    print(host_entry["hostname"])
                    #traceback.print_exception(*sys.exc_info())
                    print(e)
                except grpc.FutureTimeoutError as e:
                    print("[Flask] - update_anysec_status - unable to connect to the host "+str(host_entry["hostname"])+" will try agrain in "+str(self.timeout)+" seconds")
                    #traceback.print_exception(*sys.exc_info())
                    try:
                        gc.close()
                    except grpc._channel._MultiThreadedRendezvous as e2:
                        print(e2)
                    sleep(self.timeout)
                    #self.update_anysec_status(host_entry)
                    #print(e)

    def run(self):
        for host_entry in hosts:
            self.sub_telemetry(host_entry)
            self.update_port_status(host_entry)
            self.update_anysec_status(host_entry)
        self.monitor_links()
        self.monitor_anysecs()


app = Flask(__name__)
telemetry = Telemetry()


@app.route("/")
def index():
    host, port = request.environ["HTTP_HOST"].split(":")
    return render_template("index.html", host=host)


@app.route("/execute_gnmic/link_toggle/<link_name>")
def execute_gnmic_link_toggle(link_name):
    try:
        # Execute gnmic command directly for enabling links
        if link_name not in telemetry.links:
            return (
                '{status:"failed", message:"The link '
                + link_name
                + ' is not configured"}'
            )
        action = ""
        if telemetry.links[link_name] == "disabled":
            action = "enable"
        else:
            action = "disable"
        thislink = links[link_name]
        results = {}
        for router, port in thislink.items():
            results[router] = subprocess.run(
                "gnmic -a "
                + router
                + ":57400 -u admin -p admin --insecure set --update-path /configure/port[port-id="
                + port
                + "]/admin-state --update-value "
                + action,
                shell=True,
                capture_output=True,
                text=True,
            ).stdout
        return jsonify(results)
        # return f"gNMIc Commands Execution Result (Enable Links): {result1.stdout}\n{result2.stdout}"
    except subprocess.CalledProcessError as e:
        # return f"Error executing gNMIc commands (Enable Links): {e.stderr}"
        result = '{status:"failed", message:"' + e.stderr + '"}'
        return result


@app.route("/execute_gnmic/anysec_toggle/<anysec_name>")
def execute_gnmic_anysec_toggle(anysec_name):
    try:
        # Execute gnmic command directly for enabling links
        if anysec_name not in telemetry.anysecs:
            return (
                '{status:"failed", message:"The anysec '
                + anysec_name
                + ' is not configured"}'
            )
        action = ""
        if telemetry.anysecs[anysec_name] == "disabled":
            action = "enable"
        else:
            action = "disable"
        thisanysec = anysecs[anysec_name]
        results = {}
        for router, anysec_data in thisanysec.items():
            anysec_group = anysec_data["group_name"]
            anysec_peer = anysec_data["peer"]
            results[router] = subprocess.run(
                "gnmic -a "
                + router
                + ":57400 -u admin -p admin --insecure set --update-path /configure/anysec/tunnel-encryption/encryption-group[group-name="
                + anysec_group
                + "]/peer[peer-ip-address="
                + anysec_peer
                + "]/admin-state --update-value "
                + action,
                shell=True,
                capture_output=True,
                text=True,
            ).stdout
        return jsonify(results)
        # return f"gNMIc Commands Execution Result (Enable Links): {result1.stdout}\n{result2.stdout}"
    except subprocess.CalledProcessError as e:
        # return f"Error executing gNMIc commands (Enable Links): {e.stderr}"
        result = '{status:"failed", message:"' + e.stderr + '"}'
        return result


@app.route("/icmp_toggle/<icmp_type>/<icmp_size>/<icmp_interval>")
def icmp_toggle(icmp_type, icmp_size, icmp_interval):
    try:
        # Execute gnmic command directly for enabling links
        if icmp_type not in telemetry.icmp_status:
            return (
                '{status:"failed", message:"The '
                + icmp_type
                + ' icmp is not configured"}'
            )
        action = ""
        if telemetry.icmp_status[icmp_type] == "disabled":
            telemetry.start_icmp_trafic_red(icmp_type, icmp_size, icmp_interval)
            action = "enable"
        else:
            telemetry.icmp_request[icmp_type] = False
            action = "disable"

        return jsonify(action)
    except subprocess.CalledProcessError as e:
        # return f"Error executing gNMIc commands (Enable Links): {e.stderr}"
        result = '{status:"failed", message:"' + e.stderr + '"}'
        return result


@app.route("/execute_gnmic/get_port_status/<element>/<card>/<mda>/<conector>/<port>")
def execute_gnmic_get_port_status(element, card, mda, conector, port="1"):
    try:
        portID = card + "/" + mda + "/" + conector + "/" + port
        request = subprocess.run(
            "gnmic -a "
            + element
            + ":57400 -u admin -p admin --insecure get --path /configure/port[port-id="
            + portID
            + "]/admin-state",
            shell=True,
            capture_output=True,
            text=True,
        )
        result = request.stdout
        return result
    except subprocess.CalledProcessError as e:
        result = '{status:"failed", message:"' + e.stderr + '"}'
        return result


@app.route("/get_routers/data")
def get_routers_data():
    return jsonify(telemetry.routers)


@app.route("/get_links/data")
def get_links_data():
    return jsonify(telemetry.links)


@app.route("/get_anysecs/data")
def get_anysecs_data():
    return jsonify(telemetry.anysecs)


@app.route("/get_icmp_status")
def get_icmp_status():
    return jsonify(telemetry.icmp_status)


@app.route("/get_server_info")
def get_server_info():
    return jsonify(request.environ)


if __name__ == "app":
    # print(__name__)
    telemetry.run()
if __name__ == "__main__":
    app.run(host="0.0.0.0", debug=True)
