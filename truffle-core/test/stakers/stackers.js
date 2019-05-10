/* global artifacts */
/* global contract */
/* global web3 */
/* global assert */

const Stakers = artifacts.require("Stakers")
const timetravel = require("./helpers/timetravel");
const truffleAssert = require('truffle-assertions');

const UNENROLLED       = 0;
const ENROLLING        = 1;
const UNENROLLING      = 2;
const UNENROLLABLE     = 3;
const ENROLLED         = 4;
const CHALLENGE_TT     = 5;
const CHALLENGE_LI_RES = 6;
const CHALLENGE_TT_RES = 7;
const SLASHABLE        = 8;
const SLASHED          = 9;

const ERR_BADSTATE  = "err-state";
const ERR_BADHASH   = "err-hash";
const ERR_BADHFIN   = "err-hashfin";
const ERR_BADSTAKE  = "err-stake";
const ERR_POSTCOND  = "err-postcon";
const ERR_BADSENDER = "err-sender";

contract("stakers", (accounts) => {

    const onion0 = web3.utils.keccak256("hel24"); // finishes 2
    const onion1 = web3.utils.keccak256(onion0);  // finishes 0
    const onion2 = web3.utils.keccak256(onion1);  // finishes b
    const onion3 = web3.utils.keccak256(onion2);  // finishes 1

    const {
        0: staker,
        1: game,
        2: updater,
    } = accounts;
    let stakers;

    beforeEach(async () => {
        stakers = await Stakers.new(game, {from : game});
    });

    it("the happy path", async () => {

        const initialBalance = await web3.eth.getBalance(staker);
        const stake = await stakers.REQUIRED_STAKE()

        // check not enrolled initially
        assert.equal(UNENROLLED,await stakers.state(staker,0));

        // enroll
        await truffleAssert.passes(
            stakers.enroll(onion3,{from : staker, value: stake }),
            "failed to enroll"
        );
        assert.isTrue(initialBalance > await web3.eth.getBalance(staker));
        assert.equal(ENROLLING,await stakers.state(staker,0));

        // wait enroll blocks
        await timetravel((await stakers.MINENROLL_SECS()).toNumber());
        assert.equal(ENROLLED,await stakers.state(staker,0));

        // game updates
        await stakers.initChallenge(staker,{from : game});
        assert.equal(CHALLENGE_TT,await stakers.state(staker,0));

        // staker waits for somebody to challenge
        await timetravel((await stakers.MAXCHALL_SECS()).toNumber());
        assert.equal(CHALLENGE_TT_RES,await stakers.state(staker,0));

        // staker resolves
        await truffleAssert.passes(
            stakers.resolveChallenge(onion2,{from : staker}),
            "Failed to resolve challenge"
            );
        assert.equal(ENROLLED,await stakers.state(staker,0));

        // query to unenroll
        await truffleAssert.passes(
            stakers.queryUnenroll({from : staker}),
            "Failed to queryUnenroll"
        );
        assert.equal(UNENROLLING,await stakers.state(staker,0));

        // wait to unenroll time
        balanceBeforeUnenroll = await web3.eth.getBalance(staker).then(web3.utils.toBN);
        await timetravel((await stakers.MINUNENROLL_SECS()).toNumber());
        assert.equal(UNENROLLABLE,await stakers.state(staker,0));
        await truffleAssert.passes(
            stakers.unenroll({from : staker}),
            "Failed to unenroll"
        );

        // check that money has been returned
        balanceAfterUnenroll = await web3.eth.getBalance(staker).then(web3.utils.toBN);
        returned = Number(web3.utils.fromWei(balanceAfterUnenroll.sub(balanceBeforeUnenroll)));
        deposit = Number(web3.utils.fromWei(stake));
        assert.isTrue(initialBalance > await web3.eth.getBalance(staker), "Some gas should have been wasted");
        assert.closeTo(deposit, returned,  0.0005, "The stake has not been returned")
        assert.equal(UNENROLLED,await stakers.state(staker,0));
    });

    it("challenged", async () => {
        await truffleAssert.passes(
            stakers.enroll(onion2,{from : staker, value: await stakers.REQUIRED_STAKE() }),
            "Failed to enroll"
        );
        await timetravel((await stakers.MINENROLL_SECS()).toNumber());
        await truffleAssert.passes(
            stakers.initChallenge(staker,{from : game}),
            "Failed to start challenge"
        );

        // game challenges
        await truffleAssert.passes(
            stakers.lierChallenge(staker,{from : game}),
            "Failed to challenge lier"
        );
        assert.equal(CHALLENGE_LI_RES,await stakers.state(staker,0));

        // staker resolves
        await truffleAssert.passes(
            stakers.resolveChallenge(onion1,{from : staker}),
            "Failed to resolve challenge"
        );
        assert.equal(ENROLLED,await stakers.state(staker,0));
    });

    it("idle enrolled staker", async () => {
        await truffleAssert.passes(
            stakers.enroll(onion2,{from : staker, value: await stakers.REQUIRED_STAKE() }),
            "Failed to enroll"
        );
        await timetravel((await stakers.MINENROLL_SECS()).toNumber());
        await truffleAssert.passes(
            stakers.initChallenge(staker,{from : game}),
            "Failed to start challenge"
        );
        await timetravel((await stakers.MAXCHALL_SECS()).toNumber());

        // go idle
        await timetravel((await stakers.MAXIDLE_SECS()).toNumber());
        assert.equal(SLASHABLE,await stakers.state(staker,0));

        // slash
        await truffleAssert.passes(
            stakers.slash(staker,{from : updater}),
            "Failed to slash"
        );
        assert.equal(SLASHED,await stakers.state(staker,0));
    });

    it("touch idle enrolled staker", async () => {
        await truffleAssert.passes(
            stakers.enroll(onion2,{from : staker, value: await stakers.REQUIRED_STAKE() }),
            "Failed to enroll"
        );
        await timetravel((await stakers.MINENROLL_SECS()).toNumber());

        // go idle
        await timetravel((await stakers.MAXIDLE_SECS()).toNumber());
        assert.equal(SLASHABLE,await stakers.state(staker,0));

        // touch
        await truffleAssert.passes(
            stakers.touch({from : staker}),
            "Failed to slash"
        );
        assert.equal(ENROLLED,await stakers.state(staker,0));
    });


    it("slash idle enrolled staker", async () => {
        await truffleAssert.passes(
            stakers.enroll(onion2,{from : staker, value: await stakers.REQUIRED_STAKE() }),
            "Failed to enroll"
        );
        await timetravel((await stakers.MINENROLL_SECS()).toNumber());

        // go idle
        await timetravel((await stakers.MAXIDLE_SECS()).toNumber());
        assert.equal(SLASHABLE,await stakers.state(staker,0));

        // slash
        await truffleAssert.passes(
            stakers.slash(staker,{from : updater}),
            "Failed to slash"
        );
        assert.equal(SLASHED,await stakers.state(staker,0));
    });

    it("slash idle trueteller staker", async () => {
        await truffleAssert.passes(
            stakers.enroll(onion2,{from : staker, value: await stakers.REQUIRED_STAKE() }),
            "Failed to enroll"
        )
        await timetravel((await stakers.MINENROLL_SECS()).toNumber());
        await truffleAssert.passes(
            stakers.initChallenge(staker,{from : game}),
            "Failed to start challenge"
        );
        await timetravel((await stakers.MAXCHALL_SECS()).toNumber());

        // go idle
        await timetravel((await stakers.MAXIDLE_SECS()).toNumber());
        assert.equal(SLASHABLE,await stakers.state(staker,0));

        // slash
        await truffleAssert.passes(
            stakers.slash(staker,{from : updater}),
            "Failed to slash"
        );
        assert.equal(SLASHED,await stakers.state(staker,0));

        // check that once slashed everything else will fail
        await truffleAssert.reverts(
            stakers.enroll(onion3,{from : staker, value: await stakers.REQUIRED_STAKE() }),
            ERR_BADSTATE,
        );
        await truffleAssert.reverts(stakers.queryUnenroll({from : staker}),            ERR_BADSTATE);
        await truffleAssert.reverts(stakers.unenroll({from : staker}),                 ERR_BADSTATE);
        await truffleAssert.reverts(stakers.touch({from : staker}),                    ERR_BADSTATE);
        await truffleAssert.reverts(stakers.slash(staker, {from : updater}),           ERR_BADSTATE);
        await truffleAssert.reverts(stakers.resolveChallenge(onion2, {from : staker}), ERR_BADSTATE);
        await truffleAssert.reverts(stakers.initChallenge(staker,{from : game}),       ERR_BADSTATE);
        await truffleAssert.reverts(stakers.lierChallenge(staker,{from : game}),       ERR_BADSTATE);
    });

    it("slash idle lier staker", async () => {
        await truffleAssert.passes(
            stakers.enroll(onion2,{from : staker, value: await stakers.REQUIRED_STAKE() }),
            "Failed to enroll"
        );
        await timetravel((await stakers.MINENROLL_SECS()).toNumber());
        await truffleAssert.passes(
            stakers.initChallenge(staker,{from : game}),
            "Failed to start challenge"
        );
        await truffleAssert.passes(
            stakers.lierChallenge(staker,{from : game}),
            "Failed to challenge lier"
        );

        // go idle
        await timetravel((await stakers.MAXIDLE_SECS()).toNumber());
        assert.equal(SLASHABLE,await stakers.state(staker,0));

        // slash
        await truffleAssert.passes(
            stakers.slash(staker,{from : updater}),
            "Failed to slash"
        );
        assert.equal(SLASHED,await stakers.state(staker,0));
    });

    it("touch slashable idle trueteller staker", async () => {
        await truffleAssert.passes(
            stakers.enroll(onion2,{from : staker, value: await stakers.REQUIRED_STAKE() }),
            "Failed to enroll"
        );
        await timetravel((await stakers.MINENROLL_SECS()).toNumber());
        await truffleAssert.passes(
            stakers.initChallenge(staker,{from : game}),
            "Failed to start challenge"
        );
        await timetravel((await stakers.MAXCHALL_SECS()).toNumber());

        // go idle
        await timetravel((await stakers.MAXIDLE_SECS()).toNumber());
        assert.equal(SLASHABLE,await stakers.state(staker,0));

        // slash
        await truffleAssert.passes(
            stakers.touch({from : staker}),
            "Failed to touch"
        );
        assert.equal(ENROLLED,await stakers.state(staker,0));
    });

    it("Tests ERR_BADSTATE in enroll/unenroll", async () => {

        const initialBalance = await web3.eth.getBalance(staker).then(web3.utils.toBN);
        const stake = await stakers.REQUIRED_STAKE();
        const min_enroll_sec = (await stakers.MINENROLL_SECS()).toNumber();
        const max_idle_sec = (await stakers.MAXIDLE_SECS()).toNumber();
        const min_unenroll_sec = (await stakers.MINUNENROLL_SECS()).toNumber();

        var checkFailures = async (msg) => {
            // check any call will fail with ERR_BADSTATE
            truffleAssert.reverts(stakers.unenroll({from : staker}),                 ERR_BADSTATE, msg);
            truffleAssert.reverts(stakers.touch({from : staker}),                    ERR_BADSTATE, msg);
            truffleAssert.reverts(stakers.slash(staker, {from : updater}),           ERR_BADSTATE, msg);
            truffleAssert.reverts(stakers.resolveChallenge(onion2, {from : staker}), ERR_BADSTATE, msg);
            truffleAssert.reverts(stakers.initChallenge(staker,{from : game}),       ERR_BADSTATE, msg);
            truffleAssert.reverts(stakers.lierChallenge(staker,{from : game}),       ERR_BADSTATE, msg);
        };

        var lookup = async (seconds) => {
            touchedTimeStamp = (await stakers.stakers(staker)).touch.toNumber();
            return touchedTimeStamp /*- (await web3.eth.getBlock('latest')).timestamp */ + seconds;
        };

        // check staker is not enrolled initially
        assert.equal(UNENROLLED,await stakers.state(staker,0));
        await checkFailures("before enrolling");


        await truffleAssert.passes(
            stakers.enroll(onion3,{from : staker, value: stake}),
            "failed to enroll"
        );

        const currentBalance = await web3.eth.getBalance(staker).then(web3.utils.toBN);
        assert.closeTo(
            Number(web3.utils.fromWei(initialBalance.sub(currentBalance).sub(stake))),
            0,
            0.003,
            "The stake was not returned"
        );

        await truffleAssert.reverts(
            stakers.enroll(onion3,{from : staker, value: stake}),
            ERR_BADSTATE,
            "Re-enrolling should fail"
        );

        assert.equal(ENROLLING,await stakers.state(staker,0));
        await checkFailures("while enrolling");

        assert.equal(ENROLLED, await stakers.state(staker, await lookup(min_enroll_sec)));
        assert.equal(SLASHABLE, await stakers.state(staker, await lookup(min_enroll_sec + max_idle_sec)));

        await truffleAssert.passes(stakers.queryUnenroll({from : staker}));
        assert.equal(UNENROLLING, await stakers.state(staker,0));
        await checkFailures("while un-enrolling");

        touchedTimeStamp = (await stakers.stakers(staker)).touch.toNumber(),
        assert.equal(UNENROLLABLE, await stakers.state(staker, await lookup(min_unenroll_sec)));
        await checkFailures("while un-enrollable");

        await timetravel((await stakers.MINUNENROLL_SECS()).toNumber());
        assert.equal(UNENROLLABLE, await stakers.state(staker,0));
        await truffleAssert.passes(stakers.unenroll({from : staker}));
        assert.equal(UNENROLLED,await stakers.state(staker,0));
    });

    it("Tests ERR_BADSENDER", async () => {
        await truffleAssert.reverts(stakers.initChallenge(staker,{from : staker}), ERR_BADSENDER);
        await truffleAssert.reverts(stakers.lierChallenge(staker,{from : staker}), ERR_BADSENDER);
    });

    it("Tests ERR_BADSTAKE", async () => {

        const stake = await stakers.REQUIRED_STAKE();
        var unit = new web3.utils.BN(1, 16);
        await truffleAssert.reverts(stakers.enroll(onion3,{from : staker, value: 0}), ERR_BADSTAKE);
        await truffleAssert.reverts(stakers.enroll(onion3,{from : staker, value: stake.sub(unit)}), ERR_BADSTAKE);
    });

    it("Tests ERR_BADSTATE and ERR_BADHASH from CHALLLENGE_TT and CHALLENGE_TT_RES", async () => {

        const initialBalance = await web3.eth.getBalance(staker);
        const stake = await stakers.REQUIRED_STAKE()

        // check not enrolled initially
        assert.equal(UNENROLLED,await stakers.state(staker,0));

        // enroll
        await truffleAssert.passes(
            stakers.enroll(onion3,{from : staker, value: stake }),
            "failed to enroll"
        );
        assert.isTrue(initialBalance > await web3.eth.getBalance(staker));
        assert.equal(ENROLLING,await stakers.state(staker,0));

        // wait enroll blocks
        await timetravel((await stakers.MINENROLL_SECS()).toNumber());
        assert.equal(ENROLLED,await stakers.state(staker,0));

        // game updates
        await stakers.initChallenge(staker,{from : game});
        assert.equal(CHALLENGE_TT,await stakers.state(staker,0));

        // check any call will fail with ERR_BADSTATE
        await truffleAssert.reverts(stakers.queryUnenroll({from : staker}),            ERR_BADSTATE);
        await truffleAssert.reverts(stakers.unenroll({from : staker}),                 ERR_BADSTATE);
        await truffleAssert.reverts(stakers.touch({from : staker}),                    ERR_BADSTATE);
        await truffleAssert.reverts(stakers.slash(staker, {from : updater}),           ERR_BADSTATE);
        await truffleAssert.reverts(stakers.resolveChallenge(onion2, {from : staker}), ERR_BADSTATE);
        await truffleAssert.reverts(stakers.initChallenge(staker,{from : game}),       ERR_BADSTATE);

        await timetravel((await stakers.MAXCHALL_SECS()).toNumber());
        assert.equal(CHALLENGE_TT_RES,await stakers.state(staker,0));

        // check any call will fail with ERR_BADSTATE
        await truffleAssert.reverts(stakers.queryUnenroll({from : staker}),      ERR_BADSTATE);
        await truffleAssert.reverts(stakers.unenroll({from : staker}),           ERR_BADSTATE);
        await truffleAssert.reverts(stakers.touch({from : staker}),              ERR_BADSTATE);
        await truffleAssert.reverts(stakers.slash(staker, {from : updater}),     ERR_BADSTATE);
        await truffleAssert.reverts(stakers.initChallenge(staker,{from : game}), ERR_BADSTATE);
        await truffleAssert.reverts(stakers.lierChallenge(staker,{from : game}), ERR_BADSTATE);

        // tests ERR_BADHASH
        await truffleAssert.reverts(stakers.resolveChallenge(onion3, {from : staker}), ERR_BADHASH);
        await truffleAssert.reverts(stakers.resolveChallenge(onion1, {from : staker}), ERR_BADHASH);
        await truffleAssert.reverts(stakers.resolveChallenge(onion0, {from : staker}), ERR_BADHASH);

        await truffleAssert.passes(
            stakers.resolveChallenge(onion2, {from : staker}),
            "Failed to resolve challenge using the correct hash onion"
        );
    });

    it("Tests ERR_BADSTATE, ERR_BADHFIN and ERR_BADHASH from CHALLENGE_LI_RES", async () => {

        const initialBalance = await web3.eth.getBalance(staker);
        const stake = await stakers.REQUIRED_STAKE()

        // check not enrolled initially
        assert.equal(UNENROLLED,await stakers.state(staker,0));

        // enroll
        await truffleAssert.passes(
            stakers.enroll(onion3,{from : staker, value: stake }),
            "failed to enroll"
        );
        assert.isTrue(initialBalance > await web3.eth.getBalance(staker));
        assert.equal(ENROLLING,await stakers.state(staker,0));

        // wait enroll blocks
        await timetravel((await stakers.MINENROLL_SECS()).toNumber());
        assert.equal(ENROLLED,await stakers.state(staker,0));

        // game updates
        await stakers.initChallenge(staker,{from : game});
        assert.equal(CHALLENGE_TT,await stakers.state(staker,0));
        await truffleAssert.passes(stakers.lierChallenge(staker,{from : game}));
        assert.equal(CHALLENGE_LI_RES,await stakers.state(staker,0));

        // check any call will fail with ERR_BADSTATE
        await truffleAssert.reverts(stakers.queryUnenroll({from : staker}),      ERR_BADSTATE);
        await truffleAssert.reverts(stakers.unenroll({from : staker}),           ERR_BADSTATE);
        await truffleAssert.reverts(stakers.touch({from : staker}),              ERR_BADSTATE);
        await truffleAssert.reverts(stakers.slash(staker, {from : updater}),     ERR_BADSTATE);
        await truffleAssert.reverts(stakers.initChallenge(staker,{from : game}), ERR_BADSTATE);

        await truffleAssert.reverts(
            stakers.resolveChallenge(onion2, {from : staker}),
            ERR_BADHFIN,
            "Should fail because onion % 16 != 0"
        );
        await truffleAssert.reverts(
            stakers.resolveChallenge(onion1, {from : staker}),
            ERR_BADHASH,
            "Should fail because incorrect onion layer"
        );
    });
});
