const HDWalletProvider = require("truffle-hdwallet-provider");

module.exports = {
  compilers: {
    solc: {
      version: "0.6.3", // A version or constraint - Ex. "^0.5.0"
                         // Can also be set to "native" to use a native solc
      parser: "solcjs",  // Leverages solc-js purely for speedy parsing
      settings: {
        optimizer: {
          enabled: true,
        }
      }
    }
  },
  networks: {
    ganache: {
      network_id: '*',
      host: '127.0.0.1',
      port: 8545
    },
    xdai: { // 0xA9c0F76cA045163E28afDdFe035ec76a44f5C1F3
      provider: new HDWalletProvider(
        "a98c8730d71a46bcc40fb06fc68142edbc2fdf17b89197db0fbe41d35718d5fc",
        "https://dai.poa.network"
      ),
      network_id: 100,
      gasPrice: 1000000000
    },
    xdaidev: { // 0xA9c0F76cA045163E28afDdFe035ec76a44f5C1F3
      provider: new HDWalletProvider(
        "a98c8730d71a46bcc40fb06fc68142edbc2fdf17b89197db0fbe41d35718d5fc",
        "https://dai.poa.network"
      ),
      network_id: 100,
      gasPrice: 1000000000
    },
    local: { // 0x291081e5a1bF0b9dF6633e4868C88e1FA48900e7
      provider: new HDWalletProvider(
        "FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85",
        "http://localhost:8545"
      ),
      network_id: 63819
    },
    dev: { // 0xA9c0F76cA045163E28afDdFe035ec76a44f5C1F3
      provider: new HDWalletProvider(
        "FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85",
        "https://k8s.gorengine.com/xdai"
      ),
      network_id: 63819
    },
  },

  // Set default mocha options here, use special reporters etc.
  // mocha: {
  //   reporter: 'eth-gas-reporter',
  //   timeout: 100000
  // }
}
