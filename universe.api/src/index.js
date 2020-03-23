const express = require("express");
const { postgraphile } = require("postgraphile");
const MutationsPlugin = require("./mutations_plugin");
const Web3 = require("web3");
const assetsJSON = require("../contracts/Assets.json");
const program = require("commander");
const version = require("../package.json").version;

// Parsing command line arguments
program
  .version(version)
  .option("-p, --port <port>", "server port", "4000")
  .option("-d, --databaseUrl <url>", "set the database url", "postgres://freeverse:freeverse@localhost:5432/cryptosoccer?sslmode=disable")
  .option("-e, --ethereum <url>", "Ethereum node url", "http://localhost:8545")
  .option("-a, --assetsContractAddress <address>", "assets contract address")
  .option("-s, --sender <address>", "sender address")
  .parse(process.argv)

const { port, databaseUrl, ethereum, assetsContractAddress, sender } = program;


console.log("--------------------------------------------------------");
console.log("port              : ", port);
console.log("databaseUrl       : ", databaseUrl);
console.log("ethereum          : ", ethereum);
console.log("assets address    : ", assetsContractAddress);
console.log("sender            : ", sender);
console.log("--------------------------------------------------------");

const app = express();
const web3 = new Web3(ethereum);
const assets = new web3.eth.Contract(assetsJSON.abi, assetsContractAddress);
const from = sender;
const mutationsPlugin = MutationsPlugin(assets, from);

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
      appendPlugins: [mutationsPlugin]
    }
  )
);

app.listen(port);
