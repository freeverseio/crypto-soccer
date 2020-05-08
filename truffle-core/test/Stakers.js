const Stakers = artifacts.require("Stakers")
const expect = require('truffle-assertions');

// TODO: add more tests that execute reward

contract('Stakers', (accounts) => {
  const [owner, gameAddr, alice, bob, carol, dave, erin, frank] = accounts

  let stakers
  let stake

  beforeEach(async () => {
      stakers  = await Stakers.new({from:owner});
      stake = await stakers.kRequiredStake();
  });

////////////////////////////////////////////////////////////////////////////////////////////

  it("Tests game address", async () => {
    await expect.reverts(
      stakers.update(level = 1, alice),
      "Only gameOwner can call this function",
      "game not set yet, so it should revert"
    )
    await expect.reverts(
      stakers.setGameOwner(gameAddr, {from:alice}),
      "Only owner can call this function",
      "wrong sender, so it should revert"
    )
    await expect.passes(
      stakers.setGameOwner(gameAddr, {from:owner}),
      "failed to set game address"
    )
  })

  it("Tests owner address change", async () => {
    await expect.reverts(
      stakers.setOwner(alice, {from:alice}),
      "Only owner can call this function",
      "wrong sender, so it should revert"
    )
    await expect.passes(
      stakers.setOwner(alice, {from:owner}),
      "failed to set new owner address"
    )
    await expect.reverts(
      stakers.setGameOwner(gameAddr, {from:owner}),
      "Only owner can call this function",
      "owner is not true owner anymore, so it should revert"
    )
    await expect.passes(
      stakers.setGameOwner(gameAddr, {from:alice}),
      "failed to set game address by updated owner"
    )
  })

////////////////////////////////////////////////////////////////////////////////////////////

  it("Tests enrolling", async () => {
    await expect.reverts(
      stakers.enroll({from:alice, value: stake}),
      null,
      "alice is not yet a trusted party, so it should revert"
    )
    await expect.reverts(
      stakers.addTrustedParty(alice, {from:gameAddr}),
      null,
      "only owner can add trusted parties, so it should revert"
    )
    await expect.passes(
      stakers.addTrustedParty(alice, {from:owner}),
      "failed to add alice as trusted party"
    )
    await expect.reverts(
      stakers.addTrustedParty(alice, {from:owner}),
      null,
      "alice is already a trusted party, so it should revert"
    )
    await expect.passes(
      stakers.enroll({from:alice, value: stake}),
      "failed to enroll alice"
    )
  });

////////////////////////////////////////////////////////////////////////////////////////////

  it("Tests stake", async () => {
    assert.equal(0, await web3.eth.getBalance(stakers.address), "Initial contract should have 0 stake in place");

    parties = [alice, bob, carol, dave, erin, frank]
    await addTrustedParties(stakers, owner, parties);
    await enroll(stakers, stake, parties);

    assert.equal(parties.length*Number(stake),
      await web3.eth.getBalance(stakers.address),
      "Total stake is not the sum of all enrolled stakers stake"
    );

    await unenroll(stakers, parties);
    assert.equal(0, await web3.eth.getBalance(stakers.address));
  });

////////////////////////////////////////////////////////////////////////////////////////////

  it("Tests stakers can't unenroll after having done an update", async () => {
    stakers.setGameOwner(gameAddr, {from:owner}),
    await stakers.addTrustedParty(alice, {from:owner});
    await stakers.enroll({from:alice, value: stake});
    await stakers.update(level = 0, alice, {from:gameAddr}),
    await expect.reverts(
      stakers.unEnroll({from:alice}),
      "failed to unenroll: staker currently updating",
      "alice is currently updating, so it should revert"
    )
    await stakers.finalize({from:gameAddr});
    await expect.passes(
      stakers.unEnroll({from:alice}),
      "failed unenrolling alice"
    )
  })

////////////////////////////////////////////////////////////////////////////////////////////

  it("Tests adding reward", async () => {
    assert.equal(0, Number(await web3.eth.getBalance(await stakers.rewards())));
    await expect.passes(
      stakers.addReward({value: stake}),
      "failed to add reward")
    assert.equal(Number(stake), Number(await web3.eth.getBalance(await stakers.rewards())));
    assert.equal(Number(0), Number(await web3.eth.getBalance(stakers.address)));
    await expect.reverts(
      stakers.executeReward({from:owner}),
      "failed to execute rewards: empty array",
      "no one deserves reward cause nothing has been played, so it should revert"
    )
  })

////////////////////////////////////////////////////////////////////////////////////////////

  it("Tests L0 -> L1 true -> start -> L1 true, the usual path", async () => {

    stakers.setGameOwner(gameAddr, {from:owner}),
    parties = [alice, bob, carol, dave, erin, frank]
    await addTrustedParties(stakers, owner, parties);
    await enroll(stakers, stake, parties);
    await expect.passes(
      stakers.addReward({value: stake}),
      "failed to add reward")

    // L0
    assert.equal(0, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(level = 0, alice, {from:gameAddr}),
      "alice failed to update"
    )

    // L1
    assert.equal(1, (await stakers.level()).toNumber());
    await expect.reverts(
      stakers.update(level = 0, bob, {from:gameAddr}),
      "failed to update: resolving wrong level",
      "level 0 cannot be updated without starting a new verse, it should revert"
    )

    for (i=0; i<10; i++) {
      // start new verse
      await expect.passes(
        stakers.finalize({from:gameAddr}),
        "failed starting new verse"
      )

      // L0
      assert.equal(0, (await stakers.level()).toNumber());
      await expect.passes(
        stakers.update(level = 0, alice, {from:gameAddr}),
        "alice failed to update because he lied in previous game"
      )

      // L1
      assert.equal(1, (await stakers.level()).toNumber());
    }

    // execute reward and test that alice has more balance
    aliceBalanceBeforeRewarded = Number(await web3.eth.getBalance(alice));
    await expect.passes(
      stakers.executeReward({from:owner}),
      "failed to execute reward"
    )

    await expect.passes(
      stakers.withdraw({from:alice}),
      "failed to withdraw alice's reward"
    )

    assert.isBelow(aliceBalanceBeforeRewarded, Number(await web3.eth.getBalance(alice)),
                 "Alice's current balance should be higher since she was rewarded");
  })

////////////////////////////////////////////////////////////////////////////////////////////

  it("Tests L0 -> L1 lie  -> L2 true -> start -> L1 lie  -> L2 true", async () => {

    stakers.setGameOwner(gameAddr, {from:owner}),
    parties = [alice, bob, carol, dave, erin, frank]
    await addTrustedParties(stakers, owner, parties);
    await enroll(stakers, stake, parties);

    // L0: first updater lies
    assert.equal(0, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(level = 0, alice, {from:gameAddr}),
      "alice failed to update"
    )

    // L1: second updater (1st challenger) tells truth
    assert.equal(1, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(level = 1, bob, {from:gameAddr}),
      "bob failed to update"
    )

    // Ensure that nobody can update lev = 0 nor lev = 1 again, 
    // since we are at level = 2
    assert.equal(2, (await stakers.level()).toNumber());
    await expect.reverts(
      stakers.update(level = 0, carol, {from:gameAddr}),
      "failed to update: resolving wrong level",
      "level 0 cannot be updated without starting a new verse, it should revert"
    )
    await expect.reverts(
      stakers.update(1, carol, {from:gameAddr}),
      "failed to update: resolving wrong level",
      "level 1 is already updated, it should revert"
    )

    // ------------- start new verse ----------------
    // which implicitly slashes alice, the first updater
    await expect.passes(
      stakers.finalize({from:gameAddr}),
      "failed starting new verse"
    )

    // L0: check that Alice is not registered as staker anymore
    assert.equal(0, (await stakers.level()).toNumber());
    await expect.reverts(
      stakers.update(0, alice, {from:gameAddr}),
      "failed to update: staker not registered",
      "alice was slashed by bob and therefore it is removed from registered stakers, so it should revert"
    )
    // check that Alice cannot enroll again
    await expect.reverts(
      stakers.enroll({from:alice, value: stake}),
      "failed to enroll: staker was slashed",
      "alice was slashed by bob it can no longer enroll, so it should revert"
    )
    await expect.passes(
      stakers.update(0, bob, {from:gameAddr}),
      "bob failed to update after new verse"
    )

    // L1
    assert.equal(1, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(1, dave, {from:gameAddr}),
      "dave failed to update level 1"
    )

    // L2
    assert.equal(2, (await stakers.level()).toNumber());
  })

////////////////////////////////////////////////////////////////////////////////////////////
  it("Tests L0 -> L1 true -> L2 lie  -> L3 true -> start", async () => {

    stakers.setGameOwner(gameAddr, {from:owner}),
    parties = [alice, bob, carol, dave, erin, frank]
    await addTrustedParties(stakers, owner, parties);
    await enroll(stakers, stake, parties);

    // L0
    assert.equal(0, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(0, alice, {from:gameAddr}),
      "alice failed to update"
    )

    // L1
    assert.equal(1, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(1, bob, {from:gameAddr}),
      "bob failed to update"
    )

    // L2
    assert.equal(2, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(2, dave, {from:gameAddr}),
      "dave failed to update"
    )

    // L3
    assert.equal(3, (await stakers.level()).toNumber());

    // make sure level 2 is not updatable
    await expect.reverts(
      stakers.update(2, erin, {from:gameAddr}),
      "failed to update: resolving wrong level",
      "we are at level 3, level 2 is not updatable anymore, so it should revert"
    )

    // challenge time for L3 has passed, and also challenge time for L1 has passed.
    // In other words dave  and alice told the truth, and the game can now call start
    // resolving that dave earns bob's stake, and alice earns reward. Also bob will
    // be slashed

    daveBalance = Number(await web3.eth.getBalance(dave));
    await expect.passes(
      stakers.finalize({from:gameAddr}),
      "failed to start from L3"
    )
    assert.isBelow(daveBalance, Number(await web3.eth.getBalance(dave)),
                 "Dave's current balance should be higher now, since he earned bob's stake");

  })


  it("Tests L0 -> L1 lie  -> L2 lie  -> L3 true -> L1 -> L2 true", async () => {

    stakers.setGameOwner(gameAddr, {from:owner}),
    parties = [alice, bob, carol, dave, erin, frank]
    await addTrustedParties(stakers, owner, parties);
    await enroll(stakers, stake, parties);

    // L0
    assert.equal(0, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(0, alice, {from:gameAddr}),
      "alice failed to update"
    )

    // L1
    assert.equal(1, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(1, bob, {from:gameAddr}),
      "bob failed to update"
    )

    // L2
    assert.equal(2, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(2, dave, {from:gameAddr}),
      "dave failed to update"
    )
    daveBalance = Number(await web3.eth.getBalance(dave));

    // L3
    assert.equal(3, (await stakers.level()).toNumber());

    // challenge time has passed, resolve from L3: bob will be slashed and dave earns bob's stake
    await expect.passes(
      stakers.update(1, erin, {from:gameAddr}),
      "erin failed to update"
    )

    // L2: because L3 was resolved with an update to L1, we are now at L2
    assert.equal(2, (await stakers.level()).toNumber());

    await expect.reverts(
      stakers.enroll({from:bob, value: stake}),
      "failed to enroll: staker was slashed",
      "bob was slashed, so it should revert"
    )

    assert.isBelow(daveBalance, Number(await web3.eth.getBalance(dave)),
                 "Dave current balance should be higher now, since he earned bob's stake");

    // start new verse
    await expect.passes(
      stakers.finalize({from:gameAddr}),
      "failed starting new verse"
    )

    await expect.reverts(
      stakers.enroll({from:bob, value: stake}),
      "failed to enroll: staker was slashed",
      "bob was slashed and will never be able to enroll again, so it should revert"
    )
  })

////////////////////////////////////////////////////////////////////////////////////////////

  it("Tests L0 -> L1 true -> L2 lie  -> L3 lie  -> L4 true -> L2 -> L3 true -> start", async () => {
    stakers.setGameOwner(gameAddr, {from:owner}),
    parties = [alice, bob, carol, dave, erin, frank]
    await addTrustedParties(stakers, owner, parties);
    await enroll(stakers, stake, parties);

    // L0
    assert.equal(0, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(0, alice, {from:gameAddr}),
      "alice failed to update"
    )

    // L1 - true
    assert.equal(1, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(1, bob, {from:gameAddr}),
      "bob failed to update"
    )

    // L2 - lie
    assert.equal(2, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(2, dave, {from:gameAddr}),
      "dave failed to update"
    )

    // L3 - lie
    assert.equal(3, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(3, erin, {from:gameAddr}),
      "erin failed to update to L4"
    )
    erinBalance = Number(await web3.eth.getBalance(erin));

    // L4 - true
    assert.equal(4, (await stakers.level()).toNumber());

    // challenge time passed, resolve from L4: erin told the truth,  dave will be slashed and erin earns dave's stake
    await expect.passes(
      stakers.update(2, frank, {from:gameAddr}),
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
    // In other words frank  and alice told the truth, and the game can now call start
    // resolving that frank earns bob's stake, and alice earns reward. Also bob will
    // be slashed

    frankBalance = Number(await web3.eth.getBalance(frank));
    await expect.passes(
      stakers.finalize({from:gameAddr}),
      "failed calling start from L3"
    )

    // TODO: this fails?
    assert.isBelow(frankBalance, Number(await web3.eth.getBalance(frank)),
                 "Frank's current balance should be higher now, since he earned bob's stake");

    await expect.reverts(
      stakers.enroll({from:bob, value: stake}),
      "failed to enroll: staker was slashed",
      "bob was slashed and will never be able to enroll again, so it should revert"
    )
  })

////////////////////////////////////////////////////////////////////////////////////////////

  it("Tests L0 -> L1 lie  -> L2 true -> L3 lie  -> L4 true -> start", async () => {
    stakers.setGameOwner(gameAddr, {from:owner})
    parties = [alice, bob, carol, dave, erin, frank]
    await addTrustedParties(stakers, owner, parties);
    await enroll(stakers, stake, parties);
    // L0
    assert.equal(0, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(0, alice, {from:gameAddr}),
      "alice failed to update"
    )

    // L1 - lie
    assert.equal(1, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(1, bob, {from:gameAddr}),
      "bob failed to update"
    )

    // L2 - true
    assert.equal(2, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(2, dave, {from:gameAddr}),
      "dave failed to update"
    )

    // L3 - lie
    assert.equal(3, (await stakers.level()).toNumber());
    await expect.passes(
      stakers.update(3, erin, {from:gameAddr}),
      "erin failed to update to L4"
    )

    // L4 - true
    assert.equal(4, (await stakers.level()).toNumber());

    // make sure level 3 is not updatable
    await expect.reverts(
      stakers.update(3, frank, {from:gameAddr}),
      "failed to update: resolving wrong level",
      "we are at level 4, level 3 is not updatable anymore, so it should revert"
    )

    // challenge time for L4 has passed, and also challenge time for L2 has passed.
    // In other words erin  and bob told the truth, and the game can now call start
    // resolving that erin earns dave's stake, and bob earns alice's stake. Also both
    // dave and alice will be slashed

    bobBalance = Number(await web3.eth.getBalance(bob));
    erinBalance = Number(await web3.eth.getBalance(erin));
    await expect.passes(
      stakers.finalize({from:gameAddr}),
      "failed to start new verse from L4"
    )
    assert.isBelow(bobBalance, Number(await web3.eth.getBalance(bob)),
                 "bob's current balance should be higher now, since she earned alice's stake");
    assert.isBelow(erinBalance, Number(await web3.eth.getBalance(erin)),
                 "Erin's current balance should be higher now, since she earned Dave's stake");

    await expect.reverts(
      stakers.enroll({from:alice, value: stake}),
      "failed to enroll: staker was slashed",
      "alice was slashed and will never be able to enroll again, so it should revert"
    )
    await expect.reverts(
      stakers.enroll({from:dave, value: stake}),
      "failed to enroll: staker was slashed",
      "dave was slashed and will never be able to enroll again, so it should revert"
    )
  })

////////////////////////////////////////////////////////////////////////////////////////////

  it("Tests L0 -> L1 lie  -> L2 lie  -> L3 true  -> L4 lie -> challenge -> L3", async () => {

    // start (L0) ->  alice updates (L1) -> bob updates (L2) -> carol updates (L3) -> dave updates (L4) -> erin challenges dave (L3) -> erin updates (L4)
    stakers.setGameOwner(gameAddr, {from:owner}),
    parties = [alice, bob, carol, dave, erin, frank]
    await addTrustedParties(stakers, owner, parties);
    await enroll(stakers, stake, parties);

    assert.equal(0, (await stakers.level()).toNumber());

    await expect.reverts(
      stakers.update(1, alice, {from:gameAddr}),
      "failed to update: wrong level",
      "level to update is 1, so it should revert"
    )

    await expect.reverts(
      stakers.update(2, alice, {from:gameAddr}),
      "failed to update: wrong level",
      "level to update is 1, so it should revert"
    )

    await expect.passes(
      stakers.update(0, alice, {from:gameAddr}),
      "alice failed to update"
    )

    assert.equal(1, (await stakers.level()).toNumber());

    await expect.reverts(
      stakers.update(1, alice, {from:gameAddr}),
      null,
      "alice is already updating, cannot participate until resolved"
    )

    await expect.passes(
      stakers.update(1, bob, {from:gameAddr}),
      "bob failed to update"
    )

    assert.equal(2, (await stakers.level()).toNumber());

    await expect.passes(
      stakers.update(2, carol, {from:gameAddr}),
      "carol failed to update"
    )

    assert.equal(3, (await stakers.level()).toNumber());

    await expect.passes(
      stakers.update(3, dave, {from:gameAddr}),
      "dave failed to update"
    )

    assert.equal(4, (await stakers.level()).toNumber());

    // erin callenges the very last update L4
    await expect.passes(
      stakers.update(4, erin, {from:gameAddr}),
      "erin failed to challenge"
    )

    assert.equal(3, (await stakers.level()).toNumber());

    await expect.reverts(
      stakers.update(5, erin, {from:gameAddr}),
      "failed to update: wrong level",
      "this update level does not exist, so it should revert"
    )

    assert.equal(3, (await stakers.level()).toNumber());

    await expect.reverts(
      stakers.update(4, erin, {from:gameAddr}),
      "failed to update: wrong level",
      "after erin slashed dave, level is 3, so it should revert"
    )

    assert.equal(3, (await stakers.level()).toNumber());
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
