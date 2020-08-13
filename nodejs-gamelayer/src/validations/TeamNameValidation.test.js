const TeamNameValidation = require("./TeamNameValidation.js")
const Web3 = require('web3')

test('Keccak256 madafakka', () => {
  const web3 = new Web3(Web3.givenProvider || "ws://localhost:8545")
  teamNameValidation = new TeamNameValidation({ web3, name: '', teamId: '3', signature: ''});
  const hash = teamNameValidation.hash()
  console.log("hash", hash)
  web3.currentProvider.disconnect();
})

test('verify signature', () => {
  const web3 = new Web3(Web3.givenProvider || "ws://localhost:8545")
  teamNameValidation = new TeamNameValidation({ web3, name: 'ciao', teamId: '4', signature: '3feac668bb718f492638b9b58d1f294379cdc8bde40074f5e49c3f80f28190e121f0fd08227c64a643dd032748ef772b0d1cf1500f649345521c133290c941a91b'});

  const res = teamNameValidation.verifySignature()
  console.log("res", res)
  web3.currentProvider.disconnect();
})