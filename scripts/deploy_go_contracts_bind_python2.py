#!/usr/bin/env python

import json
import os,sys
import subprocess
from distutils.spawn import find_executable

mydir = os.path.dirname(os.path.realpath(__file__))
parentdir = os.path.dirname(mydir)
abigen=os.path.join(os.getenv("HOME"), 'go', 'bin', 'abigen')
truffle_contracts_dir=os.path.join(parentdir, 'truffle-core', 'build', 'contracts')

def openfile(filename, flags):
    if sys.version_info[0] < 3:
        return  open(filename, flags)
    return  open(filename, flags, encoding='utf-8')

def deploy_go_contract(jsonfile, pkgname, destinations):
    print 'deployning ' + pkgname
    base, ext = os.path.splitext(os.path.basename(jsonfile))
    abifile = os.path.join(mydir, base + '.abi')
    binfile = os.path.join(mydir, base + '.bin')
    with openfile(os.path.join(truffle_contracts_dir, jsonfile), 'r') as fp:
        contract = json.load(fp)
    with openfile(abifile, 'w') as outfile:
        json.dump(contract['abi'], outfile, ensure_ascii=False, indent=2)
    with openfile(binfile, 'w') as outfile:
        outfile.write(contract['bytecode'])

    for dest in destinations:
        outputdir = os.path.join(dest, pkgname)
        if not os.path.exists(outputdir):
            os.makedirs(outputdir)
        cmd = [abigen,
              '--abi',
              abifile,
              '--bin',
              binfile,
              '--pkg',
              pkgname,
              '-out',
              os.path.join(outputdir, pkgname + '.go')]
        run_cmd(cmd)

def run_cmd(cmd, working_dir = os.getcwd()):
    try:
        print ' '.join(cmd)
        output = subprocess.Popen(cmd, cwd=working_dir, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    except subprocess.CalledProcessError as e:
        print '"' + ' '.join(e.cmd) + '" failed with return code ' + str(e.returncode)
        sys.exit(e.returncode)
    return output.communicate()

if __name__ == "__main__":

    if len(sys.argv) == 2:
        abigen = sys.argv[1]
    if not os.path.exists(abigen):
        if find_executable('abigen'):
            abigen = 'abigen'
        else:
            print 'Could not find abigen'
            print 'usage: python deploy_go_contracts_bind_python2.py [abigen_path]'
            sys.exit(-1)

    dests = [os.path.join(parentdir,'go','contracts'),
             ]

    deploy_go_contract(os.path.join(truffle_contracts_dir, 'Assets.json'), 'assets', dests)
    deploy_go_contract(os.path.join(truffle_contracts_dir, 'Market.json'), 'market', dests)
    deploy_go_contract(os.path.join(truffle_contracts_dir, 'Updates.json'), 'updates', dests)
    deploy_go_contract(os.path.join(truffle_contracts_dir, 'Leagues.json'), 'leagues', dests)
    deploy_go_contract(os.path.join(truffle_contracts_dir, 'Engine.json'), 'engine', dests)
    deploy_go_contract(os.path.join(truffle_contracts_dir, 'EnginePreComp.json'), 'engineprecomp', dests)
    deploy_go_contract(os.path.join(truffle_contracts_dir, 'Evolution.json'), 'evolution', dests)


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
