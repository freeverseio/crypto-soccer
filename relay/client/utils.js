var ethUtil = require('ethereumjs-util');
var HDKey = require('hdkey')
var bip39 = require('bip39')

let bytesToHex = function bytesToHex(buff) {
  return `0x${buff.toString('hex')}`;
}

let hexToBytes = function hexToBytes(hex) {
  if (hex.substr(0, 2) === '0x') {
    return Buffer.from(hex.substr(2), 'hex');
  }

  return Buffer.from(hex, 'hex');
}

let generateKeysMnemonic = function generateKeysMnemonic(mnemonic) {
  if (mnemonic == undefined) {
    mnemonic = bip39.generateMnemonic();
  }

  let len = mnemonic.trim().split(/\s+/g).length
  if (len < 12) {
    throw "Mnemonic phrase needs 12 word at least. Only " + len + ' given.'
  }

  const hdkey = HDKey.fromMasterSeed(mnemonic);
  const masterPrivateKey = hdkey.privateKey;
  const masterPubKey = hdkey.publicKey;
  var path = "m/44'/60'/0'/0/0";
  const addrNode = hdkey.derive(path);
  let privK = addrNode._privateKey;
  const pubKey = ethUtil.privateToPublic(addrNode._privateKey);
  let address = ethUtil.privateToAddress(addrNode._privateKey);
  let addressHex = bytesToHex(address);
  let privKHex = bytesToHex(privK);
  //localStorage.setItem(addressHex, privKHex);
  return {address: addressHex, mnemonic: mnemonic};
}

module.exports = {
  hexToBytes, bytesToHex, generateKeysMnemonic
}

/*
function transact() {
	let toAddr = document.getElementById("toAddr").value;
	if(toAddr==undefined) {
		toastr.error("adreça invàlida");
		return;
	}
	if(toAddr=="") { // TODO check also if it's a valid eth address
		toastr.error("adreça invàlida");
		return;
	}
	let amount = Number(100*document.getElementById("amount").value);
	if(amount>myBalance) {
		toastr.error("no hi ha prou saldo");
		return;
	}
	if(amount<=0) {
		toastr.error("la quantitat no es correcte");
		return;
	}
	document.getElementById('spinnerTx').className = 'spinner-border';
	axios.get(RELAYURL + '/tx/nonce/' + myAddr)
	  .then(function (res) {
		    myNonce = res.data.nonce;
		    console.log(res.data);
		    console.log("myNonce " + myNonce);
		  // after getting nonce, generate & sign & send transaction
		let msg = "0x" + buf(uint8(0x19)).toString('hex') + buf(uint8(0)).toString('hex') + buf(TOKENADDR).toString('hex') + buf(uint256(myNonce)).toString('hex') + buf(myAddr).toString('hex') + buf(toAddr).toString('hex') + buf(uint256(amount)).toString('hex')
		let privK = localStorage.getItem(myAddr);
		let sig = ethUtil.ecsign(buf(sha3(msg)),buf(privK));
		let txData = {
			from: myAddr,
			to: toAddr,
			value: Number(amount),
			r: sig.r.toString('hex'),
			s: sig.s.toString('hex'),
			v: sig.v
		};
		console.log(txData);
		axios.post(RELAYURL + '/tx', txData)
		  .then(function (res) {
		    console.log(res.data);
				toastr.success("transferència realitzada");
				$('.nav-tabs a[href="#history"]').tab('show');
			  document.getElementById('spinnerTx').className += 'invisible';
		  })
		  .catch(function (error) {
		    console.log(error);
		    toastr.error(error);
		    document.getElementById('spinnerTx').className += 'invisible';
		  })

	  }) // nonce get error catch
	  .catch(function (error) {
	    console.log(error);
	    toastr.error(error);
	    document.getElementById('spinnerTx').className += 'invisible';
	  })
}
*/
