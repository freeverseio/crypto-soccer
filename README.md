Install the required dependencies:
cd truffle
npm install

Run tests by starting ganache first:
ganache-cli --deterministic

and then:
cd truffle
truffle test --network ganache