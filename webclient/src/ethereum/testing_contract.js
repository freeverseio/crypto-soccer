// import leagueJSON from '../contracts/League.json';
import testingJSON from '../contracts/Testing.json';
import f from '../jsCommons/functions';

export const createTestingContract = async web3 => {
    const contractJSON = testingJSON;
    const networkId = await web3.eth.net.getId();
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
}
