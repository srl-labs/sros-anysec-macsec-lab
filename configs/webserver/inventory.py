hosts = [
    {
        "hostname": "pe1",
        "nos": "nokia-sros",
        "port": 57400,
        "username": "admin",
        "password": "admin",
        "has_anysec": True,
        "subscribe": {
            "subscription": [
                {"path": "/configure/port/admin-state", "mode": "on_change"},
                {
                    "path": "/configure/anysec/tunnel-encryption/encryption-group/peer/admin-state",
                    "mode": "on_change",
                },
            ],
            "use_aliases": False,
            "mode": "stream",
            "encoding": "json",
        },
    },
    {
        "hostname": "pe2",
        "nos": "nokia-sros",
        "port": 57400,
        "username": "admin",
        "password": "admin",
        "has_anysec": True,
        "subscribe": {
            "subscription": [
                {"path": "/configure/port/admin-state", "mode": "on_change"},
                {
                    "path": "/configure/anysec/tunnel-encryption/encryption-group/peer/admin-state",
                    "mode": "on_change",
                },
            ],
            "use_aliases": False,
            "mode": "stream",
            "encoding": "json",
        },
    },
    {
        "hostname": "p3",
        "nos": "nokia-sros",
        "port": 57400,
        "username": "admin",
        "password": "admin",
        "has_anysec": False,
        "subscribe": {
            "subscription": [
                {"path": "/configure/port/admin-state", "mode": "on_change"}
            ],
            "use_aliases": False,
            "mode": "stream",
            "encoding": "json",
        },
    },
    {
        "hostname": "p4",
        "nos": "nokia-sros",
        "port": 57400,
        "username": "admin",
        "password": "admin",
        "has_anysec": False,
        "subscribe": {
            "subscription": [
                {"path": "/configure/port/admin-state", "mode": "on_change"}
            ],
            "use_aliases": False,
            "mode": "stream",
            "encoding": "json",
        },
    },
]

links = {
    "top": {"pe1": "1/1/c1/1", "p3": "1/1/c2/1"},
    "bottom": {"pe1": "1/1/c2/1", "p4": "1/1/c2/1"},
}

anysecs = {
    "vll": {
        "pe1": {
            "group_name": "EG_VLL-1001",
            "peer": "10.0.0.21",
        },
        "pe2": {
            "group_name": "EG_VLL-1001",
            "peer": "10.0.0.11",
        },
    },
    "vpls": {
        "pe1": {
            "group_name": "EG_VPLS-1002",
            "peer": "10.0.0.22",
        },
        "pe2": {
            "group_name": "EG_VPLS-1002",
            "peer": "10.0.0.12",
        },
    },
    "vprn": {
        "pe1": {
            "group_name": "EG_VPRN-1003",
            "peer": "10.0.0.2",
        },
        "pe2": {
            "group_name": "EG_VPRN-1003",
            "peer": "10.0.0.1",
        },
    },
}

icmp_types = {"vll": "192.168.1.8", "vpls": "192.168.2.8", "vprn": "1.1.1.8"}
