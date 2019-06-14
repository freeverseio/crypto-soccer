
team actions: A1 = (u1, a1, sig1)

where sig1 proves that (u1, a1) was submitted by someone who knew sk1, in other words:

sig1 = E( Hash(u1,a1), sk1 )

One can prove that D(sig1, pk1) = Hash(u1, a1), and hence ask:

isActionCorrect(A1, pk1):
	return D(sig1, pk1) = Hash(u1, a1)


------------

Let's hash every N minutes.
Let's limit to 1 action per user per hash.
Let's group all actions... per league.

We only commit the leagues that will be played in the next few mins.

Size of Merkle proof: log2 N (where one of these is the Merkle Root)

L1 ={A1, A2, A3, ..., A_nTeams}  

ALL_ACTIONS = {L1, L2, L3, ...}

Challenge init state

-------------

Worst case: 
	- write the hash of each league before each round (twice per game):


1 League = 2 nMatchdays = 4 x (nTeams-1)

Better case:
	- write the hash of all leagues being played that block => basically scales indep of leagues

Still, one can prove that an action was there by providing the Merkle Proof

verifyAllActionsOfMatchdayOne(L1, MerkleProof(L1)):
	verify Hash(L1), together with MerkleProof(L1) is MerkleRoot.

-




















