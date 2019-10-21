const HDWalletProvider = require("truffle-hdwallet-provider");

module.exports = {
  networks: {
    ganache: {
      network_id: '*',
      host: '127.0.0.1',
      port: 8545
    },
    devnet: { // 0x291081e5a1bF0b9dF6633e4868C88e1FA48900e7
      provider: new HDWalletProvider(
        "FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85",
        "http://localhost:8545"
      ),
      network_id: 63819
    }
  },

  // Set default mocha options here, use special reporters etc.
  mocha: {
    // reporter: 'eth-gas-reporter',
    // timeout: 100000
  },
}
