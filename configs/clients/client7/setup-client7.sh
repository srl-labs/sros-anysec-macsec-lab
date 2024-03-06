# Client7 setup script

# restart interfaces
ifconfig eth1 down
ifconfig eth2 down
ifconfig eth3 down
ifconfig eth1 up
ifconfig eth2 up
ifconfig eth3 up

# Ipv4 and IPv6 configs
ip -6 address add 2002::192:168:51:7/96 dev eth1
ip address add 192.168.51.7/24 dev eth1
ip address add 192.168.52.7/24 dev eth2
ip address add 192.168.53.7/24 dev eth3
ip route add 192.168.63.0/24 via 192.168.53.1

# Install Automation tools - gnmic, python3 and flask
echo "Install Automation tools - gnmic, python3 and flask "
/bin/bash /config/install_gnmic.sh
