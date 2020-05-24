import React, { useState } from 'react';
import { Table } from 'semantic-ui-react';

const assetsJSON = require("../../contracts/Assets.json");

const MarketCard = (props) => {
    const {web3, assetsAddress} = props;
    const [market, setMarket] = useState("");

    const assetsContract = new web3.eth.Contract(assetsJSON.abi, assetsAddress);
    assetsContract.methods.market().call()
        .then(setMarket)
        .catch(error => { setMarket("") })

    return (
        <Table.Row>
            <Table.Cell>market</Table.Cell>
            <Table.Cell>{market}</Table.Cell>
        </Table.Row>
    )
}

export default MarketCard;