# Install the required dependencies:
cd truffle-core
npm install

# Run tests by starting ganache first:
ganache-cli --deterministic

# and then:
cd truffle
truffle test --network ganache

# if you prefer using the console
truffle develop
test

# To use the UI: (make sure you have compiled first)
cd truffle
truffle compile
python -m SimpleHTTPServer
go to a browser and connect to: http://0.0.0.0:8000/UI/index.html

# To use a solidity IDE for VS, use I recommend Juan Blanco's plugin.
