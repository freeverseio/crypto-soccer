const Stakers = artifacts.require("Stakers")
const expect = require('truffle-assertions');


// TODO: add test where a malicious party is updating and tries to unenroll. Test should fail.

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

  it("Tests stake", async () => {
    assert.equal(0, await web3.eth.getBalance(stakers.address));

    parties = [bob, alice, carol, dave, erin, frank]
    await addTrustedParties(stakers, owner, parties);
    await enroll(stakers, stake, parties);

    assert.equal(parties.length*Number(stake),
                 await web3.eth.getBalance(stakers.address));

    await unenroll(stakers, parties);
    assert.equal(0, await web3.eth.getBalance(stakers.address));
  })

  it("Tests can't unenroll during update", async () => {
    stakers.setGame(game, {from:owner}),
    await stakers.addTrustedParty(bob, {from:owner});
    await stakers.enroll({from:bob, value: stake});
    await stakers.update(0, bob, {from:game}),
    await expect.reverts(
      stakers.unEnroll({from:bob}),
      "failed to unenroll: staker currently updating",
      "bob is currently updating, so it should revert"
    )
    await stakers.start({from:game});
    await expect.passes(
      stakers.unEnroll({from:bob}),
      "failed unenrolling bob"
    )
  })

  it("Tests L0 - > L1 -> start  -> L1, the usual path", async () => {

    stakers.setGame(game, {from:owner}),
    await expect.reverts(
      stakers.update(0, bob, {from:game}),
      "failed to update: staker not registered",
      "bob not yet enrolled, so it should revert"
    )

    parties = [bob, alice, carol, dave, erin, frank]
    await addTrustedParties(stakers, owner, parties);
    await enroll(stakers, stake, parties);

    assert.equal(0, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(0, bob, {from:game}),
      "bob failed to update"
    )
    assert.equal(1, (await stakers.level()).toNumber());
    await expect.reverts(
      stakers.update(0, alice, {from:game}),
      "failed to update: resolving wrong level",
      "level 0 cannot be updated without starting a new verse, it should revert"
    )

    for (i=0; i<10; i++) {
      await expect.passes(
        stakers.start({from:game}),
        "failed starting new verse"
      )

      assert.equal(0, (await stakers.level()).toNumber());
      await expect.passes(
        stakers.update(0, bob, {from:game}),
        "bob failed to update"
      )
      assert.equal(1, (await stakers.level()).toNumber());
    }
  })

  it("Tests L0 -> L1 -> L2 -> L3 -> L4 -> L3 -> L4", async () => {

    // start (L0) ->  bob updates (L1) -> alice updates (L2) -> carol updates (L3) -> dave updates (L4) -> erin challenges dave (L3) -> erin updates (L4)
    stakers.setGame(game, {from:owner}),
    await expect.reverts(
      stakers.update(0, bob, {from:game}),
      "failed to update: staker not registered",
      "bob not yet enrolled, so it should revert"
    )

    parties = [bob, alice, carol, dave, erin, frank]
    await addTrustedParties(stakers, owner, parties);
    await enroll(stakers, stake, parties);

    assert.equal(0, (await stakers.level()).toNumber());

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

    await expect.passes(
      stakers.update(0, bob, {from:game}),
      "bob failed to update"
    )

    assert.equal(1, (await stakers.level()).toNumber());

    await expect.reverts(
      stakers.update(1, bob, {from:game}),
      null,
      "bob is already updating, cannot participate until resolved"
    )

    await expect.passes(
      stakers.update(1, alice, {from:game}),
      "alice failed to update"
    )

    assert.equal(2, (await stakers.level()).toNumber());

    await expect.passes(
      stakers.update(2, carol, {from:game}),
      "carol failed to update"
    )

    assert.equal(3, (await stakers.level()).toNumber());

    await expect.passes(
      stakers.update(3, dave, {from:game}),
      "dave failed to update"
    )

    assert.equal(4, (await stakers.level()).toNumber());

    await expect.passes(
      stakers.update(4, erin, {from:game}),
      "erin failed to update"
    )

    assert.equal(3, (await stakers.level()).toNumber());

    await expect.reverts(
      stakers.update(5, erin, {from:game}),
      "failed to update: wrong level",
      "this update level does not exist, so it should revert"
    )

    assert.equal(3, (await stakers.level()).toNumber());

    await expect.reverts(
      stakers.update(4, erin, {from:game}),
      "failed to update: wrong level",
      "after erin slashed dave, level is 3, so it should revert"
    )

    assert.equal(3, (await stakers.level()).toNumber());

    await expect.reverts(
      stakers.start({from:game}),
      "failed to start: wrong level",
      "starting a new verse is only possible from level 1 or 2, so it should revert"
    )
  })
})

async function asyncForEach(array, callback) {
  for (let index = 0; index < array.length; index++) {
    await callback(array[index], index, array);
  }
}
async function addTrustedParties(contract, owner, addresses) {
  await asyncForEach(addresses, async (address) => {
    contract.addTrustedParty(address, {from:owner})
  });
}
async function enroll(contract, stake, addresses) {
  await asyncForEach(addresses, async (address) => {
    await contract.enroll({from:address, value: stake})
  });
}
async function unenroll(contract, addresses) {
  await asyncForEach(addresses, async (address) => {
    await contract.unEnroll({from:address})
  });
}
