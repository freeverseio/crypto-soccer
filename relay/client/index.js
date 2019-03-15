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

app.get(
    '/relay/v1',
    function( req, res ) {
      res.statusCode = 200    // = OK
      console.log('account: ' + req.query.account)
      console.log('mnemonic: ' + req.query.mnemonic)
      console.log('password: ' + req.query.password)
      console.log('actionType: ' + req.query.actionType)
      console.log('actionValue: ' + req.query.actionValue)
      console.log( 'GET Params: '+JSON.stringify( req.query ) )
      return  res.json( {account : "ok!" } )
    }
  )

//http://localhost:8888/relay/db
app.get( // TODO: just for debugging. remove
    //u
    '/relay/db',
    function( req, res ) {
      res.statusCode = 200    // = OK
      return  res.json( db.getData() )
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

      if (entry) {
        return res.status(200).send({
          success: 'true',
          account : useraccount,
          nonce: 0
        });
      }
      return res.status(404).send({
        success: 'false',
        message: 'account ' + useraccount + ' does not exist',
      });
    }
  )

// ----------------------------
// POST
// ----------------------------

// curl -v -H "Content-Type: application/json" -X POST -d '{"name":"your name","phonenumber":"111-111"}' http://localhost:8888/relay/v1
app.post(
    '/relay/v1/createuser',
    jsonParser,
    function( req, res ) {
      console.log( 'POST Body: '+JSON.stringify( req.body ))

      //if(!req.body.account)
      //{
      //  return res.status(400).send({
      //    success: 'false',
      //    message: 'account is required'
      //  });
      //}
      //else if(!req.body.mnemonic)
      //{
      //  return res.status(400).send({
      //    success: 'false',
      //    message: 'description is required'
      //  });
      //}

      let wallet = utils.generateKeysMnemonic(null);

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
