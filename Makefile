core: 
	cd truffle-core && truffle compire

horizon: core
	mkdir -p ./nodejs-horizon/contracts
	cp ./truffle-core/build/contracts/*.json ./nodejs-horizon/contracts
	cd nodejs-horizon && npm install

test:
	cd truffle-core && truffle test
	cd nodejs-horizon && npm test

clean: 
	rm -rf ./truffle-core/build
	rm -rf ./nodejs-horizon/contracts
	rm -rf ./nodejs-horizon/node_modules


