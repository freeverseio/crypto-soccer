function compareArraysInternal(result, expected, toNum, verbose, isBigNumber) {
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

function compareArrays(result, expected, toNum = true, isBigNumber = false) {
  try {
    compareArraysInternal(result, expected, toNum, isBigNumber)
  } 
  catch(e) {
    compareArraysInternal(result, expected, toNum, verbose = true, isBigNumber)
    throw e
  }  
}
  
  module.exports = {
    compareArrays
  }