contracts: 
	cd truffle-core && truffle compile
	mkdir -p ./nodejs-horizon/contracts
	cp ./truffle-core/build/contracts/*.json ./nodejs-horizon/contracts
	cd ./scripts && ./deploy_go_contracts_bind.py

test:
	cd truffle-core && truffle test
	cd nodejs-horizon && npm test
	cd go-synchronizer && go test ./...

clean: 
	rm -rf ./truffle-core/build
	rm -rf ./nodejs-horizon/contracts
	rm -rf ./go-synchronizer/contracts


