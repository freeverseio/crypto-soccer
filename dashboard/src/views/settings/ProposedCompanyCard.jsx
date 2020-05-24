import React, { useState } from 'react';
import { Input, Button } from 'semantic-ui-react';
import Config from '../../Config';

const proxyJSON = require("../../contracts/Proxy.json");

const ProposedCompanyWidget = ({ web3, proxyAddress, account, proposedCompany }) => {
    const [address, setAddress] = useState("");

    const proxyContract = new web3.eth.Contract(proxyJSON.abi, proxyAddress);

    const accept = () => {
        proxyContract.methods.acceptCompany().send({
            from: account,
            gasPrice: Config.gasPrice,
        })
            .on('transactionHash', hash => { console.log(hash) })
            .on('confirmation', (confirmationNumber, receipt) => {
                console.log("confiramtion");
            })
            .on('receipt', (receipt) => {
                console.log(receipt);
            })
            .on('error', (error, receipt) => { console.log(error) });
    };

    const validAddress = web3.utils.isAddress(proposedCompany);

    return (
        <Input fluid
            size='mini'
            error={!validAddress}
            disabled
            icon='ethereum'
            iconPosition='left'
            placeholder={proposedCompany}
            value={address}
            onChange={event => setAddress(event.target.value)}
            action={
                <Button size='mini' color='red' onClick={accept} disabled={!validAddress}>Accept</Button>
            }
        />
    )
}

export default ProposedCompanyWidget;