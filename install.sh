#!/bin/bash

kernel=`uname -s|tr A-Z a-z`
arch=`uname -m`

echo "input your password:"

sudo curl https://github.com/lingdor/ggrep/releases/download/v0.0.2/ggrep_v0.0.1_${kernel}_${arch} -o /usr/local/bin/ggrep


