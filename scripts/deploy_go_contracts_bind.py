#!/usr/bin/env python3

import json
import os

with open('../truffle-core/build/contracts/Assets.json', 'r') as fp:
    contract = json.load(fp)
with open('./Assets.abi', 'w', encoding='utf-8') as outfile:
    json.dump(contract['abi'], outfile, ensure_ascii=False, indent=2)
with open('./Assets.bin', 'w', encoding='utf-8') as outfile:
    outfile.write(contract['bytecode'])

os.system('abigen --abi ./Assets.abi --bin ./Assets.bin --pkg assets -out ../go-synchronizer/contracts/assets/assets.go')

with open('../truffle-core/build/contracts/TeamState.json', 'r') as fp:
    contract = json.load(fp)
with open('./States.abi', 'w', encoding='utf-8') as outfile:
    json.dump(contract['abi'], outfile, ensure_ascii=False, indent=2)
with open('./States.bin', 'w', encoding='utf-8') as outfile:
    outfile.write(contract['bytecode'])

os.system('abigen --abi ./States.abi --bin ./States.bin --pkg states -out ../go-synchronizer/contracts/states/states.go')

with open('../truffle-core/build/contracts/Leagues.json', 'r') as fp:
    contract = json.load(fp)
with open('./Leagues.abi', 'w', encoding='utf-8') as outfile:
    json.dump(contract['abi'], outfile, ensure_ascii=False, indent=2)
with open('./Leagues.bin', 'w', encoding='utf-8') as outfile:
    outfile.write(contract['bytecode'])

os.system('abigen --abi ./Leagues.abi --bin ./Leagues.bin --pkg leagues -out ../go-synchronizer/contracts/leagues/leagues.go')

os.system('rm -rf ./*abi ./*bin')