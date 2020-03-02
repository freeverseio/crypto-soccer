const BN = require('bn.js');
require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bn')(BN))
    .should();
const truffleAssert = require('truffle-assertions');
const logUtils = require('../utils/matchLogUtils.js');
const debug = require('../utils/debugUtils.js');

const Utils = artifacts.require('Utils');
const TrainingPoints = artifacts.require('TrainingPoints');
const Evolution = artifacts.require('Evolution');
const Assets = artifacts.require('Assets');
const Market = artifacts.require('Market');
const EncodingMatchLog = artifacts.require('EncodingMatchLog');
const Engine = artifacts.require('Engine');
const EnginePreComp = artifacts.require('EnginePreComp');
const EngineApplyBoosters = artifacts.require('EngineApplyBoosters');
const PlayAndEvolve = artifacts.require('PlayAndEvolve');
const Shop = artifacts.require('Shop');


contract('Evolution', (accounts) => {
    const substitutions = [6, 10, 0];
    const subsRounds = [3, 7, 1];
    const noSubstitutions = [11, 11, 11];
    const lineup0 = [0, 3, 4, 5, 6, 9, 10, 11, 12, 15, 16, 7, 13, 17];
    const lineup1 = [0, 3, 4, 5, 6, 9, 10, 11, 16, 17, 18, 7, 13, 17];
    const lineupConsecutive =  Array.from(new Array(14), (x,i) => i);
    const extraAttackNull =  Array.from(new Array(10), (x,i) => 0);
    const tacticId442 = 0; // 442
    const tacticId433 = 2; // 433
    const playersPerZone442 = [1,2,1,1,2,1,0,2,0];
    const playersPerZone433 = [1,2,1,1,1,1,1,1,1];
    const PLAYERS_PER_TEAM_MAX = 25;
    const firstHalfLog = [0, 0];
    const subLastHalf = false;
    const is2ndHalf = false;
    const isHomeStadium = true;
    const isPlayoff = false;
    const matchBools = [is2ndHalf, isHomeStadium, isPlayoff]
    const IDX_R = 1;
    const IDX_C = 2;
    const IDX_CR = 3;
    const IDX_L = 4;
    const IDX_LR = 5;
    const IDX_LC = 6;
    const IDX_LCR = 7;
    const fwd442 =  [0, 1, 1, 1, 1, 2, 2, 2, 2, 3, 3];
    const left442 = [0, IDX_L, IDX_C, IDX_C, IDX_R, IDX_L, IDX_C, IDX_C, IDX_R, IDX_C, IDX_C];
    // const now = Math.floor(new Date()/1000);
    // const dayOfBirth21 = Math.round(secsToDays(now) - 21/7);
    const now = 1570147200; // this number has the property that 7*nowFake % (SECS_IN_DAY) = 0 and it is basically Oct 3, 2019
    const dayOfBirth21 = secsToDays(now) - 21*365/7; // = exactly 17078, no need to round
    const MAX_PENALTY = 10000;
    const MAX_GOALS = 12;
    const RED_CARD = 3;

    const assistersIdx = Array.from(new Array(MAX_GOALS), (x,i) => i);
    const shootersIdx  = Array.from(new Array(MAX_GOALS), (x,i) => 1);
    const shooterForwardPos  = Array.from(new Array(MAX_GOALS), (x,i) => 1);
    const penalties  = Array.from(new Array(7), (x,i) => false);
    const typesOutOfGames = [0, 0];
    const outOfGameRounds = [0, 0];
    const yellowCardedDidNotFinish1stHalf = [false, false];
    const ingameSubs1 = [0, 0, 0]
    const ingameSubs2 = [0, 0, 0]
    const outOfGames = [14, 14]
    const yellowCards1 = [14, 14]
    const yellowCards2 = [14, 14]
    const halfTimeSubstitutions = [14, 14, 14]
    const nDefs1 = 4; 
    const nDefs2 = 4; 
    const nTot = 11; 
    const winner = 2; // DRAW = 2
    const isHomeSt = false;
    const teamSumSkillsDefault = 0;
    const trainingPointsInit = 0;
    
    // Skills: shoot, speed, pass, defence, endurance
    const SK_SHO = 0;
    const SK_SPE = 1;
    const SK_PAS = 2;
    const SK_DEF = 3;
    const SK_END = 4;
    
    const it2 = async(text, f) => {};

    function setNoSubstInLineUp(lineup, substitutions) {
        modifiedLineup = [...lineup];
        NO_SUBST = 11;
        NO_LINEUP = 25;
        for (s = 0; s < 3; s++) {
            if (substitutions[s] == NO_SUBST) modifiedLineup[s + 11] = NO_LINEUP;
        }
        return modifiedLineup;
    }

    function daysToSecs(dayz) {
        return (dayz * 24 * 3600); 
    }

    function secsToDays(secs) {
        return secs/ (24 * 3600);
    }
    
    function getDefaultTPs() {
        TP = 200;
        TPperSkill = Array.from(new Array(25), (x,i) => TP/5 - 3*i % 6);
        specialPlayer = 21;
        // make sure they sum to TP:
        for (bucket = 0; bucket < 5; bucket++){
            sum4 = 0;
            for (sk = 5 * bucket; sk < (5 * bucket + 4); sk++) {
                sum4 += TPperSkill[sk];
            }
            TPperSkill[5 * bucket + 4] = TP - sum4;
        }       
        return [TP, TPperSkill];
    }
    
    const createTeamState442 = async (engine, forceSkills, alignedEndOfLastHalfTwoVec = [false, false]) => {
        teamState = [];
        playerId = 123121;
        pot = 3;
        aggr = 0;
        alignedEndOfLastHalf = true;
        redCardLastGame = false;
        gamesNonStopping = 0;
        injuryWeeksLeft = 0;
        sumSkills = forceSkills.reduce((a, b) => a + b, 0);
        for (p = 0; p < 11; p++) {
            pSkills = await assets.encodePlayerSkills(forceSkills, dayOfBirth21, gen = 0, playerId + p, [pot, fwd442[p], left442[p], aggr],
                alignedEndOfLastHalfTwoVec[0], redCardLastGame, gamesNonStopping, 
                injuryWeeksLeft, subLastHalf, sumSkills).should.be.fulfilled 
            teamState.push(pSkills)
        }
        p = 10;
        pSkills = await assets.encodePlayerSkills(forceSkills, dayOfBirth21, gen = 0, playerId + p, [pot, fwd442[p], left442[p], aggr],
                alignedEndOfLastHalfTwoVec[1], redCardLastGame, gamesNonStopping, 
                injuryWeeksLeft, subLastHalf, sumSkills).should.be.fulfilled 
        for (p = 11; p < PLAYERS_PER_TEAM_MAX; p++) {
            teamState.push(pSkills)
        }        
        return teamState;
    };

    const createTeamStateFromSinglePlayer = async (skills, engine, forwardness = 3, leftishness = 2, alignedEndOfLastHalfTwoVec = [false, false]) => {
        teamState = []
        sumSkills = skills.reduce((a, b) => a + b, 0);
        var playerStateTemp = await assets.encodePlayerSkills(
            skills, dayOfBirth21, gen = 0, playerId = 2132321, [potential = 3, forwardness, leftishness, aggr = 0],
            alignedEndOfLastHalfTwoVec[0], redCardLastGame = false, gamesNonStopping = 0, 
            injuryWeeksLeft = 0, subLastHalf, sumSkills
        ).should.be.fulfilled;
        for (player = 0; player < 11; player++) {
            teamState.push(playerStateTemp)
        }

        playerStateTemp = await assets.encodePlayerSkills(
            skills, dayOfBirth21, gen = 0, playerId = 2132321, [potential = 3, forwardness, leftishness, aggr = 0],
            alignedEndOfLastHalfTwoVec[1], redCardLastGame = false, gamesNonStopping = 0, 
            injuryWeeksLeft = 0, subLastHalf, sumSkills
        ).should.be.fulfilled;
        for (player = 11; player < PLAYERS_PER_TEAM_MAX; player++) {
            teamState.push(playerStateTemp)
        }
        return teamState;
    };

    const createHardcodedTeam = function () {
        // returns 18 players generated with the following code. We hardcode it to avoid the "deployDate" time-dependency
        // teamState = [];
        // playerId0 = await assets.encodeTZCountryAndVal(tz = 1, countryIdx = 0, playerIdx = 0).should.be.fulfilled;
        // for (p = 0; p < 18; p++) {
        //     skills = await assets.getPlayerSkillsAtBirth(playerId0.toNumber() + p);
        //     teamState.push(skills);
        //     console.log(skills.toString(10))
        // }
        return [
            '14606248079918261338806855150670198598294524424421999',
            '14603325075249802958062362651785117246719383552393656',
            '14615017086954653606499907426763036762091679724733245',
            '14609171184243174825485386589332947715467405749846827',
            '14615017461189033969342085869889674545308663693968083',
            '14603325891317697566792669908219362044711638355411673',
            '14606249873734453245614329076439313941148075272765994',
            '14603324461979309998470701478621001103697221903123183',
            '14606248281321866413037179508268863783570851530343215',
            '14606249082057998697777445123967984023640370982880706',
            '14603327085801362263089568768708477093108613577769640',
            '14612095382001501327618929648053879079031002742916002',
            '14603326117112742701915784319947485139466656825672861',
            '14612093787498219632679532865607761507997232766977103',
            '14609173081200313275497388848716119026424650418029241',
            '14603326360330245023390630956127251848106222989410926',
            '14606249807529115937477333996086265720951632055960118',
            '14603326808435843856365497638008216685947959514366883'
        ];
    };
    
    beforeEach(async () => {
        training = await TrainingPoints.new().should.be.fulfilled;
        evo = await Evolution.new().should.be.fulfilled;
        play = await PlayAndEvolve.new().should.be.fulfilled;
        engine = await Engine.new().should.be.fulfilled;
        assets = await Assets.new().should.be.fulfilled;
        await assets.init().should.be.fulfilled;
        market = await Market.new().should.be.fulfilled;
        shop = await Shop.new().should.be.fulfilled;
        encodeLog = await EncodingMatchLog.new().should.be.fulfilled;
        precomp = await EnginePreComp.new().should.be.fulfilled;
        applyBoosters = await EngineApplyBoosters.new().should.be.fulfilled;
        await engine.setPreCompAddr(precomp.address).should.be.fulfilled;
        await engine.setApplyBoostersAddr(applyBoosters.address).should.be.fulfilled;
        await training.setAssetsAddress(assets.address).should.be.fulfilled;
        await training.setMarketAddress(market.address).should.be.fulfilled;
        await play.setEngineAddress(engine.address).should.be.fulfilled;
        await play.setTrainingAddress(training.address).should.be.fulfilled;
        await play.setEvolutionAddress(evo.address).should.be.fulfilled;
        await play.setShopAddress(shop.address).should.be.fulfilled;
        
        tactics0 = await engine.encodeTactics(substitutions, subsRounds, setNoSubstInLineUp(lineupConsecutive, substitutions), 
            extraAttackNull, tacticId442).should.be.fulfilled;
        tactics1 = await engine.encodeTactics(substitutions, subsRounds, setNoSubstInLineUp(lineupConsecutive, substitutions), 
            extraAttackNull, tacticId433).should.be.fulfilled;
        tactics1NoChanges = await engine.encodeTactics(noSubstitutions, subsRounds, setNoSubstInLineUp(lineupConsecutive, noSubstitutions), 
            extraAttackNull, tacticId433).should.be.fulfilled;
        tactics442 = await engine.encodeTactics(substitutions, subsRounds, setNoSubstInLineUp(lineupConsecutive, substitutions),
            extraAttackNull, tacticId442).should.be.fulfilled;
        tactics442NoChanges = await engine.encodeTactics(noSubstitutions, subsRounds, setNoSubstInLineUp(lineupConsecutive, noSubstitutions), 
            extraAttackNull, tacticId442).should.be.fulfilled;
        teamStateAll50Half1 = await createTeamStateFromSinglePlayer([50, 50, 50, 50, 50], engine, forwardness = 3, leftishness = 2, aligned = [false, false]).should.be.fulfilled;
        teamStateAll1Half1 = await createTeamStateFromSinglePlayer([1,1,1,1,1], engine, forwardness = 3, leftishness = 2, aligned = [false, false]).should.be.fulfilled;
        teamStateAll50Half2 = await createTeamStateFromSinglePlayer([50, 50, 50, 50, 50], engine, forwardness = 3, leftishness = 2, aligned = [true, false]).should.be.fulfilled;
        teamStateAll1Half2 = await createTeamStateFromSinglePlayer([1,1,1,1,1], engine, forwardness = 3, leftishness = 2, aligned = [true, false]).should.be.fulfilled;
        MAX_RND = await engine.MAX_RND().should.be.fulfilled;
        MAX_RND = MAX_RND.toNumber();
        kMaxRndNumHalf = Math.floor(MAX_RND/2)-200; 
        events1Half = Array.from(new Array(7), (x,i) => 0);
        events1Half = [events1Half,events1Half];
        POINTS_FOR_HAVING_PLAYED = await training.POINTS_FOR_HAVING_PLAYED().should.be.fulfilled;
        POINTS_FOR_HAVING_PLAYED = POINTS_FOR_HAVING_PLAYED.toNumber();
    });

    
    it2('test that used to fail because skills[lineUp[p]] would query skills[25]', async () => {
        seed = '0x6c94aa1a7eea1de18637d1145b6d4bd41cf5f6f8412aae446c2c699d7580ac1f';startTime = '1581951774';matchLog0 = '0';teamId0 = '274877906944';tactic0 = '232408266334649167582215536641';assignedTP0 = '0';players0 = ['14606248079918261338806855269144928920528183545627247','14603325075249802958062362770259847568953042673598904','14615017086954653606499907545237767084325338845938493','14609171184243174825485386707807678037701064871052075','14615017461189033969342085988364404867542322815173331','14603325891317697566792670026694092366945297476616921','14606249873734453245614329194914044263381734393971242','14603324461979309998470701597095731425930881024328431','14606248281321866413037179626743594105804510651548463','14606249082057998697777445242442714345874030104085954','14603327085801362263089568887183207415342272698974888','14612095382001501327618929766528609401264661864121250','14603326117112742701915784438422215461700315946878109','14612093787498219632679532984082491830230891888182351','14609173081200313275497388967190849348658309539234489','14603326360330245023390631074601982170339882110616174','14606249807529115937477334114560996043185291177165366','14603326808435843856365497756482947008181618635572131','0','0','0','0','0','0','0',];matchLog1 = '0';teamId1 = '274877906951';tactic1 = '232408266302079135077072109569';assignedTP1 = '0';players1 = ['14615016376815298690800201649220184280315730971132558','14609172511834412425521368984185260418865566827283036','14609171084586719719561567913262331453334268194587406','14609172165475963560842787370746505659732178042290961','14612094897657191547041386733102280708157489908351780','14609171364042932988648677202799875053042440135311897','14606248714792601209485990362067212005781000358003188','14609173055415076639705784028918284727348393612411594','14609171905532902340470607391083606114650385692034077','14609172622641240130721037564311250677507995239581185','14603325390944727174772193097761782592653101121733224','14603324761645573603736249750401919269415400293270169','14603324774189742777909804362708129945470638967817654','14609171585656588399047378013534405380348672917505319','14609173082594109850415535128877508619287877366448825','14612096081687381931527530703691145228948441982501521','14612093676691391927490720233815463751121253833769674','14606249096692862734323783960084670624419958191030946','0','0','0','0','0','0','0',];
        var {0: skills, 1: matchLogsAndEvents} =  await play.play1stHalfAndEvolve(
            seed, startTime, [players0, players1], [teamId0, teamId1], [tactic0, tactic1], [matchLog0, matchLog1],
            [is2nd = false, isHom = true, isPlay = false], [assignedTP0, assignedTP1]).should.be.fulfilled;
    });
    
        
    it('test that we can have two yellow cards on the same guy that also sees a red card, the latter to be interpreted as what happens when you see a 2 yellows', async () => {
        utils = await Utils.new().should.be.fulfilled;
        seed = '0x8527a891e224136950ff32ca212b45bc93f69fbb801c3b1ebedac52775f99e61';startTime = '1790899200';matchLog0 = '1809252841225230840719990802576568215898489130205763936236440661413666488320';teamId0 = '274877906944';tactic0 = '340596594427581673436941882753025';assignedTP0 = '0';players0 = ['444839120007985571215337246103345753542683081535197729558926920581886','444839120007985571215331537112574929703158848391319931578381389595390','444839120007985571215331537112574929703158848391319931578381389595390','444839120007985571215337246103345753542683081535197729558926920581886','444839120007985571215337246103345753542683081535197729558926920581886','444839120007985571215337246103345753542683081535197729558926920581886','444839120007985571215337246103345753542683081535197729558926920581886','444839120007985571215337246103345753542683081535197729558926920581886','444839120007985571215337246103345753542683081535197729558926920581886','444839120007985571215337246103345753542683081535197729558926920581886','444839120007985571215337246103345753542683081535197729558926920581886','444839120007985571215337246103345753542683081535197729558926920581886','444839120007985571215337246103345753542683081535197729558926920581886','444839120007985571215331537112574929703158848391319931578381389595390','444839120007985571215331537112574929703158848391319931578381389595390','444839120007985571215331537112574929703158848391319931578381389595390','444839120007985571215331537112574929703158848391319931578381389595390','444839120007985571215331537112574929703158848391319931578381389595390','444839120007985571215331537112574929703158848391319931578381389595390','444839120007985571215331537112574929703158848391319931578381389595390','444839120007985571215331537112574929703158848391319931578381389595390','444839120007985571215331537112574929703158848391319931578381389595390','444839120007985571215331537112574929703158848391319931578381389595390','444839120007985571215331537112574929703158848391319931578381389595390','444839120007985571215331537112574929703158848391319931578381389595390',];matchLog1 = '1809252841804448444897055329415117368158677453627034734049075616403095709285';teamId1 = '274877906945';tactic1 = '340596594427581673436941882753025';assignedTP1 = '0';players1 = ['13479973333575334506511561091326173164915009765975984341840712893241','13479973333575334506505852100555349325390776622098186361295181906745','13479973333575334506505852100555349325390776622098186361295181906745','13479973333575334506511561091326173164915009765975984341840712893241','13479973333575334506511561091326173164915009765975984341840712893241','13479973333575334506511561091326173164915009765975984341840712893241','13479973333575334506511561091326173164915009765975984341840712893241','13479973333575334506511561091326173164915009765975984341840712893241','13479973333575334506511561091326173164915009765975984341840712893241','13479973333575334506511561091326173164915009765975984341840712893241','13479973333575334506511561091326173164915009765975984341840712893241','13479973333575334506511561091326173164915009765975984341840712893241','13479973333575334506511561091326173164915009765975984341840712893241','13479973333575334506505852100555349325390776622098186361295181906745','13479973333575334506505852100555349325390776622098186361295181906745','13479973333575334506505852100555349325390776622098186361295181906745','13479973333575334506505852100555349325390776622098186361295181906745','13479973333575334506505852100555349325390776622098186361295181906745','13479973333575334506505852100555349325390776622098186361295181906745','13479973333575334506505852100555349325390776622098186361295181906745','13479973333575334506505852100555349325390776622098186361295181906745','13479973333575334506505852100555349325390776622098186361295181906745','13479973333575334506505852100555349325390776622098186361295181906745','13479973333575334506505852100555349325390776622098186361295181906745','13479973333575334506505852100555349325390776622098186361295181906745',];
        var {0: skills, 1: matchLogsAndEvents} =  await play.play1stHalfAndEvolve(
            seed, startTime, [players0, players1], [teamId0, teamId1], [tactic0, tactic1], [matchLog0, matchLog1],
            [is2nd = false, isHom = true, isPlay = false], [assignedTP0, assignedTP1]).should.be.fulfilled;
        
        decoded = await evo.decodeTactics(tactic0).should.be.fulfilled;
        // show that there are no substitutions:
        let {0: subs, 1: roun, 2: line, 3: attk, 4: tac} = decoded;
        debug.compareArrays(subs, [NO_SUBST, NO_SUBST, NO_SUBST], toNum = true, verbose = false, isBigNumber = true);
        
        // show that the
        var {0: sumSkills , 1: winner, 2: nGoals, 3: TPs, 4: outPlayer, 5: typeOut, 6: outRounds, 7: yellow1, 8: yellow2, 9: subs1, 10: subs2, 11: subs3 } = await utils.fullDecodeMatchLog(matchLogsAndEvents[0], is2nd = false).should.be.fulfilled;
        outPlayer.toNumber().should.be.equal(0);
        typeOut.toNumber().should.be.equal(RED_CARD);
        outRounds.toNumber().should.be.equal(9);
        yellow1.toNumber().should.be.equal(0);
        yellow2.toNumber().should.be.equal(0);
    
        expectedReds = Array.from(new Array(5), (x,i) => false);
        expectedReds[0] = true;
        reds = [];
        for (p=0; p < 25; p++) {       
            red = await assets.getRedCardLastGame(skills[0][p]).should.be.fulfilled;
            reds.push(red);
        }
        debug.compareArrays(reds, expectedReds, toNum = false, verbose = false, isBigNumber = false);
    });

    return
    
    it('show that a red card is stored in skills after playing 1st half', async () => {
        TP = 0;
        assignment = 0
        prev2ndHalfLog = 0;
        teamIds = [1,2]
        vSeed = '0x234a2b366'
        var {0: skills, 1: matchLogsAndEvents} =  await play.play1stHalfAndEvolve(
            vSeed, now, [teamStateAll50Half1, teamStateAll50Half1], teamIds, [tactics0, tactics1], [prev2ndHalfLog, prev2ndHalfLog],
            [is2nd = false, isHomeStadium, isPlayoff], [assignment, assignment]
        ).should.be.fulfilled;
        outType = await training.getOutOfGameType(matchLogsAndEvents[0], is2 = false).should.be.fulfilled;
        outType.toNumber().should.be.equal(3); // RED_CARD = 3
        // with this seed, player p = 8 sees the red card
        outPlayer = await training.getOutOfGamePlayer(matchLogsAndEvents[0], is2 = false).should.be.fulfilled;
        outPlayer.toNumber().should.be.equal(8);
        p = 8;    
        red = await assets.getRedCardLastGame(skills[0][p]).should.be.fulfilled;
        red.should.be.equal(true)
    });
    
    it('updateSkillsAfterPlayHalf: half 1', async () => {
        // note: substitutions = [6, 10, 0];
        // note: lineup is consecutive
        matchLog = await engine.playHalfMatch(
            123456, now, [teamStateAll50Half1, teamStateAll50Half1], [tactics0, tactics1], [0, 0], 
            [is2nd = false, isHome = true, playoff = false]
        ).should.be.fulfilled;
        newSkills = await evo.updateSkillsAfterPlayHalf(teamStateAll50Half1, matchLog[0], tactics0, is2nd = false).should.be.fulfilled;
        // players not aligned did not change state: 
        debug.compareArrays(newSkills.slice(14,25), teamStateAll50Half1.slice(14,25), toNum = false, verbose = false, isBigNumber = true);
        // those that were aligned either finished the 1st half, or were substituted:
        aligned = await evo.setAlignedEndOfFirstHalf(teamStateAll50Half1[0], true).should.be.fulfilled
        substituted = await evo.setSubstitutedFirstHalf(teamStateAll50Half1[0], true).should.be.fulfilled
        for (p = 0; p < 14; p++) {
            if (!substitutions.includes(p)) {newSkills[p].should.be.bignumber.equal(aligned);}
            else {newSkills[p].should.be.bignumber.equal(substituted);}
        }
        
        // now try the same with a red card:
        newLog = await evo.addOutOfGame(matchLog[0], player = 1, round = 2, typeOfOutOfGame = RED_CARD, is2nd = false).should.be.fulfilled;
        newSkills = await evo.updateSkillsAfterPlayHalf(teamStateAll50Half1, newLog, tactics0, is2nd = false).should.be.fulfilled;
        debug.compareArrays(newSkills.slice(14,25), teamStateAll50Half1.slice(14,25), toNum = false, verbose = false, isBigNumber = true);
        alignedRedCarded = await evo.setRedCardLastGame(aligned, true).should.be.fulfilled
        newSkills[1].should.be.bignumber.equal(alignedRedCarded);
        for (p = 0; p < 14; p++) {
            if (p != 1) {
                if (!substitutions.includes(p)) {newSkills[p].should.be.bignumber.equal(aligned);}
                else {newSkills[p].should.be.bignumber.equal(substituted);}
            } 
        }
        
        // now try the same with a hard injury:
        SOFT_INJURY = 1;
        HARD_INJURY = 2;
        WEEKS_SOFT_INJ = 2;
        WEEKS_HARD_INJ = 5;
        newLog = await evo.addOutOfGame(matchLog[0], player = 1, round = 2, typeOfOutOfGame = HARD_INJURY, is2nd = false).should.be.fulfilled;
        newSkills = await evo.updateSkillsAfterPlayHalf(teamStateAll50Half1, newLog, tactics0, is2nd = false).should.be.fulfilled;
        debug.compareArrays(newSkills.slice(14,25), teamStateAll50Half1.slice(14,25), toNum = false, verbose = false, isBigNumber = true);
        alignedInjured = await evo.setInjuryWeeksLeft(aligned, WEEKS_HARD_INJ).should.be.fulfilled
        newSkills[1].should.be.bignumber.equal(alignedInjured);
        for (p = 0; p < 14; p++) {
            if (p != 1) {
                if (!substitutions.includes(p)) {newSkills[p].should.be.bignumber.equal(aligned);}
                else {newSkills[p].should.be.bignumber.equal(substituted);}
            } 
        }
        // now try the same with a soft injury:
        newLog = await evo.addOutOfGame(matchLog[0], player = 1, round = 2, typeOfOutOfGame = SOFT_INJURY, is2nd = false).should.be.fulfilled;
        newSkills = await evo.updateSkillsAfterPlayHalf(teamStateAll50Half1, newLog, tactics0, is2nd = false).should.be.fulfilled;
        debug.compareArrays(newSkills.slice(14,25), teamStateAll50Half1.slice(14,25), toNum = false, verbose = false, isBigNumber = true);
        alignedInjured = await evo.setInjuryWeeksLeft(aligned, WEEKS_SOFT_INJ).should.be.fulfilled
        newSkills[1].should.be.bignumber.equal(alignedInjured);
        for (p = 0; p < 14; p++) {
            if (p != 1) {
                if (!substitutions.includes(p)) {newSkills[p].should.be.bignumber.equal(aligned);}
                else {newSkills[p].should.be.bignumber.equal(substituted);}
            } 
        }
    });
    
    it('updateSkillsAfterPlayHalf: half 2', async () => {
        // note: substitutions = [6, 10, 0];
        // note: lineup is consecutive
        matchLog = await engine.playHalfMatch(
            123456, now, [teamStateAll50Half2, teamStateAll50Half2], [tactics0, tactics1], [0, 0], 
            [is2nd = true, isHome = true, playoff = false]
        ).should.be.fulfilled;
        teamStateAll50Half2[1] = await evo.setInjuryWeeksLeft(teamStateAll50Half2[1], 2);
        newSkills = await evo.updateSkillsAfterPlayHalf(teamStateAll50Half2, matchLog[0], tactics0, is2nd = true).should.be.fulfilled;
        // players not aligned did not change state: 
        debug.compareArrays(newSkills.slice(14,25), teamStateAll50Half2.slice(14,25), toNum = false, verbose = false, isBigNumber = true);
        for (p = 0; p < 25; p++) {
            aligned = await evo.getAlignedEndOfFirstHalf(newSkills[p]).should.be.fulfilled
            aligned.should.be.equal(false)
            substituted = await evo.getSubstitutedFirstHalf(newSkills[p]).should.be.fulfilled
            substituted.should.be.equal(false)
        }
        weeks = await evo.getInjuryWeeksLeft(newSkills[1]);
        weeks.toNumber().should.be.equal(1);
        
        // now try the same with a red card in both halfs...
        newLog = await evo.addOutOfGame(matchLog[0], player = 1, round = 2, typeOfOutOfGame = RED_CARD, is2nd = false).should.be.fulfilled;
        newLog = await evo.addOutOfGame(newLog, player = 2, round = 2, typeOfOutOfGame = RED_CARD, is2nd = true).should.be.fulfilled;
        newSkills = await evo.updateSkillsAfterPlayHalf(teamStateAll50Half2, newLog, tactics0, is2nd = true).should.be.fulfilled;
        debug.compareArrays(newSkills.slice(14,25), teamStateAll50Half2.slice(14,25), toNum = false, verbose = false, isBigNumber = true);
        for (p = 0; p < 25; p++) {
            redCarded = await evo.getRedCardLastGame(newSkills[p]).should.be.fulfilled
            if (p == 1 || p == 2) {redCarded.should.be.equal(true);}
            else {redCarded.should.be.equal(false);}
        }
    });
    
    it('applyTrainingPoints: if assignment = 0, it works by doing absolutely nothing', async () => {
        matchStartTime = now;
        newSkills = await training.applyTrainingPoints(teamStateAll50Half2, assignment = 0, tactics = 0, matchStartTime, TPs = 0).should.be.fulfilled;
        newSkills2 = await training.applyTrainingPoints(teamStateAll50Half2, assignment = 0, tactics = 0, matchStartTime, TPs = 1).should.be.fulfilled;
        debug.compareArrays(newSkills, teamStateAll50Half2, toNum = false, verbose = false, isBigNumber = true);
        debug.compareArrays(newSkills2, teamStateAll50Half2, toNum = false, verbose = false, isBigNumber = true);
    });

    it('training leading to an actual son', async () => {
        playerSkills = await assets.encodePlayerSkills(
            skills = [100, 100, 100, 100, 100], 
            dayOfBirth = 30*365, // 30 years after unix time 
            gen = 45,
            playerId = 2132321,
            [potential = 2, forwardness, leftishness, aggr = 0],
            alignedEndOfLastHalf = true,
            redCardLastGame = false,
            gamesNonStopping = 0,
            injuryWeeksLeft = 0,
            subLastHalf,
            sumSkills = 5
        ).should.be.fulfilled;
        age = 40;
        matchStartTime = dayOfBirth*24*3600 + Math.floor(age*365*24*3600/7);
        
        TPs = 20;
        TPperSkill = Array.from(new Array(5), (x,i) => TPs/5 - 3*i % 5);
        newSkills = await training.evolvePlayer(playerSkills, TPperSkill, matchStartTime).should.be.fulfilled;

        // checks that the generation increases by 1. 
        // It sets a "32" at the beginning if it is a Academy player, otherwise it is a child
        // In this case, the randomness leads to an Academy player
        result = await assets.getGeneration(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(gen + 1)

        playerSkills = await assets.encodePlayerSkills(
            skills = [100, 100, 100, 100, 100], 
            dayOfBirth = 30*365, // 30 years after unix time 
            gen = 45,
            playerId = 3,
            [potential = 2, forwardness, leftishness, aggr = 0],
            alignedEndOfLastHalf = true,
            redCardLastGame = false,
            gamesNonStopping = 0,
            injuryWeeksLeft = 0,
            subLastHalf,
            sumSkills = 5
        ).should.be.fulfilled;
        age = 40;
        matchStartTime = dayOfBirth*24*3600 + Math.floor(age*365*24*3600/7);
        
        TPs = 20;
        TPperSkill = Array.from(new Array(5), (x,i) => TPs/5 - 3*i % 5);
        newSkills = await training.evolvePlayer(playerSkills, TPperSkill, matchStartTime).should.be.fulfilled;

        // checks that the generation increases by 1. 
        // It sets a "32" at the beginning if it is a Academy player, otherwise it is a child
        // In this case, the randomness leads to a son
        result = await assets.getGeneration(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(gen - 32 + 1)
        
        expected = [ 531, 1506, 912, 551, 1500 ];
        N_SKILLS = 5;
        results = [];
        for (sk = 0; sk < N_SKILLS; sk++) {
            result = await engine.getSkill(newSkills, sk).should.be.fulfilled;
            results.push(result);
        }
        debug.compareArrays(results, expected, toNum = true, verbose = false);
        
        expectedSumSkills = expected.reduce((a, b) => a + b, 0);
        result = await engine.getSumOfSkills(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expectedSumSkills);
        
    });
    
    it('training leading to an academy', async () => {
        // all inputs are identical to the previous test, except for a +2 in matchStatTime,
        // which changes the entire randomness
        playerSkills = await assets.encodePlayerSkills(
            skills = [100, 100, 100, 100, 100], 
            dayOfBirth = 30*365, // 30 years after unix time 
            gen = 3,
            playerId = 2132321,
            [potential = 2, forwardness, leftishness, aggr = 0],
            alignedEndOfLastHalf = true,
            redCardLastGame = false,
            gamesNonStopping = 0,
            injuryWeeksLeft = 0,
            subLastHalf,
            sumSkills = 5
        ).should.be.fulfilled;
        age = 40;
        matchStartTime = dayOfBirth*24*3600 + Math.floor(age*365*24*3600/7);
        
        TPs = 20;
        TPperSkill = Array.from(new Array(5), (x,i) => TPs/5 - 3*i % 5);
        newSkills = await training.evolvePlayer(playerSkills, TPperSkill, matchStartTime + 2).should.be.fulfilled;

        // checks that the generation increases by 1. It sets a "32" at the beginning if it is a Academy player, otherwise it is a child.
        // In this case, randomness leads to an academy.
        result = await assets.getGeneration(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(32 + gen + 1)

        expected = [ 755, 920, 1455, 762, 1107 ];
        results = []
        for (sk = 0; sk < N_SKILLS; sk++) {
            result = await engine.getSkill(newSkills, sk).should.be.fulfilled;
            results.push(result);
        }
        debug.compareArrays(results, expected, toNum = true, verbose = false);
        
        expectedSumSkills = expected.reduce((a, b) => a + b, 0);
        result = await engine.getSumOfSkills(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expectedSumSkills);
        
    });
    
    
    it('applyTrainingPoints', async () => {
        TP = 200;
        TPperSkill = Array.from(new Array(25), (x,i) => TP/5 - 3*i % 6);
        specialPlayer = 21;
        // make sure they sum to TP:
        for (bucket = 0; bucket < 5; bucket++){
            sum4 = 0;
            for (sk = 5 * bucket; sk < (5 * bucket + 4); sk++) {
                sum4 += TPperSkill[sk];
            }
            TPperSkill[5 * bucket + 4] = TP - sum4;
        }        
        assignment = await training.encodeTP(TP, TPperSkill, specialPlayer).should.be.fulfilled;
        matchStartTime = now;
        newSkills = await training.applyTrainingPoints(teamStateAll50Half2, assignment, tactics = 0, matchStartTime, TP+1).should.be.rejected;
        newSkills = await training.applyTrainingPoints(teamStateAll50Half2, assignment, tactics = 0, matchStartTime, TP).should.be.fulfilled;
        for (p = 0; p < 25; p++) {
            result = await training.getSkill(newSkills[p], SK_SHO).should.be.fulfilled;
            if (p == specialPlayer) result.toNumber().should.be.equal(110);
            else result.toNumber().should.be.equal(105);
        }
    });
    
    it('applyTrainingPoints with recovery stamina', async () => {
        const [TP, TPperSkill] = getDefaultTPs();
        assignment = await training.encodeTP(TP, TPperSkill, specialPlayer).should.be.fulfilled;
        matchStartTime = now;
        staminas = Array.from(new Array(PLAYERS_PER_TEAM_MAX), (x,i) => i % 4); 
        gamesNonStopping = Array.from(new Array(PLAYERS_PER_TEAM_MAX), (x,i) => i % 7); 
        skills = [...teamStateAll50Half2];
        for (p = 0; p < PLAYERS_PER_TEAM_MAX; p++){
            skills[p] = await evo.setGamesNonStopping(skills[p], gamesNonStopping[p]).should.be.fulfilled;
        }
        tactics = await training.setStaminaRecovery(initTactics = 0, staminas);
        newSkills = await training.applyTrainingPoints(skills, assignment, tactics, matchStartTime, TP+1).should.be.rejected;
        newSkills = await training.applyTrainingPoints(skills, assignment, tactics, matchStartTime, TP).should.be.fulfilled;
        newGamesNonStopping = [];
        expectedGamesNonStopping = [];
        for (p = 0; p < 25; p++) {
            result = await training.getSkill(newSkills[p], SK_SHO).should.be.fulfilled;
            if (p == specialPlayer) result.toNumber().should.be.equal(110);
            else result.toNumber().should.be.equal(105);
            result = await evo.getGamesNonStopping(newSkills[p]).should.be.fulfilled;
            newGamesNonStopping.push(result);
            expected = 0;
            if (staminas[p] == 0) { expected = gamesNonStopping[p] }
            else if (staminas[p] == 3 || gamesNonStopping[p] <= 2*staminas[p] ) { expected = 0 }
            else { expected = gamesNonStopping[p] - 2 * staminas[p]}
            expectedGamesNonStopping.push(expected)
        }
        debug.compareArrays(newGamesNonStopping, expectedGamesNonStopping, toNum = true, verbose = false);
    });
    
    it('applyTrainingPoints with realistic team and zero TPs', async () => {
        teamState = createHardcodedTeam();
        for (p = 18; p < 25; p++) teamState.push(0);
        TPperSkill = Array.from(new Array(25), (x,i) => 0);
        TP = TPperSkill.reduce((a, b) => a + b, 0);
        assignment = await training.encodeTP(TP, TPperSkill, specialPlayer = 0).should.be.fulfilled;
        matchStartTime = now;
        newSkills = await training.applyTrainingPoints(teamState, assignment, tactics = 0, matchStartTime, TP);
        initShoot = [];
        newShoot = [];
        expectedInitShoot = [ 623, 440, 829, 811, 723, 729, 554, 751, 815, 1474, 680, 930, 1181, 1103, 697, 622, 566, 931 ];
        expectedNewShoot  = [ 623, 440, 829, 811, 723, 702, 554, 735, 815, 1466, 680, 930, 1181, 1095, 697, 622, 566, 931 ];
        // check that if skills are different, then:
        // - the new ones are worse than the init ones,
        // - it happened because of age (older than 31 y.o.)
        for (p = 0; p < 18; p++) {
            resultInit = await training.getSkill(teamState[p], SK_SHO).should.be.fulfilled;
            resultNew = await training.getSkill(newSkills[p], SK_SHO).should.be.fulfilled;
            if (resultNew.toNumber() != resultInit.toNumber()) {
                resultId = await assets.getPlayerIdFromSkills(newSkills[p]).should.be.fulfilled;
                resultAge = await assets.getPlayerAgeInDays(resultId).should.be.fulfilled;
                (resultAge.toNumber() >= 31 * 365).should.be.equal(true);
                (resultNew.toNumber() < resultInit.toNumber()).should.be.equal(true);
            }
            initShoot.push(resultInit)
            newShoot.push(resultNew)
        }
        debug.compareArrays(newShoot, expectedNewShoot, toNum = true, verbose = false);
        debug.compareArrays(initShoot, expectedInitShoot, toNum = true, verbose = false);
    });
    
    it('applyTrainingPoints with realistic team and non-zero TPs', async () => {
        teamState = createHardcodedTeam();
        for (p = 18; p < 25; p++) teamState.push(0);
        TPperSkill = [ 40, 37, 40, 37, 46, 37, 40, 37, 40, 46, 40, 37, 40, 37, 46, 37, 40, 37, 40, 46, 40, 37, 40, 37, 46 ];
        TP = 200;
        assignment = await training.encodeTP(TP, TPperSkill, specialPlayer = 12).should.be.fulfilled;
        matchStartTime = now;
        newSkills = await training.applyTrainingPoints(teamState, assignment, tactics = 0, matchStartTime, TP);
        initShoot = [];
        newShoot = [];
        expectedNewShoot  = [ 673, 480, 869, 987, 1015, 739, 591, 772, 1009, 1506, 906, 1178, 1452, 1147, 905, 816, 603, 1120 ];
        expectedInitShoot = [ 623, 440, 829, 811, 723, 729, 554,  751,  815, 1474, 680,  930, 1181, 1103, 697, 622, 566,  931 ];
        for (p = 0; p < 18; p++) {
            result0 = await training.getSkill(teamState[p], SK_SHO);
            result1 = await training.getSkill(newSkills[p], SK_SHO);
            initShoot.push(result0)
            newShoot.push(result1)
        }
        debug.compareArrays(newShoot, expectedNewShoot, toNum = true, verbose = false);
        debug.compareArrays(initShoot, expectedInitShoot, toNum = true, verbose = false);
    });

    it('test evolvePlayer at zero potential', async () => {
        playerSkills = await assets.encodePlayerSkills(
            skills = [100, 100, 100, 100, 100], 
            dayOfBirth = 30*365, // 30 years after unix time 
            gen = 0,
            playerId = 2132321,
            [potential = 0, forwardness, leftishness, aggr = 0],
            alignedEndOfLastHalf = true,
            redCardLastGame = false,
            gamesNonStopping = 0,
            injuryWeeksLeft = 0,
            subLastHalf,
            sumSkills = 5
        ).should.be.fulfilled;
        age = 16;
        matchStartTime = dayOfBirth*24*3600 + Math.floor(age*365*24*3600/7);
        
        TPs = 20;
        TPperSkill = [10, 20, 30, 40, 50];
        newSkills = await training.evolvePlayer(playerSkills, TPperSkill, matchStartTime).should.be.fulfilled;
        expected = [110,120,130,140,150]; // at zero potential, it's easy
        results = []
        for (sk = 0; sk < N_SKILLS; sk++) {
            result = await engine.getSkill(newSkills, sk).should.be.fulfilled;
            results.push(result);
        }
        debug.compareArrays(results, expected, toNum = true, verbose = false);
    });
    
    it('test evolvePlayer with TPs= 0', async () => {
        playerSkills = await assets.encodePlayerSkills(
            skills = [12, 13, 155, 242, 32], 
            dayOfBirth = 30*365, // 30 years after unix time 
            gen = 0,
            playerId = 2132321,
            [potential = 6, forwardness, leftishness, aggr = 0],
            alignedEndOfLastHalf = true,
            redCardLastGame = false,
            gamesNonStopping = 0,
            injuryWeeksLeft = 0,
            subLastHalf,
            sumSkills = 5
        ).should.be.fulfilled;
        age = 16;
        matchStartTime = dayOfBirth*24*3600 + Math.floor(age*365*24*3600/7);
        
        TPperSkill = [0, 0, 0, 0, 00];
        newSkills = await training.evolvePlayer(playerSkills, TPperSkill, matchStartTime).should.be.fulfilled;
        result = await engine.getSkill(newSkills, SK_SHO).should.be.fulfilled;
        expected = skills;
        results = []
        for (sk = 0; sk < N_SKILLS; sk++) {
            result = await engine.getSkill(newSkills, sk).should.be.fulfilled;
            results.push(result);
        }
        debug.compareArrays(results, expected, toNum = true, verbose = false);
    });
    
    
    it('test evolvePlayer at non-zero potential', async () => {
        playerSkills = await assets.encodePlayerSkills(
            skills = [100, 100, 100, 100, 100], 
            dayOfBirth = 30*365, // 30 years after unix time 
            gen = 0,
            playerId = 2132321,
            [potential = 1, forwardness, leftishness, aggr = 0],
            alignedEndOfLastHalf = true,
            redCardLastGame = false,
            gamesNonStopping = 0,
            injuryWeeksLeft = 0,
            subLastHalf,
            sumSkills = 5
        ).should.be.fulfilled;
        age = 16;
        matchStartTime = dayOfBirth*24*3600 + Math.floor(age*365*24*3600/7);
        
        TPperSkill = [10, 20, 30, 40, 50];
        newSkills = await training.evolvePlayer(playerSkills, TPperSkill, matchStartTime).should.be.fulfilled;
        result = await engine.getSkill(newSkills, SK_SHO).should.be.fulfilled;
        expected = [ 113, 126, 140, 153, 166 ];
        results = []
        for (sk = 0; sk < N_SKILLS; sk++) {
            result = await engine.getSkill(newSkills, sk).should.be.fulfilled;
            results.push(result);
        }
        debug.compareArrays(results, expected, toNum = true, verbose = false);


        expectedSumSkills = expected.reduce((a, b) => a + b, 0);
        result = await engine.getSumOfSkills(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expectedSumSkills);
    });

    it('test evolvePlayer at non-zero potential and age', async () => {
        playerSkills = await assets.encodePlayerSkills(
            skills = [100, 100, 100, 100, 100], 
            dayOfBirth = 30*365, // 30 years after unix time 
            gen = 0,
            playerId = 2132321,
            [potential = 2, forwardness, leftishness, aggr = 0],
            alignedEndOfLastHalf = true,
            redCardLastGame = false,
            gamesNonStopping = 0,
            injuryWeeksLeft = 0,
            subLastHalf,
            sumSkills = 5
        ).should.be.fulfilled;
        age = 17;
        matchStartTime = dayOfBirth*24*3600 + Math.floor(age*365*24*3600/7);
        
        TPperSkill = [10, 20, 30, 40, 50];
        newSkills = await training.evolvePlayer(playerSkills, TPperSkill, matchStartTime).should.be.fulfilled;
        result = await engine.getSkill(newSkills, SK_SHO).should.be.fulfilled;
        expected = [121, 143, 165, 186, 208];
        results = []
        for (sk = 0; sk < N_SKILLS; sk++) {
            result = await engine.getSkill(newSkills, sk).should.be.fulfilled;
            results.push(result);
        }
        debug.compareArrays(results, expected, toNum = true, verbose = false);

        
        expectedSumSkills = expected.reduce((a, b) => a + b, 0);
        result = await engine.getSumOfSkills(newSkills).should.be.fulfilled;
        result.toNumber().should.be.equal(expectedSumSkills);
    });

    it('test evolvePlayer with old age', async () => {
        playerSkills = await assets.encodePlayerSkills(
            skills = [1000, 2000, 3000, 4000, 5000], 
            dayOfBirth = 30*365, // 30 years after unix time 
            gen = 0,
            playerId = 2132321,
            [potential = 2, forwardness, leftishness, aggr = 0],
            alignedEndOfLastHalf = true,
            redCardLastGame = false,
            gamesNonStopping = 0,
            injuryWeeksLeft = 0,
            subLastHalf,
            sumSkills = 5
        ).should.be.fulfilled;
        age = 35;
        matchStartTime = dayOfBirth*24*3600 + Math.floor(age*365*24*3600/7);
        
        TPperSkill = [0, 0, 0, 0, 0];
        newSkills = await training.evolvePlayer(playerSkills, TPperSkill, matchStartTime).should.be.fulfilled;
        expected = [968, 1968, 2968, 3968, 4968]; // -32 per game
        results = []
        for (sk = 0; sk < N_SKILLS; sk++) {
            result = await engine.getSkill(newSkills, sk).should.be.fulfilled;
            results.push(result);
        }
        debug.compareArrays(results, expected, toNum = true, verbose = false);

        
    });

    it('test that we can play a 1st half with log = assignedTPs = 0', async () => {
        TP = 0;
        assignment = 0
        prev2ndHalfLog = 0;
        teamIds = [1,2]
        verseSeed = '0x234ab3'
        await play.play1stHalfAndEvolve(
            verseSeed, now, [teamStateAll50Half1, teamStateAll50Half1], teamIds, [tactics0, tactics1], [prev2ndHalfLog, prev2ndHalfLog],
            [is2nd = false, isHomeStadium, isPlayoff], [assignment, assignment]
        ).should.be.fulfilled;

        prev2ndHalfLog = await evo.addTrainingPoints(0, TP = 2).should.be.fulfilled;
        await play.play1stHalfAndEvolve(
            verseSeed, now, [teamStateAll50Half1, teamStateAll50Half1], teamIds, [tactics0, tactics1], [prev2ndHalfLog, prev2ndHalfLog],
            [is2nd = false, isHomeStadium, isPlayoff], [assignment, assignment]
        ).should.be.fulfilled;
    });
    
    it('test that we can a 1st half and include apply training points too', async () => {
        const [TP, TPperSkill] = getDefaultTPs();
        assignment = await training.encodeTP(TP, TPperSkill, specialPlayer).should.be.fulfilled;
        // Should be rejected if we earned 0 TPs in previous match, and now we claim 200 in the assignedTPs:
        prev2ndHalfLog = 0;
        teamIds = [1,2]
        verseSeed = '0x234ab3'
        
        lineUpNew = [...lineupConsecutive];
        lineUpNew[0] = 16;
        subst = [6, 10, 0]
        tacticsNew = await engine.encodeTactics(subst, subsRounds, setNoSubstInLineUp(lineUpNew, subst), 
        extraAttackNull, tacticId433).should.be.fulfilled;
        
        await play.play1stHalfAndEvolve(
            verseSeed, now, [teamStateAll50Half1, teamStateAll50Half1], teamIds, [tacticsNew, tacticsNew], [prev2ndHalfLog, prev2ndHalfLog],
            [is2nd = false, isHomeStadium, isPlayoff], [assignment, assignment]
        ).should.be.rejected;
        
        prev2ndHalfLog = await evo.addTrainingPoints(0, TP).should.be.fulfilled;
        const {0: skills, 1: matchLogsAndEvents} = await play.play1stHalfAndEvolve(
            verseSeed, now, [teamStateAll50Half1, teamStateAll50Half1], teamIds, [tacticsNew, tacticsNew], [prev2ndHalfLog, prev2ndHalfLog],
            [is2nd = false, isHomeStadium, isPlayoff], [assignment, assignment]
        ).should.be.fulfilled;

        // matchLogsAndEvents[0].should.be.bignumber.equal('1809251596697222440607644166008099735659887273687611206471872191446072164449');
        // matchLogsAndEvents[1].should.be.bignumber.equal('1809251596697222440607644166008099735659887286364117208754702134768761309058');

        // show that after applying the training points (before the match), the teams evolved from 250 per player to 549
        sumBeforeEvolving = await evo.getSumOfSkills(teamStateAll50Half1[0]).should.be.fulfilled;
        sumBeforeEvolving.toNumber().should.be.equal(250);
        expectedSums = Array.from(new Array(25), (x,i) => 549);
        sumSkills0 = []  // sum of skills of each player for team 0
        sumSkills1 = []  // sum of skills of each player for team 1
        for (p = 0; p < 25; p++) {
            sum = await evo.getSumOfSkills(skills[0][p]).should.be.fulfilled;
            sumSkills0.push(sum)
            sum = await evo.getSumOfSkills(skills[1][p]).should.be.fulfilled;
            sumSkills1.push(sum)
        }
        debug.compareArrays(sumSkills0, expectedSums, toNum = true, verbose = false, isBigNumber = false);
        debug.compareArrays(sumSkills1, expectedSums, toNum = true, verbose = false, isBigNumber = false);

        // check that the game is played, ends up in 2-2, and that there are no TPs assigned (this is 1st half)
        expectedGoals = [1, 2];
        expectedPoints = [0, 0];
        goals = []
        points = []
        for (team = 0; team < 2; team++) {
            nGoals = await encodeLog.getNGoals(matchLogsAndEvents[team]);
            goals.push(nGoals);
            nPoints = await encodeLog.getTrainingPoints(matchLogsAndEvents[team]).should.be.fulfilled;
            points.push(nPoints);
        }
        debug.compareArrays(goals, expectedGoals, toNum = true, verbose = false, isBigNumber = false);
        debug.compareArrays(points, expectedPoints, toNum = true, verbose = false, isBigNumber = false);
        // check that the events are generated, and match whatever we got once.
        expected = [ 1, 1, 8, 1, 8, 1, 1, 7, 1, 7, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 10, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 6, 0, 1, 9, 0, 0 ];
        debug.compareArrays(matchLogsAndEvents.slice(2), expected, toNum = true, verbose = false, isBigNumber = false);

        // check that all 3 substitutions took place
        for (pos = 0; pos < 3; pos++) {
            result = await evo.getInGameSubsHappened(matchLogsAndEvents[0], pos, is2nd = false);
            result.toNumber().should.be.equal(1);
        }
        
        // check that we set the "aligned" properties properly
        // there where 3 changes in total, so was in LineUp includes the three changes
        // recall:   lineUpNew[0] = 16;  subst = [6, 10, 0]
        // So, using lineUp idx:    6 -> 11, 10 -> 12, 0 -> 13
        // Using shirtNum:          6 -> 11, 10 -> 12, 16 -> 13
        shirtNumSubst = Array.from(subst, (subst,i) => lineUpNew[subst]); 
        for (team = 0; team < 2; team++) {
            for (p = 0; p < 25; p++) {
                endedHalf = await evo.getAlignedEndOfFirstHalf(skills[team][p]).should.be.fulfilled;
                wasSubst = await evo.getSubstitutedFirstHalf(skills[team][p]).should.be.fulfilled;
                wasInLineUp = lineUpNew.includes(p);
                wasSubst = shirtNumSubst.includes(p);
                if (wasInLineUp && !wasSubst) {
                    endedHalf.should.be.equal(true);
                    wasSubst.should.be.equal(false);
                }
                if (wasInLineUp && wasSubst) {
                    endedHalf.should.be.equal(false);
                    wasSubst.should.be.equal(true);
                }
                if (!wasInLineUp) {
                    endedHalf.should.be.equal(false);
                    wasSubst.should.be.equal(false);
                }
            }
        }
    });
    
    it('test that we can play a first half with totally null players, and that they do not evolve', async () => {
        teamIds = [0, 0]
        verseSeed = '0x234ab3'
        emptyTeam = Array.from(new Array(25), (x,i) => 0); 
        assignment = 0;
        const {0: skills, 1: matchLogsAndEvents} = await play.play1stHalfAndEvolve(
            verseSeed, now, [emptyTeam, emptyTeam], teamIds, [tactics0, tactics1], [0, 0], 
            [is2nd = false, isHomeStadium, isPlayoff], [assignment, assignment]
        ).should.be.fulfilled;
        
        expectedGoals = [0, 0];
        expectedPoints = [0, 0];
        goals = []
        points = []
        for (team = 0; team < 2; team++) {
            nGoals = await encodeLog.getNGoals(matchLogsAndEvents[team]);
            goals.push(nGoals);
            nPoints = await encodeLog.getTrainingPoints(matchLogsAndEvents[team]).should.be.fulfilled;
            points.push(nPoints);
        }
        debug.compareArrays(goals, expectedGoals, toNum = true, verbose = false, isBigNumber = false);
        debug.compareArrays(points, expectedPoints, toNum = true, verbose = false, isBigNumber = false);
        debug.compareArrays(skills[0], emptyTeam, toNum = true, verbose = false, isBigNumber = false);
        debug.compareArrays(skills[1], emptyTeam, toNum = true, verbose = false, isBigNumber = false);
    });
        
    
    it('test that we can play a 2nd half, include the training points, and check gamesNonStopping', async () => {
        const [TP, TPperSkill] = getDefaultTPs();
        assignment = await training.encodeTP(TP, TPperSkill, specialPlayer).should.be.fulfilled;
        teamIds = [1,2]
        verseSeed = '0x234ab3'
        prev2ndHalfLog = await evo.addTrainingPoints(0, TP).should.be.fulfilled;

        // FIRST half:
        // add one player who will go from 6 to 7, one that will remain at 7, and two that will reset, since they were not linedUp
        teamStateAll50Half1[0] = await evo.setGamesNonStopping(teamStateAll50Half1[0], 6).should.be.fulfilled;
        teamStateAll50Half1[1] = await evo.setGamesNonStopping(teamStateAll50Half1[1], 7).should.be.fulfilled;
        teamStateAll50Half1[22] = await evo.setGamesNonStopping(teamStateAll50Half1[1], 7).should.be.fulfilled;
        teamStateAll50Half1[23] = await evo.setGamesNonStopping(teamStateAll50Half1[1], 1).should.be.fulfilled;

        // for team1, besides the previous, plan only the 1st of the substitutions
        subst = [...substitutions]; // = [6, 10, 0]
        subst[1] = NO_SUBST;
        subst[2] = NO_SUBST;
        tacticsNew = await engine.encodeTactics(subst, subsRounds, setNoSubstInLineUp(lineupConsecutive, subst), 
        extraAttackNull, tacticId433).should.be.fulfilled;

        // play the 1st half:
        const {0: skills0, 1: matchLogsAndEvents0} = await play.play1stHalfAndEvolve(
            verseSeed, now, [teamStateAll50Half1, teamStateAll50Half1], teamIds, [tactics0, tacticsNew], [prev2ndHalfLog, prev2ndHalfLog],
            [is2nd = false, isHomeStadium, isPlayoff], [assignment, assignment]
        ).should.be.fulfilled;

        // first: check correct properties for team1:
            // recall:   lineUp = consecutive,  subst = [6, NO_SUBST, NO_SUBST]
            // So, using lineUp idx, the sust was::     6 -> 11, 
            // Same as using shirtNum:                  6 -> 11,
        for (p=0; p<25; p++){ 
            result = await engine.getAlignedEndOfFirstHalf(skills0[1][p]).should.be.fulfilled;
            if ((p < 12) && (p!= 6)) result.should.be.equal(true);
            else result.should.be.equal(false);
        }

        // do 1 change at half time for team1, that still had 2 remaining changes.
        lineUpNew = [...lineupConsecutive];
        lineUpNew[3] = 16;
        tactics1NoChangesNew = await engine.encodeTactics(noSubstitutions, subsRounds, setNoSubstInLineUp(lineUpNew, noSubstitutions), 
            extraAttackNull, tacticId433).should.be.fulfilled;
            
        // play half 2:
        const {0: skills, 1: matchLogsAndEvents} = await play.play2ndHalfAndEvolve(
            verseSeed, now, skills0, teamIds, [tactics1NoChanges, tactics1NoChangesNew], matchLogsAndEvents0.slice(0,2), 
            [is2nd = true, isHomeStadium, isPlayoff]
        ).should.be.fulfilled;

        // check that we find the correct halfTimeSubs in the matchLog.
        // note that what is stored is: lineUp[p] + 1 = 17
        expectedHalfTimeSubs = [17,0,0];
        halfTimeSubs = []
        for (pos = 0; pos < 3; pos ++) {
            result = await evo.getHalfTimeSubs(matchLogsAndEvents[1], pos).should.be.fulfilled;
            halfTimeSubs.push(result);
        }
        debug.compareArrays(halfTimeSubs, expectedHalfTimeSubs, toNum = true, verbose = false, isBigNumber = false);

        // check Training Points (and Goals)
        expectedGoals = [3, 5];
        expectedPoints = [23, 49];
        goals = []
        points = []
        for (team = 0; team < 2; team++) {
            nGoals = await encodeLog.getNGoals(matchLogsAndEvents[team]);
            goals.push(nGoals);
            nPoints = await encodeLog.getTrainingPoints(matchLogsAndEvents[team]).should.be.fulfilled;
            points.push(nPoints);
        }
        debug.compareArrays(goals, expectedGoals, toNum = true, verbose = false, isBigNumber = false);
        debug.compareArrays(points, expectedPoints, toNum = true, verbose = false, isBigNumber = false);

        // test that the states did not change the intrinsics of the players:
        sumBeforeEvolving = await evo.getSumOfSkills(skills0[0][2]).should.be.fulfilled;
        sumBeforeEvolving.toNumber().should.be.equal(549);
        expectedSums = Array.from(new Array(25), (x,i) => 549);
        sumSkills0 = []  // sum of skills of each player for team 0
        sumSkills1 = []  // sum of skills of each player for team 1
        for (p = 0; p < 25; p++) {
            sum = await evo.getSumOfSkills(skills[0][p]).should.be.fulfilled;
            sumSkills0.push(sum)
            sum = await evo.getSumOfSkills(skills[1][p]).should.be.fulfilled;
            sumSkills1.push(sum)
        }
        debug.compareArrays(sumSkills0, expectedSums, toNum = true, verbose = false, isBigNumber = false);
        debug.compareArrays(sumSkills1, expectedSums, toNum = true, verbose = false, isBigNumber = false);

        // check that we correctly reset the "played game" and gamesNonStopping properties
        // team0 went through subst: [6, 10, 0], so 6 -> 11, 10 -> 12, 0 -> 13
        // team1 went through subst: [6], so 6 -> 11
        // so we expect that team0 has [0,..13] increasing gamesNonStopping
        // so we expect that team1 has [0,..11] increasing gamesNonStopping
        expectedGamesNonStopping = Array.from(new Array(25), (x,i) => 0);
        for (p=0; p < 14; p++) expectedGamesNonStopping[p] = 1;
        expectedGamesNonStopping[0] = 7;    // 6 -> 7
        expectedGamesNonStopping[1] = 7;    // 7 -> 7
        expectedGamesNonStopping[22] = 0;   // 6 -> 0
        expectedGamesNonStopping[23] = 0;   // 1 -> 0
        expected = [];
        expected.push([...expectedGamesNonStopping]);
        expected[0][23] = 0;
        // team1 particular cases:
        expectedGamesNonStopping[12] = 0;   // subst was not planned for team1
        expectedGamesNonStopping[13] = 0;   // subst was not planned for team1
        expectedGamesNonStopping[16] = 1;   // he joined at half time for team1
        expected.push([...expectedGamesNonStopping]);
        
        for (team = 0; team < 2; team++) {
            nonStoppingGames = [];
            for (p = 0; p < 25; p++) {
                endedHalf = await evo.getAlignedEndOfFirstHalf(skills[team][p]).should.be.fulfilled;
                wasSubst = await evo.getSubstitutedFirstHalf(skills[team][p]).should.be.fulfilled;
                nGamesNonStopping = await evo.getGamesNonStopping(skills[team][p]).should.be.fulfilled;
                endedHalf.should.be.equal(false);
                wasSubst.should.be.equal(false);
                nonStoppingGames.push(nGamesNonStopping);
            }
            debug.compareArrays(nonStoppingGames, expected[team], toNum = true, verbose = false, isBigNumber = false);
        }
    });

    it('training points with random inputs', async () => {
        typeOut = [3, 0];
        outRounds = [7, 0];
        outGames = [9, 14]
        yellows1 = [14, 0]
        yellows2 = [0, 0]
        defs1 = 4; 
        defs2 = 0; 
        numTot = 10; 
        win = 0; 
        isHome = true;
        
        log0 = await logUtils.encodeLog(encodeLog, nGoals = 3, assistersIdx, shootersIdx, shooterForwardPos, penalties,
            outGames, outRounds, typeOut, yellowCardedDidNotFinish1stHalf,
            isHome, ingameSubs1, ingameSubs2, yellows1, yellows2, 
            halfTimeSubstitutions, defs1, defs2, numTot, win, teamSumSkillsDefault, trainingPointsInit);
        
        logFinal = await training.computeTrainingPoints(log0, log0)
        expected = [36, 25];
        for (team = 0; team < 2; team++) {
            points = await encodeLog.getTrainingPoints(logFinal[team]).should.be.fulfilled;
            points.toNumber().should.be.equal(expected[team]);
        }
    });

    it('training points with no goals nor anything else', async () => {
        log0 = await logUtils.encodeLog(encodeLog, nGoals = 0, assistersIdx, shootersIdx, shooterForwardPos, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHomeSt, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkillsDefault, trainingPointsInit);
        
        logFinal = await training.computeTrainingPoints(log0, log0)
        // expect: POINTS_FOR_HAVING_PLAYED(10) + cleanSheet(24+8) = 42
        expected = [42, 42];
        for (team = 0; team < 2; team++) {
            points = await encodeLog.getTrainingPoints(logFinal[team]).should.be.fulfilled;
            points.toNumber().should.be.equal(expected[team]);
        }
    });    

    it('training points with many goals by attackers', async () => {
        goals = 5;
        ass     = Array.from(new Array(goals), (x,i) => 10);
        shoot   = Array.from(new Array(goals), (x,i) => 10);
        fwd     = Array.from(new Array(goals), (x,i) => 3);
        
        log0 = await logUtils.encodeLog(encodeLog, goals, ass, shoot, fwd, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHomeSt, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkillsDefault, trainingPointsInit);
        
        logFinal = await training.computeTrainingPoints(log0, log0)
        // expect: POINTS_FOR_HAVING_PLAYED(10) + GOALS_BY_ATTACKERS(4 * 5) - GOALS_OPPONENT(5)  
        expected = [25, 25];
        for (team = 0; team < 2; team++) {
            points = await encodeLog.getTrainingPoints(logFinal[team]).should.be.fulfilled;
            points.toNumber().should.be.equal(expected[team]);
        }
    });    

    it('training points with many goals by mids', async () => {
        goals = 5;
        ass     = Array.from(new Array(goals), (x,i) => 6);
        shoot   = Array.from(new Array(goals), (x,i) => 6);
        fwd     = Array.from(new Array(goals), (x,i) => 2);
        
        log0 = await logUtils.encodeLog(encodeLog, goals, ass, shoot, fwd, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHomeSt, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkillsDefault, trainingPointsInit);
        
        logFinal = await training.computeTrainingPoints(log0, log0)
        // expect: POINTS_FOR_HAVING_PLAYED(10) + GOALS_BY_MIDS(5 * 5) - GOALS_OPPONENT(5)  
        expected = [30, 30];
        for (team = 0; team < 2; team++) {
            points = await encodeLog.getTrainingPoints(logFinal[team]).should.be.fulfilled;
            points.toNumber().should.be.equal(expected[team]);
        }
    });    

    it('training points with many goals by defs with assists', async () => {
        goals = 5;
        ass     = Array.from(new Array(goals), (x,i) => 6);
        shoot   = Array.from(new Array(goals), (x,i) => 2);
        fwd     = Array.from(new Array(goals), (x,i) => 1);
        
        log0 = await logUtils.encodeLog(encodeLog, goals, ass, shoot, fwd, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHomeSt, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkillsDefault, trainingPointsInit);
        
        logFinal = await training.computeTrainingPoints(log0, log0)
        // expect: POINTS_FOR_HAVING_PLAYED(10) + GOALS_BY_DEFS(4 * 5) + ASSISTS(3*5) - GOALS_OPPONENT(5)  
        expected = [50, 50];
        for (team = 0; team < 2; team++) {
            points = await encodeLog.getTrainingPoints(logFinal[team]).should.be.fulfilled;
            points.toNumber().should.be.equal(expected[team]);
        }
    });    

    it('training points with many goals with a winner at home', async () => {
        win = 0;
        isHome = true;

        goals = 5;
        ass     = Array.from(new Array(goals), (x,i) => 10);
        shoot   = Array.from(new Array(goals), (x,i) => 10);
        fwd     = Array.from(new Array(goals), (x,i) => 3);
        log0 = await logUtils.encodeLog(encodeLog, goals, ass, shoot, fwd, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHome, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, win, teamSumSkillsDefault, trainingPointsInit);

        goals = 4;
        ass     = Array.from(new Array(goals), (x,i) => 10);
        shoot   = Array.from(new Array(goals), (x,i) => 10);
        fwd     = Array.from(new Array(goals), (x,i) => 3);
        log1 = await logUtils.encodeLog(encodeLog, goals, ass, shoot, fwd, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHome, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, win, teamSumSkillsDefault, trainingPointsInit);
            
        logFinal = await training.computeTrainingPoints(log0, log1)
        // expect: POINTS_FOR_HAVING_PLAYED(10) + WIN_AT_HOME(11) + GOALS_BY_ATTACKERS(4 * 5) - GOALS_OPPONENT(4)  
        // expect: POINTS_FOR_HAVING_PLAYED(10) + GOALS_BY_ATTACKERS(4 * 4) - GOALS_OPPONENT(5)  
        expected = [37, 21];
        for (team = 0; team < 2; team++) {
            points = await encodeLog.getTrainingPoints(logFinal[team]).should.be.fulfilled;
            points.toNumber().should.be.equal(expected[team]);
        }
    });    

    it('training points with many goals with a winner away', async () => {
        win = 1;
        isHome = true;

        goals = 5;
        ass     = Array.from(new Array(goals), (x,i) => 10);
        shoot   = Array.from(new Array(goals), (x,i) => 10);
        fwd     = Array.from(new Array(goals), (x,i) => 3);
        log0 = await logUtils.encodeLog(encodeLog, goals, ass, shoot, fwd, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHome, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, win, teamSumSkillsDefault, trainingPointsInit);

        goals = 6;
        ass     = Array.from(new Array(goals), (x,i) => 10);
        shoot   = Array.from(new Array(goals), (x,i) => 10);
        fwd     = Array.from(new Array(goals), (x,i) => 3);
        log1 = await logUtils.encodeLog(encodeLog, goals, ass, shoot, fwd, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHome, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, win, teamSumSkillsDefault, trainingPointsInit);
            
        logFinal = await training.computeTrainingPoints(log0, log1)
        // expect: POINTS_FOR_HAVING_PLAYED(10) + GOALS_BY_ATTACKERS(4 * 5) - GOALS_OPPONENT(6)  
        // expect: POINTS_FOR_HAVING_PLAYED(10) + WIN_AWAY(22) + GOALS_BY_ATTACKERS(4 * 6) - GOALS_OPPONENT(5)  
        expected = [24, 51];
        for (team = 0; team < 2; team++) {
            points = await encodeLog.getTrainingPoints(logFinal[team]).should.be.fulfilled;
            points.toNumber().should.be.equal(expected[team]);
        }
    });    
    
    it('training points with no goals but cards', async () => {
        outGames    = [4, 6];
        types       = [RED_CARD, RED_CARD];
        yellows1    = [3, 7];
        yellows2    = [1, 2];
        
        log0 = await logUtils.encodeLog(encodeLog, nGoals = 0, assistersIdx, shootersIdx, shooterForwardPos, penalties,
            outGames, outOfGameRounds, types, yellowCardedDidNotFinish1stHalf,
            isHomeSt, ingameSubs1, ingameSubs2, yellows1, yellows2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkillsDefault, trainingPointsInit);
        
        logFinal = await training.computeTrainingPoints(log0, log0)
        // expect: POINTS_FOR_HAVING_PLAYED(10) + cleanSheet(23+8) - REDS(3*2) - YELLOWS(4) 
        expected = [31, 31];
        for (team = 0; team < 2; team++) {
            points = await encodeLog.getTrainingPoints(logFinal[team]).should.be.fulfilled;
            points.toNumber().should.be.equal(expected[team]);
        }
    });    
    
    it('training points with many goals by attackers... and different teamSumSkills', async () => {
        // first get the resulting Traning points with teamSkills difference: [25, 25]
        goals = 5;
        ass     = Array.from(new Array(goals), (x,i) => 10);
        shoot   = Array.from(new Array(goals), (x,i) => 10);
        fwd     = Array.from(new Array(goals), (x,i) => 3);
        
        log0 = await logUtils.encodeLog(encodeLog, goals, ass, shoot, fwd, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHomeSt, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkillsDefault, trainingPointsInit);
        
        logFinal = await training.computeTrainingPoints(log0, log0)
        // expect: POINTS_FOR_HAVING_PLAYED(10) + GOALS_BY_ATTACKERS(4 * 5) - GOALS_OPPONENT(5)  
        expected = [25, 25];
        for (team = 0; team < 2; team++) {
            points = await encodeLog.getTrainingPoints(logFinal[team]).should.be.fulfilled;
            points.toNumber().should.be.equal(expected[team]);
        }

        // second: get the resulting Traning points with teamSkills difference
        teamSumSkills = 1000;
        log0 = await logUtils.encodeLog(encodeLog, goals, ass, shoot, fwd, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHomeSt, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkills, trainingPointsInit);
        teamSumSkills = 2000;
        log1 = await logUtils.encodeLog(encodeLog, goals, ass, shoot, fwd, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHomeSt, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkills, trainingPointsInit);
            
        logFinal = await training.computeTrainingPoints(log0, log1)
        expected = [50, 12];
        for (team = 0; team < 2; team++) {
            points = await encodeLog.getTrainingPoints(logFinal[team]).should.be.fulfilled;
            points.toNumber().should.be.equal(expected[team]);
        }
        // third: same as above but inverse
        teamSumSkills = 2000;
        log0 = await logUtils.encodeLog(encodeLog, goals, ass, shoot, fwd, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHomeSt, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkills, trainingPointsInit);
        teamSumSkills = 1000;
        log1 = await logUtils.encodeLog(encodeLog, goals, ass, shoot, fwd, penalties,
            outOfGames, outOfGameRounds, typesOutOfGames, yellowCardedDidNotFinish1stHalf,
            isHomeSt, ingameSubs1, ingameSubs2, yellowCards1, yellowCards2, 
            halfTimeSubstitutions, nDefs1, nDefs2, nTot, winner, teamSumSkills, trainingPointsInit);
            
        logFinal = await training.computeTrainingPoints(log0, log1)
        expected = [12, 50];
        for (team = 0; team < 2; team++) {
            points = await encodeLog.getTrainingPoints(logFinal[team]).should.be.fulfilled;
            points.toNumber().should.be.equal(expected[team]);
        }

    });    
});