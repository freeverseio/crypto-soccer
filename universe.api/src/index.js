const express = require("express");
const { postgraphile } = require("postgraphile");
const program = require("commander");
const version = require("../package.json").version;
const mutationsPlugin = require("./mutations_plugin");
const mutationsWrapperPlugin =  require("./mutation_wrapper_plugin");

// Parsing command line arguments
program
  .version(version)
  .option("-p, --port <port>", "server port", "4000")
  .option("-d, --databaseUrl <url>", "set the database url", "postgres://freeverse:freeverse@localhost:5432/cryptosoccer?sslmode=disable")
  .parse(process.argv)

const { port, databaseUrl } = program;

console.log("--------------------------------------------------------");
console.log("port              : ", port);
console.log("databaseUrl       : ", databaseUrl);
console.log("--------------------------------------------------------");

const app = express();

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
      appendPlugins: [mutationsPlugin, mutationsWrapperPlugin],
    }
  )
);

app.listen(port);
