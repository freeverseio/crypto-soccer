import leagueJSON from '../contracts/League.json';
import testingJSON from '../contracts/Testing.json';

export default class EthLeagueManager {
    static async createAsync(web3, account) {
        const contractJSON = testingJSON;
        const networkId = await web3.eth.net.getId();
        const contract = new web3.eth.Contract(
            contractJSON.abi,
            contractJSON.networks[networkId].address
        )

        return new EthLeagueManager(contract, account);
    }

    constructor(contract, account) {
        this.contract = contract
        this.address = account;
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
}
