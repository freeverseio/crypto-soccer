import React, { useState } from 'react';
import { Table } from 'semantic-ui-react';
import Config from '../../Config';
import RoleCard from './RoleCard';

const marketJSON = require("../../contracts/Market.json");

const CryptoMarketCard = (props) => {
    const {web3, marketAddress, account} = props;
    const [cryptoMarket, setCryptoMarket] = useState("");

    const marketContract = new web3.eth.Contract(marketJSON.abi, marketAddress);
    marketContract.methods.cryptoMktAddr().call()
        .then(setCryptoMarket)
        .catch(error => { setCryptoMarket("") })

    const setAddress = (address) => {
        marketContract.methods.setCryptoMarketAddress(address).send({ from: account, gasPrice: Config.gasPrice })
            .on('error', (error, receipt) => { console.error(error) });
    }

    return (
        <Table.Row>
            <Table.Cell>cryptoMarket</Table.Cell>
            <Table.Cell>{cryptoMarket}</Table.Cell>
            <Table.Cell>
                <RoleCard onChange={setAddress} disabled={account ? false : true} />
            </Table.Cell>
        </Table.Row>
    )
}

export default CryptoMarketCard;