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
  .option("-e, --ethereum <url>", "Ethereum node url", "http://localhost:8545")
  .option("-k, --privateKey <pk>", "private key", "3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
  .option("-a, --assetsContractAddress <address>", "assets contract address")
  .option("-i, --interval <sec>", "interval in sec", "5")
  .parse(process.argv)

const { ethereum, privateKey, assetsContractAddress } = program;


console.log("--------------------------------------------------------");
console.log("ethereum       : ", ethereum);
console.log("🔥  account p.k.   : ", privateKey);
console.log("assets address    : ", assetsContractAddress);
console.log("--------------------------------------------------------");

const provider = new HDWalletProvider(privateKey, ethereum);
const web3 = new Web3(provider, null, {});
const assets = new web3.eth.Contract(assetsJSON.abi, assetsContractAddress);
const from = provider.addresses[0];

const resolvers = new Resolvers({
  assets,
  from
});

const server = new ApolloServer({
  typeDefs,
  resolvers
});

server.listen().then(({ url }) => {
  console.log(`🚀  Server ready at ${url}`);
});
