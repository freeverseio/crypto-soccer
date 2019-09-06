const Accounts = require('web3-eth-accounts');

class Validator {
    constructor() {
        this.accounts = new Accounts();
    }

    recoverSignerAddress(msg, signature) {
        return this.accounts.recover(msg, signature);
    }
}

module.exports = Validator;