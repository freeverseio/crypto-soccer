setup:
	cd truffle-core && npm install

contracts:
	cd truffle-core && ./node_modules/.bin/truffle compile
	mkdir -p nodejs-horizon
	cp -r truffle-core/build/contracts ./relay/nodejs-api
	cp -r truffle-core/build/contracts ./staker
	cd scripts && ./deploy_go_contracts_bind_python2.py

clean:
	rm -rf ./truffle-core/build
	rm -rf ./relay/nodejs-api/contracts
	rm -rf ./go-synchronizer/contracts
	rm -rf ./market/notary/contracts

