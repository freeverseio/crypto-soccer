import React, { useState, useEffect } from 'react';
import { Table } from 'semantic-ui-react';
import Config from '../../Config';
import RoleCard from './RoleCard';

const RelayCard = ({account, assetsContract}) => {
    const [relay, setRelay] = useState();

    useEffect(() => {
        assetsContract.methods.relay().call()
            .then(setRelay)
            .catch(error => {
                console.error(error);
                setRelay("error");
            });
    }, [assetsContract]);

    const setAddress = (address) => {
        assetsContract.methods.setRelay(address).send({ from: account, gasPrice: Config.gasPrice })
        .catch(console.error);
    }

    return (
        <Table.Row>
            <Table.Cell>relay</Table.Cell>
            <Table.Cell>{relay}</Table.Cell>
            <Table.Cell>
                <RoleCard onChange={setAddress} account={account}/>
            </Table.Cell>
        </Table.Row>
    )
}

export default RelayCard;