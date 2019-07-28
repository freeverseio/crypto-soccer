core: 
	cd truffle-core && truffle compire

horizon: core
	mkdir -p ./nodejs-horizon/contracts
	cp ./truffle-core/build/contracts/*.json ./nodejs-horizon/contracts
	cd nodejs-horizon && npm install

synchronizer: core
	cd ./scripts && ./deploy_go_contracts_bind.py

test:
	cd truffle-core && truffle test
	cd nodejs-horizon && npm test
	cd go-synchronizer && go test ./...

clean: 
	rm -rf ./truffle-core/build
	rm -rf ./nodejs-horizon/contracts
	rm -rf ./nodejs-horizon/node_modules


