import React, { useState } from 'react';
import { Table } from 'semantic-ui-react';
import Config from '../../Config';
import RoleCard from './RoleCard';

const assetsJSON = require("../../contracts/Assets.json");

const CryptoMarketCard = (props) => {
    const {web3, assetsAddress, account} = props;
    const [cryptoMarket, setCryptoMarket] = useState("");

    const assetsContract = new web3.eth.Contract(assetsJSON.abi, assetsAddress);
    assetsContract.methods.cryptoMktAddr().call()
        .then(setCryptoMarket)
        .catch(error => { setCryptoMarket("") })

    const setAddress = (address) => {
        assetsContract.methods.setCryptoMarketAddress(address).send({ from: account, gasPrice: Config.gasPrice })
            .on('error', (error, receipt) => { console.error(error) });
    }

    return (
        <Table.Row>
            <Table.Cell>cryptoMarket</Table.Cell>
            <Table.Cell>{cryptoMarket}</Table.Cell>
            {/* <RoleCard onChange={setAddress} disabled={account ? false : true} /> */}
        </Table.Row>
    )
}

export default CryptoMarketCard;