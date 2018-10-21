// import leagueJSON from '../contracts/League.json';
import testingJSON from '../contracts/Testing.json';
import f from '../jsCommons/functions';
import k from '../jsCommons/constants';

export const createTestingContract = async web3 => {
    const connected = await web3.eth.net.isListening();
    const contractJSON = testingJSON;
    const networkId = await web3.eth.net.getId();

    if (!contractJSON.networks[networkId])
        return;
    
    const contract = new web3.eth.Contract(
        contractJSON.abi,
        contractJSON.networks[networkId].address
    )

    return contract;
}

function unixMonthToAge(unixMonthOfBirth) {
    // in July 2018, we are at month 582 after 1970.
    const age = (582 - unixMonthOfBirth) / 12;
    return parseInt(age * 10) / 10;
}

export class TestingFacade {
    constructor(contract, account) {
        this.contract = contract
        this.address = account;
        this.skillNumber = 7;
        this.bitPerState = 14;
    }

    async createTeam(name) {
        console.log(this.contract)
        const instance = this.contract;
        const teamName = name;
        const playerBasename = name;
        const maxPlayersPerTeam = 11;
        const playerRoles = f.createAlineacion(4, 3, 3);

        const newTeamIdx = await instance.methods.test_getNCreatedTeams().call({ from: this.address });
        console.log("creating team: " + newTeamIdx + ", name " + teamName);
        await instance.methods.test_createTeam(teamName).send(
            {
                from: this.address,
                gasPrice: 20000000000,
                gas: 6721975
            }
        );
        const userChoice = 1;

        for (var p = 0; p < maxPlayersPerTeam; p++) {
            const thisName = playerBasename + p.toString();
            await instance.methods.test_createBalancedPlayer(
                thisName,
                newTeamIdx,
                userChoice,
                p,
                playerRoles[p]
            ).send(
                {
                    from: this.address,
                    gasPrice: 20000000000,
                    gas: 6721975
                }
            );
        }
        const nCreatedPlayers = await instance.methods.test_getNCreatedPlayers.call();
        console.log('Final nPlayers in the entire game = ' + nCreatedPlayers);
    }

    async countTeams() {
        return this.contract.methods.test_getNCreatedTeams().call({ from: this.address });
    }

    async teamName(index){
        return this.contract.methods.test_getTeamName(index).call({ from: this.address });
    }

    async playerSkills(teamIndex, index) {
        const serialized = await this.contract.methods.test_getStatePlayerInTeam(index, teamIndex).call({ from: this.address });
        const result = await this.contract.methods.test_decode(this.skillNumber, serialized, this.bitPerState).call({ from: this.address });
        result[0] = unixMonthToAge(result[0]);

        return result;
    }

    async playerName(teamIndex, index){
        const absIndex = 10 * teamIndex + index + 1 + teamIndex;
        return await this.contract.methods.test_getPlayerName(absIndex).call();
    }

    async playGame(teamAIndex, teamBIndex) {
        const seed = Math.floor(Math.random() * 10000);
        const tx = await this.contract.methods.test_playGame(teamAIndex, teamBIndex, seed).send(
            {
                from: this.address,
                gas: 6721975
            });

        const gameId = await this.contract.methods.test_getGameId(teamAIndex, teamBIndex, seed).call();

        // catches events and prints them out
        const gameEvents = this.catchGameResults(tx.events, gameId);
        
        console.log(gameEvents);
        //return result;
    }

    catchGameResults(logs, gameId) {
        const teamAttacksEvents = logs.TeamAttacks;
        let teamThatAttacks = teamAttacksEvents.map(attack => ({
            homeOrAway: attack.returnValues.homeOrAway,
            round: attack.returnValues.round
        }));

        const shootResultEvents = logs.ShootResult;
        let shootResult = shootResultEvents.map(shoot => ({
            isGoal: shoot.returnValues.isGoal,
            attackerIdx: shoot.returnValues.attackerIdx,
            round: shoot.returnValues.round
        }));
        
        const result = {
            teamThatAttacks: teamThatAttacks,
            shootResult: shootResult
        };
        this.printGameEvents(result)
        return result;
    }

    printGameEvents(gameEvents) {
        console.log("EVENTS: ");
        for (var r = 0; r < k.RoundsPerGame; r++) {
            // we add a bit of noise so that events are not always at minute 5,10,15...
            var rndNoise = Math.round(-2 + Math.floor(Math.random() * 4));
            var thisMinute = (r + 1) * 5 + rndNoise;
            var t = f.getEntryForAGivenRound(gameEvents.teamThatAttacks, r);
            console.log("Min " + thisMinute + ": Opportunity for team " + t[1] + "...");
            var result = f.getEntryForAGivenRound(gameEvents.shootResult, r);
            if (result.length == 0) { console.log("  ... well tackled by defenders, did not prosper!"); }
            else {
                console.log("  ... that leads to a shoot by attacker " + result[2]);
                if (result[1]) { console.log("  ... and GOAAAAL!") }
                else { console.log("  ... blocked by the goalkeeper!!"); }
            }
        }
    }
}
