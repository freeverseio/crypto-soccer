// local storage for storing user wallets

let fs = require('fs'); // to write to a file
var data = [] // TODO use a proper database
const path = '/tmp/useraccounts.txt'
load(path)

function getData() {
  return data;
}

function load() {
  if (fs.existsSync(path)) {
    try {
      data = JSON.parse(fs.readFileSync(path))
    }
    catch(error)
    {
      console.error(error);
    }
  }
}

function save() {
  fs.writeFileSync(path, JSON.stringify(data), function(err){
      if(err) {
          console.log(err)
      } else {
          console.log('Database updated');
      }
  });
}

function getAccount(useraccount) {
  for (var entry of data) {
    if (entry.account === useraccount) {
      return entry
    }
  }
  return null
}

function addAccount(account, privatekey, mnemonic) {
      const new_entry = {
        id : data.length + 1,
        account : account,
        privatekey: privatekey,
        mnemonic : mnemonic
      }
      data.push(new_entry);
      save();
      return new_entry
}

module.exports = {
  getData, getAccount, addAccount
}
