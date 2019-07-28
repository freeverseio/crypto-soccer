#!/usr/bin/env python3

import json
import os

with open('../truffle-core/build/contracts/Assets.json', 'r') as fp:
    contract = json.load(fp)
with open('./Assets.abi', 'w', encoding='utf-8') as outfile:
    json.dump(contract['abi'], outfile, ensure_ascii=False, indent=2)
with open('./Assets.bin', 'w', encoding='utf-8') as outfile:
    outfile.write(contract['bytecode'])

os.system('mkdir -p ../go-synchronizer/contracts/assets')
os.system('abigen --abi ./Assets.abi --bin ./Assets.bin --pkg assets -out ../go-synchronizer/contracts/assets/assets.go')

with open('../truffle-core/build/contracts/TeamState.json', 'r') as fp:
    contract = json.load(fp)
with open('./States.abi', 'w', encoding='utf-8') as outfile:
    json.dump(contract['abi'], outfile, ensure_ascii=False, indent=2)
with open('./States.bin', 'w', encoding='utf-8') as outfile:
    outfile.write(contract['bytecode'])

os.system('mkdir -p ../go-synchronizer/contracts/states')
os.system('abigen --abi ./States.abi --bin ./States.bin --pkg states -out ../go-synchronizer/contracts/states/states.go')

with open('../truffle-core/build/contracts/Leagues.json', 'r') as fp:
    contract = json.load(fp)
with open('./Leagues.abi', 'w', encoding='utf-8') as outfile:
    json.dump(contract['abi'], outfile, ensure_ascii=False, indent=2)
with open('./Leagues.bin', 'w', encoding='utf-8') as outfile:
    outfile.write(contract['bytecode'])

os.system('mkdir -p ../go-synchronizer/contracts/leagues')
os.system('abigen --abi ./Leagues.abi --bin ./Leagues.bin --pkg leagues -out ../go-synchronizer/contracts/leagues/leagues.go')

with open('../truffle-core/build/contracts/Engine.json', 'r') as fp:
    contract = json.load(fp)
with open('./Engine.abi', 'w', encoding='utf-8') as outfile:
    json.dump(contract['abi'], outfile, ensure_ascii=False, indent=2)
with open('./Engine.bin', 'w', encoding='utf-8') as outfile:
    outfile.write(contract['bytecode'])

os.system('mkdir -p ../go-synchronizer/contracts/engine')
os.system('abigen --abi ./Engine.abi --bin ./Engine.bin --pkg engine -out ../go-synchronizer/contracts/engine/engine.go')

os.system('rm -rf ./*abi ./*bin')