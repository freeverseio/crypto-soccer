#!/usr/bin/env python3

import os

os.system('mkdir ../nodejs-horizon/contracts')
os.system('cp ../truffle-core/build/contracts/Assets.json ../nodejs-horizon/contracts')
os.system('cp ../truffle-core/build/contracts/Leagues.json ../nodejs-horizon/contracts')
os.system('cp ../truffle-core/build/contracts/PlayerState.json ../nodejs-horizon/contracts')
