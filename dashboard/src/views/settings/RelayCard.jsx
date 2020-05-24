import React, { useState } from 'react';
import { Table } from 'semantic-ui-react';

const assetsJSON = require("../../contracts/Assets.json");

const RelayCard = (props) => {
    const {web3, assetsAddress} = props;
    const [relay, setRelay] = useState("");

    const assetsContract = new web3.eth.Contract(assetsJSON.abi, assetsAddress);
    assetsContract.methods.relay().call()
        .then(setRelay)
        .catch(error => { setRelay("") })

    return (
        <Table.Row>
            <Table.Cell>relay</Table.Cell>
            <Table.Cell>{relay}</Table.Cell>
        </Table.Row>
    )
}

export default RelayCard;