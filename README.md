When deploying new contracts:

* add them to scripts/deploy_go_contracts_bind_python2.py
* make contracts
* in go/testutils/blockchain_node.go:
  * add to ContractAddresses all contracts that you want to call from Go code
  * 
