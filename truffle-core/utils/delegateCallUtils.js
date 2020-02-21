const UNDEF = undefined;
    
function extractSelectorsFromAbi(abi) {
    functions = [];
    for (i = 0; i < abi.length; i++) { 
        if (abi[i].type == "function") {
            functions.push(web3.eth.abi.encodeFunctionSignature(abi[i]));
        }
    }    
    return functions;
}

function removeDuplicatesFromFirstContract(selectors1, selectors2) {
    selectors1_filtered = [];
    for (s = 0; s < selectors1.length; s++) {
        if (!selectors2.includes(selectors1[s])) {
            selectors1_filtered.push(selectors1[s]);
        }
    }
    return selectors1_filtered;
}
  
module.exports = {
    extractSelectorsFromAbi,
    removeDuplicatesFromFirstContract
}