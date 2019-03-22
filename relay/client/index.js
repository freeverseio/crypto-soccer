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

// http://localhost:8888/relay/db
// shows local user wallets (just for debugging)
app.get(
    '/relay/debug',
    function( req, res ) {
      res.statusCode = 200    // = OK
      return  res.json( db.getData() )
    }
  )

// http://localhost:8888/relay/v1/1234
// will post a Relay server to create a user with given address
app.get(
    '/relay/v1/:useraddr',
    function( req, res ) {
      console.log( 'GET Params: '+JSON.stringify( req.params ) )
      const useraddr = req.params.useraddr;
      const entry = db.getAccount(useraddr)

      if (!entry) {
        return res.status(404).send({
          success: 'false',
          message: 'account ' + useraddr + ' does not exist',
        });
      }

      utils.createUser(useraddr)
      .then(function(r) {
        console.log(r.data);
        return res.status(200).send({
          user: r.data.user,
          message: r.data.message
        });
      })
      .catch(function (error) {
        console.log(error);
      });
    }
)

// http://localhost:8888/account/:useraddr/action?type=xyz&value=123
app.get(
    '/relay/v1/:useraddr/action',
    function( req, res ) {
      //console.log( 'GET Params: '+JSON.stringify( req.params ) )
      const useraddr = req.params.useraddr;
      const entry = db.getAccount(useraddr)

      if (!entry) {
        return res.status(404).send({
          success: 'false',
          message: 'account ' + useraddr + ' does not exist',
        });
      }

      const actionType = req.query.type
      const actionValue = req.query.value
      utils.submitAction(
        useraddr,
        entry.privatekey,
        actionType,
        actionValue
      )
      .then(function(r) {
        console.log(r.data)
        res.status(200).send({
          success: true,
          user: r.data.user,
          action: r.data.action,
          verified: r.data.verified,
          message: r.data.message
        })
      })
      .catch(function(error) {
        console.log(error);
      });
    }
  )

// ----------------------------
// POST
// ----------------------------

// curl -v -H "Content-Type: application/json" -X POST -d '{"mnemonic":"a b c d e"}' http://localhost:8888/createwallet
// creates a wallet with the given mnemonic phrase and adds it to local db. If mnemonic is empty, mnemonic will be generated
app.post(
    '/createwallet',
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
      const privatekey = wallet.privatekey

      console.log('account: ', account, "mnemonic: ", mnemonic)
      const entry = db.getAccount(account)

      if (entry) {
        return res.status(400).send({
          success: 'false',
          message: 'account already exists',
          entry,
       });
      }

      const new_entry = db.addAccount(account, privatekey, mnemonic);

      return res.status(201).send({
        success: 'true',
        entry: new_entry
      });
    }
  )
