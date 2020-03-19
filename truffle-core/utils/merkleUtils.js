function hash_node(x, y) {
  return web3.utils.keccak256(web3.eth.abi.encodeParameters(['bytes32', 'bytes32'], [x,y]));
}

function merkleRoot(leafs, nLevels) {
  _leafs = [...leafs];
  nLeafs = 2**nLevels;
  assert.equal(_leafs.length, nLeafs, "number of leafs is not = pow(2,nLevels)");
  for (level = 0; level < nLevels; level++) {
      nLeafs = Math.floor(nLeafs/2);
      for (pos = 0; pos < nLeafs; pos++) {
        _leafs[pos] = hash_node(_leafs[2 * pos], _leafs[2 * pos + 1]);      
      }
  }
  return _leafs[0];
}

function verify(root, proof, leafHash, leafPos) {
  for (pos = 0; pos < proof.length; pos++) {
      if ((leafPos % 2) == 0) {
          leafHash = hash_node(leafHash, proof[pos]);
      } else {
          leafHash = hash_node(proof[pos], leafHash);
      }
      leafPos = Math.floor(leafPos/2);
  }     
  return root == leafHash;   
}

  module.exports = {
    hash_node,
    merkleRoot
  }