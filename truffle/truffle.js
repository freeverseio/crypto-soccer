module.exports = { 
  networks: {
    ganache: {
    host: "127.0.0.1", port: 8545, network_id: "*"
    } 
  },
  mocha: {
    reporter: 'eth-gas-reporter',
    reporterOptions : {
      currency: 'EUR',
      gasPrice: 21
    }
  }
};
