/*
 * NB: since truffle-hdwallet-provider 0.0.5 you must wrap HDWallet providers in a 
 * function when declaring them. Failure to do so will cause commands to hang. ex:
 * ```
 * mainnet: {
 *     provider: function() { 
 *       return new HDWalletProvider(mnemonic, 'https://mainnet.infura.io/<infura-key>') 
 *     },
 *     network_id: '1',
 *     gas: 4500000,
 *     gasPrice: 10000000000,
 *   },
 */
const HDWalletProvider = require("truffle-hdwallet-provider");
 
const mnemonic = "twelve_words";
const infura_apikey = "XXXXXX"; // infura api key

module.exports = {
    // See <http://truffleframework.com/docs/advanced/configuration>
  // to customize your Truffle configuration!
  networks: {
    ganache: {
      network_id: '*',
      host: '127.0.0.1',
      port: 8545
    },
    infura_ropsten: {
      provider: new HDWalletProvider(mnemonic, "https://ropsten.infura.io/"+infura_apikey),
      network_id: 3
    },
    devnet: {
      provider: new HDWalletProvider("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54", "https://devnet.busyverse.com/web3"),
      network_id: 63819
    }
  },
  // mocha: {
  //   reporter: 'eth-gas-reporter',
  //   reporterOptions: {
  //     currency: 'EUR',
  //     gasPrice: 21
  //   }
  // },
  compilers: {
    solc: {
      version: "0.5.8"
    }
  }
}
