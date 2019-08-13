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

////////////////////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////////////////////

  it("Tests L0 - > L1 true -> start -> L1 true, the usual path", async () => {

    stakers.setGame(game, {from:owner}),
    parties = [bob, alice, carol, dave, erin, frank]
    await addTrustedParties(stakers, owner, parties);
    await enroll(stakers, stake, parties);

    // L0
    assert.equal(0, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(0, bob, {from:game}),
      "bob failed to update"
    )

    // L1
    assert.equal(1, (await stakers.level()).toNumber());
    await expect.reverts(
      stakers.update(0, alice, {from:game}),
      "failed to update: resolving wrong level",
      "level 0 cannot be updated without starting a new verse, it should revert"
    )

    for (i=0; i<10; i++) {
      // start new verse
      await expect.passes(
        stakers.start({from:game}),
        "failed starting new verse"
      )

      // L0
      assert.equal(0, (await stakers.level()).toNumber());
      await expect.passes(
        stakers.update(0, bob, {from:game}),
        "bob failed to update"
      )

      // L1
      assert.equal(1, (await stakers.level()).toNumber());
    }
  })

////////////////////////////////////////////////////////////////////////////////////////////

  it("Tests L0 - > L1 lie -> L2 true -> start -> L1 lie -> L2 true", async () => {

    stakers.setGame(game, {from:owner}),
    parties = [bob, alice, carol, dave, erin, frank]
    await addTrustedParties(stakers, owner, parties);
    await enroll(stakers, stake, parties);

    // L0
    assert.equal(0, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(0, bob, {from:game}),
      "bob failed to update"
    )

    // L1
    assert.equal(1, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(1, alice, {from:game}),
      "alice failed to update"
    )

    // L2
    assert.equal(2, (await stakers.level()).toNumber());
    await expect.reverts(
      stakers.update(0, carol, {from:game}),
      "failed to update: resolving wrong level",
      "level 0 cannot be updated without starting a new verse, it should revert"
    )
    await expect.reverts(
      stakers.update(1, carol, {from:game}),
      "failed to update: resolving wrong level",
      "level 1 is already updated, it should revert"
    )

    // ------------- start new verse ----------------
    await expect.passes(
      stakers.start({from:game}),
      "failed starting new verse"
    )

    // L0
    assert.equal(0, (await stakers.level()).toNumber());
    await expect.reverts(
      stakers.update(0, bob, {from:game}),
      "failed to update: staker not registered",
      "bob was slashed by alice and therefore it is removed from registered stakers, so it should revert"
    )
    await expect.reverts(
      stakers.enroll({from:bob, value: stake}),
      "failed to enroll: staker was slashed",
      "bob was slashed by alice it can no longer enroll, so it should revert"
    )
    await expect.passes(
      stakers.update(0, alice, {from:game}),
      "alice failed to update after new verse"
    )

    // L1
    assert.equal(1, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(1, dave, {from:game}),
      "dave failed to update level 1"
    )

    // L2
    assert.equal(2, (await stakers.level()).toNumber());
  })

////////////////////////////////////////////////////////////////////////////////////////////

  it("Tests L0 -> L1 lie -> L2 lie -> L3 true -> L1 -> L2 true", async () => {

    stakers.setGame(game, {from:owner}),
    parties = [bob, alice, carol, dave, erin, frank]
    await addTrustedParties(stakers, owner, parties);
    await enroll(stakers, stake, parties);

    // L0
    assert.equal(0, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(0, bob, {from:game}),
      "bob failed to update"
    )

    // L1
    assert.equal(1, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(1, alice, {from:game}),
      "alice failed to update"
    )

    // L2
    assert.equal(2, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(2, dave, {from:game}),
      "dave failed to update"
    )
    daveBalance = Number(await web3.eth.getBalance(dave));

    // L3
    assert.equal(3, (await stakers.level()).toNumber());

    //await expect.reverts(
    //  stakers.start({from:game}),
    //  "failed to start: wrong level",
    //  "can't start a new verse from L3, so it should revert"
    //)

    // challenge time has passed, resolve from L3: alice will be slashed and dave earns alice's stake
    await expect.passes(
      stakers.update(1, erin, {from:game}),
      "erin failed to update"
    )

    // L2: because L3 was resolved with an update to L1, we are now at L2
    assert.equal(2, (await stakers.level()).toNumber());

    await expect.reverts(
      stakers.enroll({from:alice, value: stake}),
      "failed to enroll: staker was slashed",
      "alice was slashed, so it should revert"
    )

    assert.isBelow(daveBalance, Number(await web3.eth.getBalance(dave)),
                 "Dave current balance should be higher now, since he earned Alice's stake");

    // start new verse
    await expect.passes(
      stakers.start({from:game}),
      "failed starting new verse"
    )

    await expect.reverts(
      stakers.enroll({from:alice, value: stake}),
      "failed to enroll: staker was slashed",
      "alice was slashed and will never be able to enroll again, so it should revert"
    )
  })

