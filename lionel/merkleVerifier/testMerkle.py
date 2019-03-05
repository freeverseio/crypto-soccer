import hashlib
import math
from sha3 import keccak_256
from merkle_tree import *


def prepareProofForIdxs(idxsToProve, tree):
    # neededHashes
    neededHashes = proof(tree, idxsToProve)
    values = {}
    for idx in idxsToProve:
        values[idx] = leafs[idx]
    return neededHashes, values

leafs = range(16)

tree, depth = make_tree(leafs)

print("Root: " + str(root(tree)))

neededHashes0 = proof(tree, [0])
neededHashes1 = proof(tree, [1])

print(neededHashes0[1]==neededHashes1[1])

neededHashes2 = proof(tree, [2])
neededHashes3 = proof(tree, [3])

print(neededHashes2[1]==neededHashes3[1])

neededHashes01 = proof(tree, [0,1])

print(neededHashes01[0]==neededHashes0[1])


# Actual test. Choose some Idxs. Get the needed hashes. Pass the values of the leafs at those idxs,
# prepared as a 'dictionary' of the form { idx: value }

# Test 1
idxsToProve = [0,1,2,3]
neededHashes, values = prepareProofForIdxs(idxsToProve, tree)
print("To prove these %i leafs you need %i hashes, in a tree with %i leafs, and depth %i" \
        % (len(idxsToProve), len(neededHashes),len(leafs), depth)
      )
print(verify(root(tree), depth, values, neededHashes,  debug_print=False))

# Test 2
idxsToProve = [0,1,2,4]
neededHashes, values = prepareProofForIdxs(idxsToProve, tree)
print("To prove these %i leafs you need %i hashes, in a tree with %i leafs, and depth %i" \
        % (len(idxsToProve), len(neededHashes),len(leafs), depth)
      )
print(verify(root(tree), depth, values, neededHashes,  debug_print=False))

# Test 3
idxsToProve = [0]
neededHashes, values = prepareProofForIdxs(idxsToProve, tree)
print("To prove these %i leafs you need %i hashes, in a tree with %i leafs, and depth %i" \
        % (len(idxsToProve), len(neededHashes),len(leafs), depth)
      )
print(verify(root(tree), depth, values, neededHashes,  debug_print=False))
