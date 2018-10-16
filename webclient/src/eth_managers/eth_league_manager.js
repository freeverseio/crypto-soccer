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
    }

    async createTeam(address, name) {
        this.contract.test_createTeam(name).send({
            from: address,
            gasPrice: 20000000000,
            gas: 6721975
        });
    }
}