const express = require("express");
const { postgraphile } = require("postgraphile");
const mutations = require("./mutations");

const { ApolloServer } = require("apollo-server");
const typeDefs = require("./schema");
const Resolvers = require("./resolvers");
const Web3 = require("web3");
const assetsJSON = require("../contracts/Assets.json");
const HDWalletProvider = require("truffle-hdwallet-provider");
const program = require("commander");
const version = require("../package.json").version;

// Parsing command line arguments
program
  .version(version)
  .option("-p, --port <port>", "server port", "4000")
  .option("-d, --databaseUrl <url>", "set the database url")
  .option("-e, --ethereum <url>", "Ethereum node url", "http://localhost:8545")
  .option("-k, --privateKey <pk>", "private key", "3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
  .option("-a, --assetsContractAddress <address>", "assets contract address")
  .option("-i, --interval <sec>", "interval in sec", "5")
  .parse(process.argv)

const { port, databaseUrl, ethereum, privateKey, assetsContractAddress } = program;


console.log("--------------------------------------------------------");
console.log("databaseUrl       : ", databaseUrl);
console.log("ethereum       : ", ethereum);
console.log("🔥  account p.k.   : ", privateKey);
console.log("assets address    : ", assetsContractAddress);
console.log("--------------------------------------------------------");

const app = express();
const provider = new HDWalletProvider(privateKey, ethereum);
const web3 = new Web3(provider, null, {});
const assets = new web3.eth.Contract(assetsJSON.abi, assetsContractAddress);
const from = provider.addresses[0];

app.use(
  postgraphile(
    databaseUrl,
    "public",
    {
      watchPg: true,
      graphiql: true,
      enhanceGraphiql: true,
      retryOnInitFail: true,
      // disableDefaultMutations: true,
      appendPlugins: [mutations]
    }
  )
);

app.listen(port);
