import React, { useState, useEffect } from 'react';
import { Table } from 'semantic-ui-react';
import Config from '../../Config';
import RoleCard from './RoleCard';

const CryptoMarketCard = ({account, marketContract}) => {
    const [cryptoMarket, setCryptoMarket] = useState("");

    useEffect(() => {
        marketContract.methods.cryptoMktAddr().call()
            .then(setCryptoMarket)
            .catch(error => {
                console.error(error);
                setCryptoMarket("error");
            });
    },[marketContract]);

    const setAddress = (address) => {
        marketContract.methods.setCryptoMarketAddress(address).send({ from: account, gasPrice: Config.gasPrice })
        .catch(console.error);
    }

    return (
        <Table.Row>
            <Table.Cell singleLine>Crypto Market Role</Table.Cell>
            <Table.Cell>{cryptoMarket}</Table.Cell>
            <Table.Cell>
                <RoleCard onChange={setAddress} account={account}/>
            </Table.Cell>
        </Table.Row>
    )
}

export default CryptoMarketCard;