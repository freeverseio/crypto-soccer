all: 
	cd truffle-core && truffle compile
	mkdir -p ./nodejs-horizon/contracts
	cp ./truffle-core/build/contracts/*.json ./nodejs-horizon/contracts
	cd ./scripts && ./deploy_go_contracts_bind.py

clean: 
	rm -rf ./truffle-core/build
	rm -rf ./nodejs-horizon/contracts
	rm -rf ./go-synchronizer/contracts


