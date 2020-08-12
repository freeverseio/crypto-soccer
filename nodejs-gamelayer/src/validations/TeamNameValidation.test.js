const TeamNameValidation = require("./TeamNameValidation.js")

test('Keccak256 madafakka', () => {
  const hash = TeamNameValidation.hash()
  console.log("hash", hash)
})
