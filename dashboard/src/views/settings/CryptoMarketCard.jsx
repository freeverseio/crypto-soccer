import React, { useState } from 'react';
import { Table } from 'semantic-ui-react';

const assetsJSON = require("../../contracts/Assets.json");

const CryptoMarketCard = (props) => {
    const {web3, assetsAddress} = props;
    const [cryptoMarket, setCryptoMarket] = useState("");

    const assetsContract = new web3.eth.Contract(assetsJSON.abi, assetsAddress);
    assetsContract.methods.cryptoMktAddr().call()
        .then(setCryptoMarket)
        .catch(error => { setCryptoMarket("") })

    return (
        <Table.Row>
            <Table.Cell>cryptoMarket</Table.Cell>
            <Table.Cell>{cryptoMarket}</Table.Cell>
        </Table.Row>
    )
}

export default CryptoMarketCard;