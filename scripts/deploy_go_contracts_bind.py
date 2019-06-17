#!/usr/bin/env python3

import json
import os

with open('../truffle-core/build/contracts/Assets.json', 'r') as fp:
    abi = json.load(fp)['abi']

with open('./Assets.abi', 'w', encoding='utf-8') as outfile:
    json.dump(abi, outfile, ensure_ascii=False, indent=2)

os.system('abigen --abi ./Assets.abi --pkg assets -out ../go-synchronizer/contracts/assets/assets.go')

with open('../truffle-core/build/contracts/TeamState.json', 'r') as fp:
    abi = json.load(fp)['abi']

with open('./States.abi', 'w', encoding='utf-8') as outfile:
    json.dump(abi, outfile, ensure_ascii=False, indent=2)

os.system('abigen --abi ./States.abi --pkg states -out ../go-synchronizer/contracts/states/states.go')

with open('../truffle-core/build/contracts/Leagues.json', 'r') as fp:
    abi = json.load(fp)['abi']

with open('./Leagues.abi', 'w', encoding='utf-8') as outfile:
    json.dump(abi, outfile, ensure_ascii=False, indent=2)

os.system('abigen --abi ./Leagues.abi --pkg leagues -out ../go-synchronizer/contracts/leagues/leagues.go')