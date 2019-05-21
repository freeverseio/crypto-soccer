
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

------------

Cost estimates:  1 uint = 20K

Per day: 4 x 24 hashes (all actions per league), one every 15 min.
	- negligible: 100 uints per day, say 200 if you need blocknum too => 4M


Per League (when finished): 
	- UpdateLeague (3 uints: leaguehash, updaterAddr, blockn) => 60K => 80K


Per League (when created): 
	- append to leagues[]  (2 uint)
	- initDataHash, blockSteps, etc (2 uint)
	- sign teams: nTeams x 5K => 100K
	=> 180K => 200K


Exchanging players
	- change state: 2, change teams playerIdxs, 2
	- First time: 40K + 10K = 50K => 70K
	- Afterwards: 10K + 10K = 20K => 40K

Create team:
	- teams[] 2
	- owners[] 2
	- 80K => 100K

Exchange teams:
	- owners[] 2 => 10K => 30K


gorss number= 200K per cada cosa

estressar la BC seria fotre-li un 1% extra de gas


add about 4 slots per league per stakers


--------------------

Total GAS Ethreum consumes in 1 day:
https://etherscan.io/chart/tx

8M gas per block, 1 block every 15s (240 per hour) 
	=> 8M x 240 x 24 = 46e9   GAS per day (suposant els block 100% plens)

Say we want to consume 1% of that GAS => 46e7 GAS per day = 460M GAS per day

Say any of the prev actions is 200K = 0.2M => 2300 accions (comrpa-venda, crear equip, join league, etc) per day

AWESOME

Cost of gas.

1 ETH = 1e9 GWei = 1e18 Wei

1 GAS = 1.5 Gwei = 1.5e-9 ETH = 2e-7 EUR

460M GAS per day ... = 850 Eur per day!!!! = 311K EUR per year!


v0 del cost estimate



-----

From CK: 
The investment also comes though activity on the CryptoKitties platform appears to have waned somewhat since it briefly went viral last year. In December, the platform processed 1.3 million transactions, based on data from Bloxy. This month, roughly 163,340 transactions in October despite a push into Greater China and Singapore earlier this year.

500K kitties sold in 1 year







	
























