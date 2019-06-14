
https://etherscan.io/chart/tx

8M gas per block
1 block every 15s
1 Tx without anything else = 21K gas

380 tx per block => 25 TXps

since they do stuff => 10 TXps

100K per day
31M Tx per year

--------

Cost estimates:

Per day: 4 x 24 hashes (all actions per league), one every 15 min.
	- negligible

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











	