////////////////////////////////////////////////////////////////////////////////////////////

  it("Tests L0 -> L1 true -> L2 lie -> L3 lie -> L4 true -> L2 -> L3 true -> start", async () => {
    stakers.setGame(game, {from:owner}),
    parties = [bob, alice, carol, dave, erin, frank]
    await addTrustedParties(stakers, owner, parties);
    await enroll(stakers, stake, parties);

    // L0
    assert.equal(0, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(0, bob, {from:game}),
      "bob failed to update"
    )

    // L1 - true
    assert.equal(1, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(1, alice, {from:game}),
      "alice failed to update"
    )

    // L2 - lie
    assert.equal(2, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(2, dave, {from:game}),
      "dave failed to update"
    )

    // L3 - lie
    assert.equal(3, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(3, erin, {from:game}),
      "erin failed to update to L4"
    )
    erinBalance = Number(await web3.eth.getBalance(erin));

    // L4 - true
    assert.equal(4, (await stakers.level()).toNumber());

    // challenge time passed, resolve from L4: erin told the truth,  dave will be slashed and erin earns dave's stake
    await expect.passes(
      stakers.update(2, frank, {from:game}),
      "frank failed to update to L4"
    )
    // L3
    assert.equal(3, (await stakers.level()).toNumber());
    assert.isBelow(erinBalance, Number(await web3.eth.getBalance(erin)),
                 "Erin current balance should be higher now, since she earned Dave's stake");
    await expect.reverts(
      stakers.enroll({from:dave, value: stake}),
      "failed to enroll: staker was slashed",
      "dave was slashed and will never be able to enroll again, so it should revert"
    )

    // challenge time for L3 has passed, and also challenge time for L1 has passed.
    // In other words frank  and bob told the truth, and the game can now call start
    // resolving that frank earns alice's stake, and bob earns reward. Also alice will
    // be slashed

    frankBalance = Number(await web3.eth.getBalance(frank));
    await expect.passes(
      stakers.start({from:game}),
      "failed calling start from L3"
    )

    // TODO: this fails?
    assert.isBelow(frankBalance, Number(await web3.eth.getBalance(frank)),
                 "Frank's current balance should be higher now, since he earned Alice's stake");

    await expect.reverts(
      stakers.enroll({from:alice, value: stake}),
      "failed to enroll: staker was slashed",
      "alice was slashed and will never be able to enroll again, so it should revert"
    )
  })

////////////////////////////////////////////////////////////////////////////////////////////

  it("Tests L0 -> L1 lie -> L2 true -> L3 lie -> L4 true -> start", async () => {
    stakers.setGame(game, {from:owner})
    parties = [bob, alice, carol, dave, erin, frank]
    await addTrustedParties(stakers, owner, parties);
    await enroll(stakers, stake, parties);
    // L0
    assert.equal(0, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(0, bob, {from:game}),
      "bob failed to update"
    )

    // L1 - lie
    assert.equal(1, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(1, alice, {from:game}),
      "alice failed to update"
    )

    // L2 - true
    assert.equal(2, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(2, dave, {from:game}),
      "dave failed to update"
    )

    // L3 - lie
    assert.equal(3, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(3, erin, {from:game}),
      "erin failed to update to L4"
    )

    // L4 - true
    assert.equal(4, (await stakers.level()).toNumber());

    // challenge time for L4 has passed, and also challenge time for L2 has passed.
    // In other words erin  and alice told the truth, and the game can now call start
    // resolving that erin earns dave's stake, and alice earns bob's stake. Also both
    // dave and bob will be slashed

    aliceBalance = Number(await web3.eth.getBalance(alice));
    erinBalance = Number(await web3.eth.getBalance(erin));
    await expect.passes(
      stakers.start({from:game}),
      "failed to start new verse from L4"
    )
    assert.isBelow(aliceBalance, Number(await web3.eth.getBalance(alice)),
                 "Alice's current balance should be higher now, since she earned Bob's stake");
    assert.isBelow(erinBalance, Number(await web3.eth.getBalance(erin)),
                 "Erin's current balance should be higher now, since she earned Dave's stake");

    await expect.reverts(
      stakers.enroll({from:bob, value: stake}),
      "failed to enroll: staker was slashed",
      "bob was slashed and will never be able to enroll again, so it should revert"
    )
    await expect.reverts(
      stakers.enroll({from:dave, value: stake}),
      "failed to enroll: staker was slashed",
      "dave was slashed and will never be able to enroll again, so it should revert"
    )
  })

////////////////////////////////////////////////////////////////////////////////////////////

  it("Tests L0 -> L1 -> L2 -> L3 -> L4 -> L3 -> L4", async () => {

    // start (L0) ->  bob updates (L1) -> alice updates (L2) -> carol updates (L3) -> dave updates (L4) -> erin challenges dave (L3) -> erin updates (L4)
    stakers.setGame(game, {from:owner}),
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

    //await expect.reverts(
    //  stakers.start({from:game}),
    //  "failed to start: wrong level",
    //  "starting a new verse is only possible from level 1 or 2, so it should revert"
    //)
  })
})

////////////////////////////////////////////////////////////////////////////////////////////

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
