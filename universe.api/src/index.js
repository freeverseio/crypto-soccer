const express = require("express");
const { postgraphile } = require("postgraphile");;
const program = require("commander");
const MutationsPlugin = require("./mutations_plugin");
const version = require("../package.json").version;

// Parsing command line arguments
program
  .version(version)
  .option("-d, --databaseUrl <url>", "set the database url", "localhost:5432")
  .parse(process.argv);

const { databaseUrl } = program;

console.log("--------------------------------------------------------");
console.log("databaseUrl       : ", databaseUrl);
console.log("--------------------------------------------------------");

const app = express();
const mutationsPlugin = MutationsPlugin();

app.use(
  postgraphile(
    databaseUrl,
    "public",
    {
      watchPg: true,
      graphiql: true,
      enhanceGraphiql: true,
      retryOnInitFail: true,
      disableDefaultMutations: true,
      appendPlugins: [mutationsPlugin],
    }
  )
);

app.listen(4000);