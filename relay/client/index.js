/** Simple example: Create a web page with form */
//import db from './db/db'
var gui = require ( 'easy-web-app' )
var log = require( 'npmlog' )
var fs = require('fs'); // to write to a file
var utils = require('./utils.js')

/** Initialize the framework, the default page and define a title */
var port = 8888
gui.init ( 'Relay client', port )

var formConfig = {
    id   : 'myForm',
    title: 'User Actions',
    type : 'pong-form',
    resourceURL: 'hello',
    height: '750px'
  }

var formPlugInConfig = {
    id: 'myFormDef',
    description: 'shows first form',
    fieldGroups: [
      {
        columns: [
          {
            formFields: [
              { id: 'account', label: 'Account', type: 'text' /*, request: 'header'*/ },
              { id: 'mnemonic', label: 'Mnemonic', type: 'text', rows:'3', defaultVal:'' },
              { id: 'password', label: 'Password', type: 'password' },
              { id: 'actionType', label: 'Action type', type: 'select',
                options:[
                  {option:'tactic',  value:'tactic', selected:true },
                  {option:'sell', value:'sell' },
                  {option:'buy', value:'buy' }
                ]
              },
              { id: 'actionValue', label: 'Action value', type: 'text' },
            ]
          },
          {
            formFields: [
            ]
          }
        ]
      }
    ],
    actions: [
      {  id: 'actionGet',  actionName: 'GET',  method: 'GET',  actionURL: '/relay/v1', setData:[ {resId: 'myForm'} ] },
      {  id: 'actionPost', actionName: 'POST', method: 'POST', actionURL: '/relay/v1' },
    ]
  }

gui.addView ( formConfig, formPlugInConfig )

let app  = gui.getExpress()
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
      res.statusCode = 200    // = OK
      const useraccount  = req.params.useraccount;
      db.map((entry) => {
        if (entry.account === useraccount) {
          return res.status(200).send({
          success: 'true',
          user_account : useraccount,
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
    '/relay/v1',
    jsonParser,
    function( req, res ) {
      res.statusCode = 201    // = created
      console.log( 'POST query: '+JSON.stringify( req.query ))
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
