import hashlib
import math
from sha3 import keccak_256
# from pylio import serialHash
# import pylio

def hash_leaf(leaf_value):
    # Toni: I do not need to convert this to anything
    # return pylio.serialHash(leaf_value)
    '''Convert a leaf value to a digest'''
    assert(leaf_value < 2**256)
    return leaf_value.to_bytes(32, 'big')

def hash_node(left_hash, right_hash):
    # I need to make this more general
    return
    '''Convert two digests to their Merkle node's digest'''
    return keccak_256(left_hash + right_hash).digest()

def paddToPowerTwo(leafs):
    num_leafs = len(leafs)
    assert num_leafs > 0, "Merkle Tree undefined for zero leafs"
    depth = int(math.log2(num_leafs))
    if num_leafs > 2**depth:
        for new_leaf in range(num_leafs, 2**(depth+1)):
            leafs.append(0)



def make_tree(leafs, hashFunction):
    '''Compute the Merkle tree of a list of values.
    The result is returned as a list where each value represents one hash in the
    tree. The indices in the array are as in a bbinary heap array.
    '''
    # INPORTANT: you need to pass here a deepcopy of leafs, otherwise, it'll modify them!!!
    # it first doubles the tree so that in the 2nd half it places the hashes of the pre-hashed data
    # then it stores the parent hashes. For example, in a tree with 4 leaves, it stores:
    # step 1:   undef0,      undef1, undef2, undef3, data0Hash, data1Hash, data2Hash, data3Hash
    # step 2:   undef0, hash(01,23), hash01, hash23, data0Hash, data1Hash, data2Hash, data3Hash
    # note that hash(01,23) = MerkleRoot
    #
    # the hash function "hash_leaf" is here just dummy (converts to bytes)
    paddToPowerTwo(leafs)
    num_leafs = len(leafs)
    depth = int(math.log2(num_leafs))
    assert(num_leafs == 2**depth)
    num_nodes = 2 * num_leafs
    tree = [None] * num_nodes
    for i in range(num_leafs):
        tree[2**depth + i] = hashFunction(leafs[i]).to_bytes(32, 'big')
    for i in range(2**depth - 1, 0, -1):
        tree[i] = hashFunction(tree[2*i] + tree[2*i + 1])
    return tree, depth

def root(tree):
    return tree[1]

def proof(tree, indices):
    '''Given a Merkle tree and a set of indices, provide a list of decommitments
    required to reconstruct the merkle root.'''
    depth = int(math.log2(len(tree))) - 1
    num_leafs = 2**depth
    num_nodes = 2*num_leafs
    known = [False] * num_nodes
    decommitment = []
    for i in indices:
        known[2**depth + i] = True
    for i in range(2**depth - 1, 0, -1):
        left = known[2*i]
        right = known[2*i + 1]
        if left and not right:
            decommitment += [tree[2*i + 1]]
        if not left and right:
            decommitment += [tree[2*i]]
        known[i] = left or right
    return decommitment

def verify(root, depth, values, decommitment, hashFunction, debug_print=False):
    '''Verify a set of leafs in the Merkle tree.
    
    Parameters
    ------------------------
    root
        Merkle root that is commited to.
    depth
        Depth of the Merkle tree. Equal to log2(number of leafs)
    values
        Mapping leaf index => value of the values we want to decommit.
    decommitments
        List of intermediate values required for deconstruction.
    '''
    
    # Create a list of pairs [(tree_index, leaf_hash)] with tree_index decreasing
    queue = []
    for index in sorted(values.keys(), reverse=True):
        tree_index = 2**depth + index
        hash = hashFunction(values[index]).to_bytes(32, 'big')
        queue += [(tree_index, hash)]

    while True:
        assert(len(queue) >= 1)

        # Take the top from the queue
        (index, hash) = queue[0]
        queue = queue[1:]
        if debug_print:
            print(index, hash.hex())

        # The merkle root has tree index 1
        if index == 1:
            return hash == root
        
        # Even nodes get merged with a decommitment hash on the right
        elif index % 2 == 0:
            queue += [(index // 2, hashFunction(hash + decommitment[0]))]
            decommitment = decommitment[1:]
        
        # Odd nodes can get merged with their neighbour
        elif len(queue) > 0 and queue[0][0] == index - 1:
                # Take the sibbling node from the stack
                (_, sibbling_hash) = queue[0]
                queue = queue[1:]

                # Merge the two nodes
                queue += [(index // 2, hashFunction(sibbling_hash + hash))]
        
        # Remaining odd nodes are merged with a decommitment on the left
        else:
            # Merge with a decommitment hash on the left
            queue += [(index // 2, hashFunction(decommitment[0] + hash))]
            decommitment = decommitment[1:]


def get_depth(tree):
    # the tree is twice larger than the number of leafs, hence our /2 below
    return int(math.log2(len(tree)/2))
