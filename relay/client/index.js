//import db from './db/db'
let express = require('express')
const app = express();
let fs = require('fs'); // to write to a file
let utils = require('./utils.js')

const PORT = 8888

app.listen(PORT, () => {
  console.log(`server running on port ${PORT}`)
});
let bodyParser  = require( 'body-parser' )
let jsonParser  = bodyParser.json()
let formParser  = bodyParser.urlencoded( { extended: true } );
let db = [] // TODO use a proper database
let mnemonic = 'public lion man jelly mom fitness awkward muscle target cactus coast depth'

let wallet = utils.generateKeysMnemonic(mnemonic);
console.log('wallet: ', wallet.address, "mnemonic: " , wallet.mnemonic)


function writeDatabase() {
  const path = '/tmp/relaydb.txt'
  fs.writeFileSync(path, JSON.stringify(db), function(err){
      if(err) {
          console.log(err)
      } else {
          console.log('Database updated');
      }
  });
}

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
app.get( // just for debugging
    //u
    '/relay/db',
    function( req, res ) {
      res.statusCode = 200    // = OK
      return  res.json( {db} )
    }
  )

//http://localhost:8888/relay/v1/1234
app.get(
    '/relay/v1/:useraccount',
    function( req, res ) {
      console.log( 'GET Params: '+JSON.stringify( req.params ) )
      const useraccount  = req.params.useraccount;

      db.map((entry) => {
        if (entry.account === useraccount) {
          return res.status(200).send({
          success: 'true',
          message: 'account retrieved successfully',
          entry,
          });
        }
      });

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
      //res.statusCode = 200    // = OK
      const useraccount  = req.params.useraccount;
      db.map((entry) => {
        if (entry.account === useraccount) {
          return res.status(200).send({
          success: 'true',
          account : useraccount,
          nonce: 0,
          });
        }
      });

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
      res.statusCode = 201    // = created
      console.log( 'POST Body: '+JSON.stringify( req.body ))

      if(!req.body.account)
      {
        return res.status(400).send({
          success: 'false',
          message: 'account is required'
        });
      }
      else if(!req.body.mnemonic)
      {
        return res.status(400).send({
          success: 'false',
          message: 'description is required'
        });
      }

      const useraccount = req.body.account
      db.map((entry) => {
        if (entry.account === useraccount) {
          return res.status(400).send({
          success: 'false',
          message: 'user already exists',
          entry,
          });
        }
      });

      const entry = {
        id : db.length + 1,
        account : req.body.account,
        mnemonic : req.body.mnemonic
      }

      db.push(entry);
      writeDatabase();

      return res.json({success: 'true', entry: entry})
    }
  )
