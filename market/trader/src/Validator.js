const Accounts = require('web3-eth-accounts');

class Validator {
    constructor() {
        this.accounts = new Accounts();
    }

    recoverSignerAddress(msgHash, signature) {
        return this.accounts.recover(msgHash, signature, true);
    }

    recoverRSV(signature) {
        signature = signature.substr(2); //remove 0x
        const r = '0x' + signature.slice(0, 64)
        const s = '0x' + signature.slice(64, 128)
        const v = '0x' + signature.slice(128, 130)
        return { r, s, v };
    }
}

module.exports = Validator;