import React, { useState } from 'react';
import { Table } from 'semantic-ui-react';
import Config from '../../Config';
import RoleCard from './RoleCard';

const assetsJSON = require("../../contracts/Assets.json");

const MarketCard = (props) => {
    const { web3, account, assetsAddress } = props;
    const [market, setMarket] = useState("");

    const assetsContract = new web3.eth.Contract(assetsJSON.abi, assetsAddress);
    assetsContract.methods.market().call()
        .then(setMarket)
        .catch(error => { setMarket("") })

    const setAddress = (address) => {
        assetsContract.methods.setMarket(address).send({ from: account, gasPrice: Config.gasPrice })
            .on('error', (error, receipt) => { console.error(error) });
    }

    return (
        <React.Fragment>
            <Table.Row>
                <Table.Cell>market</Table.Cell>
                <Table.Cell>{market}</Table.Cell>
                {/* <RoleCard onChange={setAddress} disabled={account ? false : true} /> */}
            </Table.Row>
        </React.Fragment>

    )
}

export default MarketCard;