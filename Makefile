setup:
	cd truffle-core && npm install
	cd nodejs-horizon && npm install

core:
	cd truffle-core && ./node_modules/.bin/truffle compile

contracts_deploy:
	mkdir -p nodejs-horizon
	cp -r truffle-core/build/contracts ./nodejs-horizon
	cd scripts && ./deploy_go_contracts_bind.py

core_test:
	cd truffle-core && ./node_modules/.bin/truffle test

horizon_test:
	cd nodejs-horizon && npm test

synchronizer:
	cd go-synchronizer && go build

synchronizer_test:
	cd go-synchronizer && go test ./...

clean:
	rm -rf ./truffle-core/build
	rm -rf ./nodejs-horizon/contracts
	rm -rf ./go-synchronizer/contracts

deepclean: clean
	rm -rf ./truffle-core/node_modules
	rm -rf ./nodejs-horizon/node_modules

build: core contracts_deploy synchronizer
test: core_test horizon_test synchronizer_test
