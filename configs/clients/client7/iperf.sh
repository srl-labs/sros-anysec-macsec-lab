# Copyright 2020 Nokia
# Licensed under the BSD 3-Clause License.
# SPDX-License-Identifier: BSD-3-Clause

# Start iperf3 server in the background
# with 8 parallel tcp streams, each 200 Kbit/s == 1.6Mbit/s
# using ipv6 interfaces

pkill iperf3

iperf3 -c 2002::192:168:1:8 -t 10000 -i 1 -p 5201 -B 2002::192:168:1:7 -P 32 -b 125K -M 1400 &
iperf3 -c 192.168.1.8 -t 10000 -i 1 -p 5201 -B 192.168.1.7 -P 32 -b 125K -M 1400 &

iperf3 -c 2002::192:168:1:8 -t 10000 -i 1 -p 5201 -B 2002::192:168:1:7 -P 32 -b 125K -M 1400 &

iperf3 -c 2002::192:168:1:8 -t 10000 -i 1 -p 5201 -B 2002::192:168:1:7 -P 32 -b 125K -M 1400 &

#Client7
#        - ip -6 address add 2002::192:168:1:7/96 dev eth1
#        - ip address add 192.168.1.7/24 dev eth1
#        - ip address add 192.168.2.7/24 dev eth2
#        - ip address add 192.168.3.7/24 dev eth3
#        - route add default gw 192.168.3.1/24 eth3

#Client 8
#        - ip -6 address add 2002::192:168:1:8/96 dev eth1
#        - ip address add 192.168.1.8/24 dev eth1
#        - ip address add 192.168.2.8/24 dev eth2
#        - ip address add 1.1.1.8/24 dev eth3
#        - route add default gw 1.1.1.1/24 eth3
