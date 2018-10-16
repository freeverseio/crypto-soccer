import Web3 from 'web3';

export default class EthLeagueManager {
    static async createAsync(provider, contractJSON) {
        const web3 = new Web3(provider);
        const networkId = await web3.eth.net.getId();
        const contract = new web3.eth.Contract(
            contractJSON.abi,
            contractJSON.networks[networkId].address
        )

        return new EthLeagueManager(contract)
    }

    constructor(contract) {
        this.contract = contract
        this.address = '0x82cC3f53b9DD7Fc8F546DB9eBC497b8D69B1AebA';
    }

    async createTeam(name) {
        this.contract.methods.test_createTeam(name).send({
            from: this.address,
            gasPrice: 20000000000,
            gas: 6721975
        });
    }

    async countTeams() {
        return await this.contract.methods.test_getNCreatedTeams.call().call({ from: this.address });
    }
}