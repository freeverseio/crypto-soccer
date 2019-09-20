const { ApolloServer } = require("apollo-server");
const typeDefs = require("./schema");
const Resolvers = require("./resolvers");
const Web3 = require("web3");
const marketJSON = require("../contracts/Market.json");
const HDWalletProvider = require("truffle-hdwallet-provider");
const program = require("commander");
const version = require("../package.json").version;

// Parsing command line arguments
program
  .version(version)
  .option("-c, --config <path>", "set config path. defaults to config.json")
  .parse(process.argv);

let configFile = "../";
if (typeof program.config !== "undefined") configFile += program.config;
else configFile += "config.json";

console.log("Configuration file: " + configFile);
const config = require(configFile);
const {
  providerUrl,
  address,
  privateKey,
  marketContractAddress,
} = config;

console.log("--------------------------------------------------------");
console.log("providerUrl       : ", providerUrl);
console.log("account           : ", address);
console.log("🔥  account p.k.   : ", privateKey);
console.log("market address    : ", marketContractAddress);
console.log("--------------------------------------------------------");

const provider = new HDWalletProvider(privateKey, providerUrl);
const web3 = new Web3(provider, null, {});
const market = new web3.eth.Contract(marketJSON.abi, marketContractAddress);

const resolvers = new Resolvers({
  market,
  from: address
});

const server = new ApolloServer({
  typeDefs,
  resolvers
});

server.listen().then(({ url }) => {
  console.log(`🚀  Server ready at ${url}`);
});
