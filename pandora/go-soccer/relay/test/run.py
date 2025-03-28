#! /usr/bin/env python

import os, sys
import subprocess
import argparse
import time
import requests
import json
from inspect import currentframe

MY_DIR = os.path.dirname(os.path.abspath(__file__))
server_process = None
client_process = None

def get_linenumber():
    cf = currentframe()
    return cf.f_back.f_lineno

def run_async_task(cmd, working_dir = os.getcwd()):
    print 'running', ' '.join(cmd), 'from', working_dir
    return subprocess.Popen(cmd, cwd=working_dir, shell=False, stdout=subprocess.PIPE)

def run_go_server():
    working_dir=os.path.join(MY_DIR,'..','server','test')
    p = run_async_task(['go','run','main.go'], working_dir)
    wait_for_process_to_start(p)
    return p

def run_js_client():
    working_dir=os.path.join(MY_DIR,'..','client')
    p = run_async_task(['node', 'index.js'], working_dir)
    wait_for_process_to_start(p)
    return p

def wait_for_process_to_start(p, timeout=5):
    if not p:
        return
    waiting_for = 0
    while p.poll():
        wait_sec = 5
        waiting_for += wait_sec
        time.sleep(wait_sec)
        if waiting_for > timeout:
            print 'Timed out while waiting for server or client to start. Exiting...'
            return

def stop_process(p):
    if not p:
        return
    try:
        p.terminate()
        p.wait()
        if not p.poll():
            p.kill()
    except OSError as e:
        print 'Error terminating process', e

def stop():
    stop_process(server_process)
    stop_process(client_process)

def FAIL(msg, line_number):
    print 'FAILED at line number', line_number
    print msg
    stop()
    sys.exit(-1)

def ASSERT_EQ(expected,actual, line_number):
    if not expected == actual:
        stop()
        print 'FAILED at line number', line_number
        print "expected:", expected,'\nactual:',actual
        sys.exit(-1)

accounts = [
        {
            "id":1,
            "account":"0x277dac33e16dcbfbd5a2ef9314cc8c232af838fb",
            "privatekey":"0x56060abe29061e8608f4e7f830573b8915328bbfe64d50a34e71d79aa70fa125",
            "mnemonic":"ocean merge switch power planet proud woman cargo vendor brass small lens"
        },
        {
            "id":2,
            "account":"0x8cd6e7c05d127dd63d433a4ccfb1ab97cb889b70",
            "privatekey":"0x5a26f4631043a8a8ffee55a19a885ae37e58ebac3c7d56e015119fbf156b2079",
            "mnemonic":"cancel liar pull crater trial across polar logic cool other force squirrel"
        },
        {
            "id":3,
            "account":"0x7781d1843be87a6d7927b5aa3d4603d3fdc113bc",
            "privatekey":"0x656af771e6633101072b4d8c93639cef0b519f98010b0c175829ccd06506a3c3",
            "mnemonic":"surge arctic network orphan script fortune gown scheme rebuild congress emerge mountain"
        },
]

client_db = "/tmp/useraccounts.txt"

if __name__ == "__main__":

    if os.path.exists(client_db):
        os.remove(client_db)

    server_process = run_go_server()
    if not server_process:
        print 'Error server process not running'
        sys.exit(-1)

    client_process = run_js_client()
    if not client_process:
        print 'Error client not running'
        stop_process(server_process)
        sys.exit(-1)

    time.sleep(5) # TODO: wait until server and client are fully available

    print 'TEST: creating wallets in client app'

    for account in accounts:
        endpoint='http://localhost:8888/createwallet'
        data = json.dumps({'mnemonic' : account['mnemonic']})
        headers = {'Content-Type': 'application/json'}
        r = requests.post(url = endpoint, data = data, headers=headers)
        if not r:
            FAIL('POST failed for account ' + str(account['id']), get_linenumber())
        d = r.json()
        if not d['success'] == 'true':
            FAIL('response not successful for account ' + str(account['id']), get_linenumber())

        entry = d['entry']
        ASSERT_EQ(account['account'], entry['account'], get_linenumber())
        ASSERT_EQ(account['privatekey'], entry['privatekey'], get_linenumber())
        ASSERT_EQ(account['mnemonic'], entry['mnemonic'], get_linenumber())

    print 'TEST: create the same wallet'
    ASSERT_EQ("false", requests.post(
            url = 'http://localhost:8888/createwallet',
            data = json.dumps({'mnemonic' : accounts[0]['mnemonic']}),
            headers={'Content-Type': 'application/json'}
            ).json()['success'], get_linenumber())


    print 'TEST: register client accounts to the server'
    for account in accounts:
        endpoint='http://localhost:8888/relay/v1/' + account['account']
        r = requests.get(url = endpoint)
        if not r:
            FAIL('GET failed for account ' + str(account['id']), get_linenumber())
        d = r.json()
        ASSERT_EQ('User created', d['message'], get_linenumber())
        ASSERT_EQ(account['account'], d['user'], get_linenumber())

    print 'TEST: submit user actions'
    #http://localhost:8888/relay/v1/0x82973f0ceed111576c508bcd999c92c9e83e49f0/action?type=sell&value=player
    for account in accounts:
        endpoint='http://localhost:8888/relay/v1/%(account)s/action'%{'account' : account['account']}
        payload = {
                'type': 'type_' + str(account['id']),
                'value': 'value_' + str(account['id']),
                }
        r = requests.get(url = endpoint, params = payload)
        if not r:
            FAIL('GET failed for account ' + str(account['id']), get_linenumber())
        d = r.json()
        ASSERT_EQ(True, d['success'], get_linenumber())
        ASSERT_EQ(True, d['verified'], get_linenumber())
        ASSERT_EQ(account['account'], d['useraddr'], get_linenumber())
        action = d['action']
        ASSERT_EQ(payload['type'], action['Type'], get_linenumber())
        ASSERT_EQ(payload['value'], action['Value'], get_linenumber())

    print "SUCCESS"

    stop_process(server_process)
    stop_process(client_process)
