let express = require('express')
const app = express();
let utils = require('./utils.js')
let db = require('./db.js')
let bodyParser  = require( 'body-parser' )

const PORT = 8888
let jsonParser  = bodyParser.json()
let formParser  = bodyParser.urlencoded( { extended: true } );

app.listen(PORT, () => {
  console.log(`server running on port ${PORT}`)
});

// ----------------------------
// GET
// ----------------------------
//

//http://localhost:8888/relay/db
app.get( // TODO: just for debugging. remove
    '/relay/db',
    function( req, res ) {
      res.statusCode = 200    // = OK
      return  res.json( db.getData() )
    }
  )


//http://localhost:8888/relay/v1?account=1234
app.get(
    '/relay/v1',
    function( req, res ) {
      console.log( 'GET Params: '+JSON.stringify( req.query ) )
      const useraccount = req.query.account;
      return res.redirect('/relay/v1/'+useraccount)
    }
  )

//http://localhost:8888/relay/v1/1234
app.get(
    '/relay/v1/:useraccount',
    function( req, res ) {
      console.log( 'GET Params: '+JSON.stringify( req.params ) )
      const useraccount = req.params.useraccount;
      const entry = db.getUserEntry(useraccount)

      if (entry) {
        return res.status(200).send({
          success: 'true',
          message: 'account retrieved successfully',
          entry,
        });
      }
      return res.status(404).send({
        success: 'false',
        message: 'account ' + useraccount + ' does not exist',
      });
    }
  )

//http://localhost:8888/relay/v1/1234/nonce
app.get(
    '/relay/v1/:useraccount/nonce',
    function( req, res ) {
      console.log( 'GET Params: '+JSON.stringify( req.params ) )
      const useraccount = req.params.useraccount;
      const entry = db.getUserEntry(useraccount)

      var nonce = -1
      if (entry) {
        utils.getAccountNonce(useraccount).then(
        function(nonce) {
          return res.status(200).send({
            success: 'true',
            account : useraccount,
            nonce: nonce
          });
        }).catch(function(err) {
          console.error(err)
        });
      }
      else
      {
        return res.status(404).send({
          success: 'false',
          message: 'account ' + useraccount + ' does not exist',
        });
      }
    }
  )

// ----------------------------
// POST
// ----------------------------

// curl -v -H "Content-Type: application/json" -X POST -d '{"mnemonic":"a b c d e"}' http://localhost:8888/relay/v1
app.post(
    '/relay/v1/createuser',
    jsonParser,
    function( req, res ) {
      console.log( 'POST Body: '+JSON.stringify( req.body ))

      var usermnemonic;
      if(req.body.mnemonic)
      {
        usermnemonic = req.body.mnemonic
      }

      let wallet = utils.generateKeysMnemonic(usermnemonic);

      const account = wallet.address
      const mnemonic = wallet.mnemonic

      console.log('account: ', account, "mnemonic: ", mnemonic)
      const entry = db.getUserEntry(account)

      if (entry) {
        return res.status(400).send({
          success: 'false',
          message: 'user already exists',
          entry,
       });
      }

      const new_entry = db.addUserEntry(account, mnemonic);

      return res.status(201).send({
        success: 'true',
        entry: new_entry
      });
    }
  )
