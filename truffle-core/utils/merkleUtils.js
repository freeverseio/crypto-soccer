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

function buildProof(leafPos, leafs, nLevels) {
  _leafs = [...leafs];
  nLeafs = 2**nLevels;
  assert.equal(_leafs.length, nLeafs, "number of leafs is not = pow(2,nLevels)");
  proof = [];
  // The 1st element is just its pair
  proof.push( 
    ((leafPos % 2) == 0) ? _leafs[leafPos+1] : _leafs[leafPos-1]
  );
  // The rest requires computing all hashes
  for (level = 0; level < nLevels-1; level++) {
      nLeafs = Math.floor(nLeafs/2);
      leafPos = Math.floor(leafPos/2);
      for (pos = 0; pos < nLeafs; pos++) {
          _leafs[pos] = hash_node(_leafs[2 * pos], _leafs[2 * pos + 1]);      
      }
      proof.push(
        ((leafPos % 2) == 0) ? _leafs[leafPos+1] : _leafs[leafPos-1]
      );
  }
  return proof;
}

function getBaseLog(base, x) {
  return Math.log(x) / Math.log(base);
}

function buildMerkleStruct(leafs, nLeafsPerRoot) {
  levelsPerRoot = Math.floor(Math.log2(nLeafsPerRoot));
  assert.equal(nLeafsPerRoot, 2**levelsPerRoot, "nLeafsPerRoot must be a power of 2");
  
  nTotalLeafs = leafs.length;
  nChallenges = getBaseLog(nLeafsPerRoot, nTotalLeafs);
  assert.equal(nTotalLeafs, nLeafsPerRoot**nChallenges, "nTotalLeafs should be a power of nLeafsPerRoot");
  

  rootsPerLevel = [];
  leafsAtThisLevel = [...leafs];

  for (ch = 0; ch < nChallenges - 1; ch++) {
      rootsAtThisLevel = [];
      assert.equal(leafsAtThisLevel.length % nLeafsPerRoot, 0, "wrong number of leafs");
      nRootsToCompute = leafsAtThisLevel.length/nLeafsPerRoot;
      console.log("new level: ", ch, nLeafsPerRoot)
      for (n = 0; n < nRootsToCompute; n++) {
          left = n * nLeafsPerRoot;
          right = (n+1)*nLeafsPerRoot
          thisRoot = merkleRoot(leafsAtThisLevel.slice(left, right), levelsPerRoot);
          rootsAtThisLevel.push(thisRoot)
      }
      leafsAtThisLevel = [...rootsAtThisLevel];
      rootsPerLevel.push([...rootsAtThisLevel]);
  }
  assert.equal(
      merkleRoot(leafsAtThisLevel, levelsPerRoot),
      merkleRoot(leafs, nLev = Math.log2(nTotalLeafs)),
      "the merkle struct built does not have a correct merkle root"
  );
}
  
  module.exports = {
    hash_node,
    merkleRoot,
    verify,
    buildProof,
    buildMerkleStruct
  }