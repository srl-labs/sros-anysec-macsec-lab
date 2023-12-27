# Copyright 2023 Nokia
# Licensed under the BSD 3-Clause License.
# SPDX-License-Identifier: BSD-3-Clause

# Generated traffic Approx 3.7Mbit/s
# using ipv6 interfaces



#### VLL

# Must be executed from host CLI directly
#ping -s 1450 -c 200000 -i 0.04 2002::192:168:1:8
ping -s 1450 -c 600 -i 0.01 192.168.1.8

# Must be executed from CLAB/Hypervisor server
#ip netns exec client7 ping -s 1450 -c 200000 -i 0.04 192.168.1.8 

# Requires user password input => Cannot be used for the webserver
#ssh root@10.82.182.179 "ip netns exec client7 ping -s 1450 -c 200000 -i 0.04 192.168.1.8"



#### VPLS
#ping -s 1450 -c 200000 -i 0.04 192.168.2.8





### VPRN
#ping -s 1450 -c 200000 -i 0.04 1.1.1.8





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
