const Web3 = require("web3");
const updatesJSON = require("../contracts/Updates.json");
const HDWalletProvider = require("truffle-hdwallet-provider");
const program = require("commander");
const version = require("../package.json").version;

// Parsing command line arguments
program
  .version(version)
  .option("-p, --providerUrl <url>", "Ethereum node url", "https://devnet.busyverse.com/web3")
  .option("-k, --privateKey <pk>", "private key", "3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
  .option("-u, --updatesContractAddress <address>", "updates contract address", "0xB3F24a605312b744973823194267A99F059e8936")
  .parse(process.argv)

const { providerUrl, privateKey, updatesContractAddress } = program;

console.log("--------------------------------------------------------");
console.log("providerUrl        : ", providerUrl);
console.log("🔥  account p.k.    : ", privateKey);
console.log("updates address    : ", updatesContractAddress);
console.log("--------------------------------------------------------");

const provider = new HDWalletProvider(privateKey, providerUrl);
const web3 = new Web3(provider, null, {});
const updates = new web3.eth.Contract(updatesJSON.abi, updatesContractAddress);

const loop = () => {
  setTimeout(() => {
    console.log("tick");
    loop();
  }, 3000);
};

loop();

