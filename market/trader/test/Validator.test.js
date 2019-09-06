const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();

describe('Validator', () => {
    it('create signature of "ciao"', async () => {
        
    });
}) 