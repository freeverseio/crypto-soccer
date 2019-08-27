const express = require("express");
const { postgraphile } = require("postgraphile");

const app = express();

app.use(
  postgraphile(
    process.env.DATABASE_URL || "postgres://freeverse:freeverse@localhost:5432/market",
    "public",
    {
      watchPg: true,
      graphiql: true,
      enhanceGraphiql: true,
      retryOnInitFail: true
    }
  )
);

app.listen(process.env.PORT || 3000);