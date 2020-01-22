function compareArrays(result, expected, toNum = true, verbose = false, isBigNumber = false) {
    verb = [];
    for (i = 0; i < expected.length; i++) {
        if (toNum) verb.push(result[i].toNumber());
        else verb.push(result[i]);
        if (!verbose) {
            if (toNum) result[i].toNumber().should.be.equal(expected[i]);
            else {
              if (isBigNumber) result[i].should.be.bignumber.equal(expected[i]);
              else result[i].should.be.equal(expected[i]);
            }
        }            
    }
    if (verbose) console.log(verb);
}
  
  module.exports = {
    compareArrays
  }