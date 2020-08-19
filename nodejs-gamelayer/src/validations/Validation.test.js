const TeamValidation = require("./TeamValidation.js")
const PlayerValidation = require("./PlayerValidation.js")
const Web3 = require('web3')

test('Team Hash', async () => {
  const web3 = new Web3(Web3.givenProvider || "ws://localhost:8545")
  teamValidation = new TeamValidation({ web3, name: '', teamId: '3', signature: ''});
  const hash = teamValidation.hash()

  expect(hash).toBe("0x074b4277787bca36334cf57f0507141ef743a08d7690dba02af123626e6955d0")
})

test('Team PrefixedHash', async () => {
  const web3 = new Web3(Web3.givenProvider || "ws://localhost:8545")
  const signature = '3feac668bb718f492638b9b58d1f294379cdc8bde40074f5e49c3f80f28190e121f0fd08227c64a643dd032748ef772b0d1cf1500f649345521c133290c941a91b'
  teamValidation = new TeamValidation({ web3, name: 'ciao', teamId: '4', signature });
  const hash = teamValidation.prefixedHash()

  expect(hash).toBe("0x34e71acf56fc4ad1a6e219dcc96bf5111d8092a6ab64308281a7c77525d2a404")
})

test('Team signerAddress', async () => {
  const web3 = new Web3(Web3.givenProvider || "ws://localhost:8545")
  const signature = '3feac668bb718f492638b9b58d1f294379cdc8bde40074f5e49c3f80f28190e121f0fd08227c64a643dd032748ef772b0d1cf1500f649345521c133290c941a91b'
  teamValidation = new TeamValidation({ web3, name: 'ciao', teamId: '4', signature });
  const signerAddress = await teamValidation.signerAddress()

  expect(signerAddress).toBe('0x291081e5a1bF0b9dF6633e4868C88e1FA48900e7')
})

test('Player Hash', async () => {
  const web3 = new Web3(Web3.givenProvider || "ws://localhost:8545")
  playerValidation = new PlayerValidation({ web3, name: '', playerId: '3', signature: ''});
  const hash = playerValidation.hash()

  expect(hash).toBe("0x074b4277787bca36334cf57f0507141ef743a08d7690dba02af123626e6955d0")
})

test('Player PrefixedHash', async () => {
  const web3 = new Web3(Web3.givenProvider || "ws://localhost:8545")
  const signature = '3feac668bb718f492638b9b58d1f294379cdc8bde40074f5e49c3f80f28190e121f0fd08227c64a643dd032748ef772b0d1cf1500f649345521c133290c941a91b'
  playerValidation = new PlayerValidation({ web3, name: 'ciao', playerId: '4', signature });
  const hash = playerValidation.prefixedHash()

  expect(hash).toBe("0x34e71acf56fc4ad1a6e219dcc96bf5111d8092a6ab64308281a7c77525d2a404")
})

test('Player signerAddress', async () => {
  const web3 = new Web3(Web3.givenProvider || "ws://localhost:8545")
  const signature = '3feac668bb718f492638b9b58d1f294379cdc8bde40074f5e49c3f80f28190e121f0fd08227c64a643dd032748ef772b0d1cf1500f649345521c133290c941a91b'
  playerValidation = new PlayerValidation({ web3, name: 'ciao', playerId: '4', signature });
  const signerAddress = await playerValidation.signerAddress()

  expect(signerAddress).toBe('0x291081e5a1bF0b9dF6633e4868C88e1FA48900e7')
})

