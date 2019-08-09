const Stakers = artifacts.require("Stakers")
const expect = require('truffle-assertions');

contract('Stakers', (accounts) => {
  const [owner, game, bob, alice, carol, dave, erin, frank] = accounts

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

  it("Tests enrolling", async () => {
    await expect.reverts(
      stakers.enroll({from:bob, value: stake}),
      null,
      "bob is not yet a trusted party, so it should revert"
    )
    await expect.reverts(
      stakers.addTrustedParty(bob, {from:game}),
      null,
      "only owner can add trusted parties, so it should revert"
    )
    await expect.passes(
      stakers.addTrustedParty(bob, {from:owner}),
      "failed to add bob as trusted party"
    )
    await expect.reverts(
      stakers.addTrustedParty(bob, {from:owner}),
      null,
      "bob is already a trusted party, so it should revert"
    )
    await expect.passes(
      stakers.enroll({from:bob, value: stake}),
      "failed to enroll bob"
    )
  })

  it("Tests level to update", async () => {
    stakers.setGame(game, {from:owner}),
    await expect.reverts(
      stakers.update(1, bob, {from:game}),
      "failed to update: wrong level",
      "level to update is 1, so it should revert"
    )
    await expect.reverts(
      stakers.update(2, bob, {from:game}),
      "failed to update: wrong level",
      "level to update is 1, so it should revert"
    )
    await expect.reverts(
      stakers.update(0, bob, {from:game}),
      "failed to update: staker not registered",
      "bob not yet enrolled, so it should revert"
    )

    parties = [bob, alice, carol, dave, erin, frank]
    await parties.forEach(function(address) {
      stakers.addTrustedParty(address, {from:owner})
    })
    await parties.forEach(function(address) {
      stakers.enroll({from:address, value: stake})
    })

    await expect.passes(
      stakers.update(0, bob, {from:game}),
      "bob failed to update"
    )
    await expect.reverts(
      stakers.update(1, bob, {from:game}),
      null,
      "bob is already updating, cannot participate until next verse"
    )
    await expect.passes(
      stakers.update(1, alice, {from:game}),
      "alice failed to update"
    )
    await expect.passes(
      stakers.update(2, carol, {from:game}),
      "carol failed to update"
    )
    await expect.reverts(
      stakers.update(2, dave, {from:game}),
      "failed to update: wrong level",
      "level to update is 4, so it should revert"
    )
    await expect.passes(
      stakers.update(3, dave, {from:game}),
      "dave failed to update"
    )
    await expect.passes(
      stakers.update(4, erin, {from:game}),
      "erin failed to update"
    )

    // TODO: this fails
    await expect.reverts(
      stakers.update(5, erin, {from:game}),
      "failed to update: level too large",
      "this update level does not exist, so it should revert"
    )

    await expect.reverts(
      stakers.start({from:game}),
      "failed starting new verse from current level",
      "starting a new verse is only possible from level 1 or 2, so it should revert"
    )
  })
})
