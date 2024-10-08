# Ipv4 and IPv6 configs
ip -6 address add 2002::192:168:51:8/96 dev eth1
ip address add 192.168.51.8/24 dev eth1
ip address add 192.168.52.8/24 dev eth2
ip address add 192.168.63.8/24 dev eth3
ip route add 192.168.53.0/24 via 192.168.63.1