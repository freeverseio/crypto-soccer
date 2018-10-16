import Web3 from 'web3';

export default class EthLeagueManager {
    static async createAsync(provider, contractJSON) {
        const web3 = new Web3(provider);
        const networkId = await web3.eth.net.getId();
        const contract = new web3.eth.Contract(
            contractJSON.abi,
            contractJSON.networks[networkId].address
        )
        const addresses = await web3.eth.getAccounts();

        return new EthLeagueManager(contract, addresses[0]);
    }

    constructor(contract, address) {
        this.contract = contract
        this.address = address;
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
