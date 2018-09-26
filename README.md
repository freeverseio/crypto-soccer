[![CircleCI](https://circleci.com/gh/asiniscalchi/crypto-soccer/tree/master.svg?style=svg)](https://circleci.com/gh/asiniscalchi/crypto-soccer/tree/master)

Install the required dependencies:
cd truffle
npm install

Run tests by starting ganache first:
ganache-cli --deterministic

and then:
cd truffle
truffle test --network ganache
