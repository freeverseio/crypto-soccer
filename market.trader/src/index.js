const express = require("express");
const { postgraphile } = require("postgraphile");
const program = require("commander");
const version = require("../package.json").version;
const ConnectionFilterPlugin = require("postgraphile-plugin-connection-filter");

// Parsing command line arguments
program
  .version(version)
  .option("-d, --databaseUrl <url>", "set the database url")
  .option("-p, --port <port>", "server port", "4000")
  .option("-o, --enableCors <bool>", "enables some generous CORS settings for the GraphQL endpoint. There are some costs associated when enabling this, if at all possible try to put your API behind a reverse proxy", "false")
  .parse(process.argv);

const { databaseUrl, port , enableCors} = program;

console.log("--------------------------------------------------------");
console.log("databaseUrl       : ", databaseUrl);
console.log("server port       : ", port);
console.log("enable CORS       : ", enableCors);
console.log("--------------------------------------------------------");

const app = express();

app.use(
  postgraphile(
    databaseUrl,
    "public",
    {
      enableCors: enableCors,
      watchPg: true,
      graphiql: true,
      enhanceGraphiql: true,
      retryOnInitFail: true,
      disableDefaultMutations: true,
      appendPlugins: [ConnectionFilterPlugin],
    }
  )
);

app.listen(port);