import React, { useState, useEffect } from 'react';
import { Table } from 'semantic-ui-react';
import Config from '../../Config';
import RoleCard from './RoleCard';

const MarketCard = ({ account, assetsContract }) => {
    const [market, setMarket] = useState();

    useEffect(() => {
        assetsContract.methods.market().call()
            .then(setMarket)
            .catch(error => {
                console.error(error);
                setMarket("error");
            });
    }, [assetsContract]);

    const setAddress = (address) => {
        assetsContract.methods.setMarket(address).send({ from: account, gasPrice: Config.gasPrice })
        .catch(console.error);
    }

    return (
        <React.Fragment>
            <Table.Row>
                <Table.Cell>market</Table.Cell>
                <Table.Cell>{market}</Table.Cell>
                <Table.Cell>
                    <RoleCard onChange={setAddress} account={account}/>
                </Table.Cell>
            </Table.Row>
        </React.Fragment>

    )
}

export default MarketCard;