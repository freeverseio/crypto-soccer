#!/usr/bin/env python

import json
import os,sys
import subprocess

mydir = os.path.dirname(os.path.realpath(__file__))
parentdir = os.path.dirname(mydir)
abigen=os.path.join(os.getenv("HOME"), 'go', 'bin', 'abigen')

def run_cmd(cmd, working_dir = os.getcwd()):
    try:
        print 'running ' + ' '.join(cmd)
        output = subprocess.Popen(cmd, cwd=working_dir, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    except subprocess.CalledProcessError as e:
        print '"' + ' '.join(e.cmd) + '" failed with return code ' + str(e.returncode)
        sys.exit(e.returncode)
    return output.communicate()

if not os.path.exists(abigen):
    print 'Could not find abigen in ' + abigen
    sys.exit(-1)

with open(os.path.join(parentdir, 'truffle-core', 'build', 'contracts', 'Assets.json'), 'r') as fp:
    contract = json.load(fp)
with open(os.path.join(mydir, 'Assets.abi'), 'w') as outfile:
    json.dump(contract['abi'], outfile, ensure_ascii=False, indent=2)
with open(os.path.join(mydir, 'Assets.bin'), 'w') as outfile:
    outfile.write(contract['bytecode'])

assets_dests = [os.path.join(parentdir,'go-synchronizer','contracts','assets'),
                os.path.join(parentdir,'market', 'notary', 'contracts','assets')]

for dest in assets_dests:
    if not os.path.exists(dest):
        os.makedirs(dest)
    cmd = [abigen,
          '--abi',
          os.path.join(mydir,'Assets.abi'),
          '--bin',
          os.path.join(mydir, 'Assets.bin'),
          '--pkg',
          'assets',
          '-out',
          os.path.join(dest, 'assets.go')]
    run_cmd(cmd)

with open(os.path.join(parentdir, 'truffle-core', 'build', 'contracts', 'Market.json'), 'r') as fp:
    contract = json.load(fp)
with open(os.path.join(mydir, 'Market.abi'), 'w') as outfile:
    json.dump(contract['abi'], outfile, ensure_ascii=False, indent=2)
with open(os.path.join(mydir,'Market.bin'), 'w') as outfile:
    outfile.write(contract['bytecode'])

market_dests = [os.path.join(parentdir, 'go-synchronizer', 'contracts', 'market'),
               os.path.join(parentdir, 'market', 'notary', 'contracts', 'market')]

for dest in market_dests:
    if not os.path.exists(dest):
        os.makedirs(dest)
    cmd = [abigen,
          '--abi',
          os.path.join(mydir, 'Market.abi'),
          '--bin',
          os.path.join(mydir, 'Market.bin'),
          '--pkg',
          'market',
          '-out',
          os.path.join(dest, 'market.go')]
    run_cmd(cmd)

# with open('../truffle-core/build/contracts/TeamState.json', 'r') as fp:
#     contract = json.load(fp)
# with open('./States.abi', 'w', encoding='utf-8') as outfile:
#     json.dump(contract['abi'], outfile, ensure_ascii=False, indent=2)
# with open('./States.bin', 'w', encoding='utf-8') as outfile:
#     outfile.write(contract['bytecode'])
# os.system('mkdir -p ../go-synchronizer/contracts/states')
# os.system('abigen --abi ./States.abi --bin ./States.bin --pkg states -out ../go-synchronizer/contracts/states/states.go')
# os.system('mkdir -p ../market/notary/contracts/states')
# os.system('abigen --abi ./States.abi --bin ./States.bin --pkg states -out ../market/notary/contracts/states/states.go')

# with open('../truffle-core/build/contracts/Leagues.json', 'r') as fp:
#     contract = json.load(fp)
# with open('./Leagues.abi', 'w', encoding='utf-8') as outfile:
#     json.dump(contract['abi'], outfile, ensure_ascii=False, indent=2)
# with open('./Leagues.bin', 'w', encoding='utf-8') as outfile:
#     outfile.write(contract['bytecode'])
# os.system('mkdir -p ../go-synchronizer/contracts/leagues')
# os.system('abigen --abi ./Leagues.abi --bin ./Leagues.bin --pkg leagues -out ../go-synchronizer/contracts/leagues/leagues.go')

# with open('../truffle-core/build/contracts/Engine.json', 'r') as fp:
#     contract = json.load(fp)
# with open('./Engine.abi', 'w', encoding='utf-8') as outfile:
#     json.dump(contract['abi'], outfile, ensure_ascii=False, indent=2)
# with open('./Engine.bin', 'w', encoding='utf-8') as outfile:
#     outfile.write(contract['bytecode'])
# os.system('mkdir -p ../go-synchronizer/contracts/engine')
# os.system('abigen --abi ./Engine.abi --bin ./Engine.bin --pkg engine -out ../go-synchronizer/contracts/engine/engine.go')

os.system('rm -rf ./*.abi ./*.bin')
