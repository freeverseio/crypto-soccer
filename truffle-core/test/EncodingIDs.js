const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();;

const EncodingIDs = artifacts.require('EncodingIDs');

contract('EncodingIDs', (accounts) => {

    beforeEach(async () => {
        encoding = await EncodingIDs.new().should.be.fulfilled;
    });
   
    it('encoding of TZ and country in teamId and playerId', async () =>  {
        encoded = await encoding.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 3, val = 4).should.be.fulfilled;
        decoded = await encoding.decodeTZCountryAndVal(encoded).should.be.fulfilled;
        const {0: timeZone, 1: country, 2: value} = decoded;
        timeZone.toNumber().should.be.equal(tz);
        country.toNumber().should.be.equal(countryIdxInTZ);
        value.toNumber().should.be.equal(val);
    });

    it('get playerID of timezone 1, country 0, index in country 0', async () => {
        encoded = await encoding.encodeTZCountryAndVal(tz = 1, countryIdxInTZ = 0, indexInCountry = 0).should.be.fulfilled;
        encoded.should.be.bignumber.equal('274877906944');
    });
 
});