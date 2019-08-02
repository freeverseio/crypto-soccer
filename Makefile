contracts:
	cd truffle-core && npm install
	cd truffle-core && ./node_modules/.bin/truffle compile
	mkdir -p nodejs-horizon
	cp -r truffle-core/build/contracts ./nodejs-horizon
	cd scripts && ./deploy_go_contracts_bind.py

