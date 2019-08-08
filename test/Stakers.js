const Stakers = artifacts.require("Stakers")
const expect = require('truffle-assertions');

contract('Stakers', (accounts) => {
  const [owner, game, bob, alice, carol] = accounts

  let stakers
  let stake

  beforeEach(async () => {
      stakers  = await Stakers.new({from:owner})
      stake = await stakers.kRequiredStake()
  });

  it("Tests game address", async () => {
    await expect.reverts(
      stakers.update(1, bob),
      "Only game can call this function",
      "game not set yet, so it should revert"
    )
    await expect.reverts(
      stakers.setGame(game, {from:bob}),
      "Only owner can call this function",
      "wrong sender, so it should revert"
    )
    await expect.passes(
      stakers.setGame(game, {from:owner}),
      "failed to set game address"
    )
  })

  it("Tests level to update", async () => {
    stakers.setGame(game, {from:owner}),
    await expect.reverts(
      stakers.update(0, bob, {from:game}),
      "cannot update: wrong level",
      "level to update is 1, so it should revert"
    )
    await expect.reverts(
      stakers.update(2, bob, {from:game}),
      "cannot update: wrong level",
      "level to update is 1, so it should revert"
    )
    await expect.reverts(
      stakers.update(1, bob, {from:game}),
      "cannot update: staker not registered",
      "bob not yet enrolled, so it should revert"
    )
    await expect.passes(
      stakers.enroll(bob, {from:game, value: stake}),
      "failed enrolling bob"
    )
    // TODO: check that the money is taken from game, but it should be taken from bob
    await expect.passes(
      stakers.update(1, bob, {from:game}),
      "failed to update"
    )
  })
})
