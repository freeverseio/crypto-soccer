// import leagueJSON from '../contracts/League.json';
import testingJSON from '../contracts/Testing.json';

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
        this.contract.methods.test_createTeam(name).send({
            from: this.address,
            gasPrice: 20000000000,
            gas: 6721975
        });
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
        for (let i=0; i < this.skillNumber ; i++)
            result[i] = unixMonthToAge(result[i]);

        return result;
    }
}
