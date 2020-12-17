import React, { useState, useEffect } from 'react';
import { Table } from 'semantic-ui-react';
import Config from '../../Config';

const multisigJSON = require("../../contracts/MultiSigWallet.json");

const PermissionTable = ({ web3, account, proxyAddress }) => {
    const [seconds, setSeconds] = useState(0);
    const [owners, setOwners] = useState([]);
    const multisigContract = new web3.eth.Contract(multisigJSON.abi, Config.multiSigAddress);

    console.log(multisigContract)

    useEffect(() => {
        const interval = setInterval(() => {
            setSeconds(seconds => seconds + 1);
        }, 5000);
        return () => clearInterval(interval);
    },[seconds]);

    useEffect(() => {
        multisigContract.methods.getOwners().call()
            .then(setOwners)
            .catch(error => {
                console.error(error);
                setOwners("error");
            });
    }, [multisigContract]);

    return (
        <Table color='green' columns={2}>
            <Table.Header>
                <Table.Row>
                    <Table.HeaderCell colSpan='2' textAlign='center'>multisig wallet</Table.HeaderCell>
                </Table.Row>
            </Table.Header>
            <Table.Body>
                <Table.Row>
                    <Table.Cell>contract address</Table.Cell>
                    <Table.Cell>{multisigContract.options.address}</Table.Cell>
                </Table.Row>
                <Table.Row>
                    <Table.Cell singleLine>owners</Table.Cell>
                    <Table.Cell>
                        {owners.map(owner => owner + " ")}
                    </Table.Cell>
                </Table.Row>
            </Table.Body>
        </Table>
    )
}

export default PermissionTable;