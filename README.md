When deploying new contracts:

* add them to scripts/deploy_go_contracts_bind_python2.py
* make contracts
* in go/testutils/blockchain_node.go:
  * add to ContractAddresses all contracts that you want to call from Go code
  * 
To update golden tests in go: go test matches_test.go setup_test.go -test.update-golden
