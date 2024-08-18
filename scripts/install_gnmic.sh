#!/bin/bash

### Install gnmic
echo "..."
echo "Start gNMIc Install script execution!"
bash -c "$(curl -sL https://get-gnmic.openconfig.net)" -- -v 0.36.1
echo "END gNMIc Install script execution!"
