setup:
	cd truffle-core && npm install

contracts:
	cd truffle-core && ./node_modules/.bin/truffle compile
	cp -r truffle-core/build/contracts ./relay.api
	cd scripts && ./deploy_go_contracts_bind_python2.py

clean:
	rm -rf ./truffle-core/build
	rm -rf ./relay.api/contracts
	rm -rf ./go-synchronizer/contracts
	rm -rf ./market.notary/contracts

