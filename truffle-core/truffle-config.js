const HDWalletProvider = require("truffle-hdwallet-provider");
 
const mnemonic = "twelve_words";
const infura_apikey = "XXXXXX"; // infura api key

module.exports = {
  networks: {
    ganache: {
      network_id: '*',
      host: '127.0.0.1',
      port: 8545
    },
    dockertest: { // 83a909262608c650bd9b0ae06e29d90d0f67ac5e
      network_id: 63819,
      provider: new HDWalletProvider(
        "FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85",
        "http://localhost:8545"
      ),
      // gasLimit: 2000000
    },
    // infura_ropsten: {
    //   provider: new HDWalletProvider(mnemonic, "https://ropsten.infura.io/"+infura_apikey),
    //   network_id: 3
    // },

    devnet: {
      // public key 0x291081e5a1bF0b9dF6633e4868C88e1FA48900e7
      provider: new HDWalletProvider("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54", "https://devnet.busyverse.com/web3"),
      network_id: 63819
    }
    // infura_ropsten: {
    //   provider: new HDWalletProvider(mnemonic, "https://ropsten.infura.io/"+infura_apikey),
    //   network_id: 3
    // },
  }
}
