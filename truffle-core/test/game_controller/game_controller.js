const GameController = artifacts.require("GameController");
const Stakers = artifacts.require("Stakers")
const truffleAssert = require('truffle-assertions');
const util = require('util')

const UNENROLLED       = 0;
const ENROLLING        = 1;
const ENROLLED         = 4;
const ERR_BADSTATE     = "err-state"

const ERR_NO_STAKERS          = "err-no-stakers-contract-set";
const ERR_WINDOW_NOT_STARTED  = "err-window-not-started";
const ERR_WINDOW_FINISHED     = "err-window-finished";
const ERR_WINDOW_RESTRICTED   = "err-window-restricted";

contract('game_controller', (accounts) => {
    const [owner, bob, alice] = accounts

    const onion0 = web3.utils.keccak256("hel24"); // finishes 2
    const onion1 = web3.utils.keccak256(onion0);  // finishes 0
    const onion2 = web3.utils.keccak256(onion1);  // finishes b
    const onion3 = web3.utils.keccak256(onion2);  // finishes 1

    let controller
    let stakers

    beforeEach(async () => {
        controller = await GameController.new()
        stakers = await Stakers.new(controller.address, {from:owner});
        stake = await stakers.REQUIRED_STAKE()

        await truffleAssert.passes(
          stakers.enroll(onion3,{from:bob, value:stake}),
          "failed to enroll"
        );
         assert.equal(ENROLLING,await stakers.state(bob,0));
    });

    it("Tests invalid stakers address", async () => {
        await truffleAssert.reverts(
          controller.updated(1, 0, bob),
          ERR_NO_STAKERS,
          "Stakers address is not valid so it should revert"
        )
        await truffleAssert.reverts(
          controller.challenged(1),
          ERR_NO_STAKERS,
          "Stakers address is not valid so it should revert"
        )
    })

    it("Tests setting stakers address", async () => {
        await truffleAssert.reverts(
          controller.setStakersContractAddress(stakers.address, {from:bob})
        )
        await truffleAssert.passes(
          controller.setStakersContractAddress(stakers.address, {from:owner})
        )
        await truffleAssert.reverts(
          controller.setStakersContractAddress("0x0000000000000000000000000000000000000000", {from:owner})
        )
        assert.equal(
          await controller.getStakersContractAddress(),
          stakers.address,
          "Stakers contract address differs from expected"
        )
    })

    it("Tests updated and challenged", async () => {
        const restrictedPeriod = (await controller.kWindowBlocksRestricted()).toNumber()
        const windowLength = (await controller.kWindowBlocks()).toNumber()

        await truffleAssert.passes(
          controller.setStakersContractAddress(stakers.address, {from:owner})
        )

        latestBlock = await web3.eth.getBlock('latest')
        leagueStartBlock = latestBlock.number
        leagueDuration = 10
        leagueEndBlock = leagueStartBlock + leagueDuration

        leagueId = 1
        windowStart = leagueEndBlock
        windowEveryone = windowStart + restrictedPeriod
        windowEnd  = windowStart + windowLength

        await truffleAssert.reverts(
          controller.updated(leagueId, windowStart, alice),
          ERR_WINDOW_NOT_STARTED,
          "League updated before league duration"
        )

        await jumpSeconds((await stakers.MINENROLL_SECS()).toNumber())
        assert.equal(ENROLLED,await stakers.state(bob,0));

        // jump beyond restricted period
        await jumpBlocks(leagueDuration + restrictedPeriod + 1)
        latestBlock = await web3.eth.getBlock('latest')

        assert.isTrue(
          latestBlock.number > windowEveryone,
          "Everyone should be able to update"
        )
        assert.isTrue(
          latestBlock.number < windowEnd,
          "window should not be ended"
        )
        await truffleAssert.passes(
          controller.updated(leagueId, windowStart, bob),
          "Failed updating league after league duration"
        )
        await truffleAssert.reverts(
          controller.updated(leagueId, windowStart + 10000, bob),
          ERR_WINDOW_NOT_STARTED,
          "Was able to update before window start"
        )
        await truffleAssert.passes(
          controller.challenged(leagueId),
          "Failed to challenge"
        )
        await truffleAssert.reverts(
          controller.challenged(leagueId),
          ERR_BADSTATE,
          "Alice was challenged without being an updater"
        )
    })
})

async function jumpBlocks(n) {
    const sendAsync = util.promisify(web3.currentProvider.send);
    await Promise.all(
      [...Array(n).keys()].map(i =>
        sendAsync({
          jsonrpc: '2.0',
          method: 'evm_mine',
          id: i
        })
      )
    );
}
const jumpSeconds = function(duration) {
  const id = Date.now()
  const sendAsync = util.promisify(web3.currentProvider.send);

  return new Promise((resolve, reject) => {
    sendAsync({
      jsonrpc: '2.0',
      method: 'evm_increaseTime',
      params: [duration],
      id: id,
    }, err1 => {
      if (err1) return reject(err1)

      sendAsync({
        jsonrpc: '2.0',
        method: 'evm_mine',
        id: id+1,
      }, (err2, res) => {
        return err2 ? reject(err2) : resolve(res)
      })
    })
  })
}
