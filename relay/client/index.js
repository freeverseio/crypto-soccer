/** Simple example: Create a web page with form */
var gui = require ( 'easy-web-app' )
var log = require( 'npmlog' )

/** Initialize the framework, the default page and define a title */
var port = 8888
gui.init ( 'Relay client', port )

/**
 * Add a view of type "pong-easy-form" (= plug-in) to the default page the first
 * parameter of addView is the view configuration, a second parameter can define
 * the plug-in configuration, a third parameter can specify the page.
 */
var view = gui.addView (
  {
    'id'   : 'UserActions',
    'type' : 'pong-easyform'
  },
  {
    "id" : "tstFormId",
    "easyFormFields" : [
        "id"
      , "c1|User~address"
      , "c1|Pass~Word"
      , "c1|Mnemonic|10rows"
      //, "c1|Date|date"
      , "c1|Action"
      //, "c1|separator"
      //, "c1|Remark|3rows"
      //, "c2|Mailings|label"   // second starts here column
      //, "c2|Send~Ads~~|checkbox_infomails_ads"
      //, "c2|Newsletter|checkbox_infomails_newsletter"
    ],
    "actions" : [
      {
        id : "userActionPost",
        actionName : "send",
        actionURL  : "/relay/v1",
        method : 'POST',
        setData:[ {resId: 'myForm'} ]
      },
      {
        id : "userActionGet",
        actionName : "get",
        actionURL  : "/relay/v1",
        method : 'GET',
        setData:[ {resId: 'tstFormId'} ]
      }
    ]
  }
)

var app  = gui.getExpress()
var bodyParser  = require( 'body-parser' )
var jsonParser  = bodyParser.json()

app.get(
    '/relay/v1',
    function( req, res ) {
      res.statusCode = 200    // = OK
      console.log( 'GET Params: '+JSON.stringify( req.query ) )
      return  res.json( {} )
    }
  )
app.post(
    '/relay/v1',
    jsonParser,
    function( req, res ) {
      log.info( 'app.post:' + view)
      console.log( 'POST Body: '+JSON.stringify( req.body ))
      res.statusCode = 200    // = OK
      return res.json( {
        useraddr: 'foobar', action:'Hello world!'
      } )
    }
  )

/*
  gui.getExpress().post( '/test', ( req, res ) => {
  res.status( 200 ).send(
    'Hello to you!!\n'+
    'and to everyone'
  )
})
*/
