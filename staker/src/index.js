const Web3 = require("web3");
const updatesJSON = require("../contracts/Updates.json");
const HDWalletProvider = require("truffle-hdwallet-provider");
const program = require("commander");
const version = require("../package.json").version;

// Parsing command line arguments
program
  .version(version)
  .option("-e, --ethereum <url>", "Ethereum node url", "http://localhost:8545")
  .option("-k, --privateKey <pk>", "private key", "3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
  .option("-u, --updatesContractAddress <address>", "updates contract address", "0xB3F24a605312b744973823194267A99F059e8936")
  .option("-i, --interval <sec>", "interval in sec", "5")
  .parse(process.argv)

const { ethereum, privateKey, updatesContractAddress, interval } = program;

console.log("--------------------------------------------------------");
console.log("ethereum           : ", ethereum);
console.log("🔥  account p.k.    : ", privateKey);
console.log("updates address    : ", updatesContractAddress);
console.log("interval           : ", interval, "sec");
console.log("--------------------------------------------------------");

const provider = new HDWalletProvider(privateKey, ethereum);
const web3 = new Web3(provider, null, {});
const updates = new web3.eth.Contract(updatesJSON.abi, updatesContractAddress);
const from = provider.addresses[0];

const loop = async () => {
  try {
    const currentVerse = await updates.methods.currentVerse().call();
    process.stdout.write("[VERSE: " + currentVerse + "] ");
    process.stdout.write("submitActionsRoot ... ");
    const root = '0x' + (Math.floor(Math.random() * 10000000)).toString(16);
    let gas = await updates.methods.submitActionsRoot(root).estimateGas();
    await updates.methods.submitActionsRoot(root).send({ from, gas });

    // process.stdout.write(", updateTZ ... ")
    // gas = await updates.methods.updateTZ(root).estimateGas();
    // await updates.methods.updateTZ(root).send({ from, gas });

    console.log("done");
  } catch (err) {
    console.error(err);
  }

  setTimeout(() => {
    loop();
  }, interval * 1000);
};

loop();

