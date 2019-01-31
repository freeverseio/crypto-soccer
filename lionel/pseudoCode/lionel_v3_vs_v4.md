Notation:

D = 'Data' = {initHash, finalStatesHash[10], scores[90] }
M = MerkleRoot:   D ---> M(D) 


Rule 1: provided D anyone can prove that one of its elems is wrong, by playing 1 team across league as a TX

Rule 2: only M is not enough to be proven wrong


## Lionel V3

- updater:  aporta D
- challenger: demostra que algun dels elem de D is wrong => if wrong, reset



## Lionel V4

- Updater escriutx $M(D)$
- Challanger1 escriutx $M_{c1}(D')$ on $D'$= {initHash FS1...FS9  R1..R90} pero $FS1_{ch1}$ es FALS
- Challanger2 simulategametx per Team=1 i demostra que $FS1_{ch1}$ es FALS, fa un slash de Challanger1

## Aventatges de LionelV4 sobre LionelV3

- El updater no cal que escrigui D gas=14*20000 = 280k gas
- 'Certificar una lliga' en LionelV3 = 14 cops mes car que LionelV4

## TODOS
    - (toni) canvis tactics
    - incentius
    - (toni) random num
















