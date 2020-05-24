import React, { useState } from 'react';
import { Table } from 'semantic-ui-react';
import Config from '../../Config';
import RoleCard from './RoleCard';

const assetsJSON = require("../../contracts/Assets.json");

const RelayCard = (props) => {
    const {web3, account, assetsAddress} = props;
    const [relay, setRelay] = useState("");

    const assetsContract = new web3.eth.Contract(assetsJSON.abi, assetsAddress);
    assetsContract.methods.relay().call()
        .then(setRelay)
        .catch(error => { setRelay("") })

    const setAddress = (address) => {
        assetsContract.methods.setRelay(address).send({ from: account, gasPrice: Config.gasPrice })
            .on('error', (error, receipt) => { console.error(error) });
    }

    return (
        <Table.Row>
            <Table.Cell>relay</Table.Cell>
            <Table.Cell>{relay}</Table.Cell>
            {/* <RoleCard onChange={setAddress} disabled={account ? false : true} /> */}
        </Table.Row>
    )
}

export default RelayCard;