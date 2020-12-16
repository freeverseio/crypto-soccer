import React, { useState, useEffect } from 'react';
import { Table } from 'semantic-ui-react';
import Config from '../../Config';

const proxyJSON = require("../../contracts/Proxy.json");
const multisigJSON = require("../../contracts/MultiSigWallet.json");

const MultiSigTable = ({ web3, account, proxyAddress }) => {
    const [seconds, setSeconds] = useState(0);
    const [owners, setOwners] = useState([]);
    const proxyContract = new web3.eth.Contract(proxyJSON.abi, proxyAddress);
    const multisigContract = new web3.eth.Contract(multisigJSON.abi, Config.multiSigAddress);

    useEffect(() => {
        multisigContract.methods.getOwners().call()
            .then(setOwners)
            .catch(error => {
                console.error(error);
                setOwners("error");
            });
    }, [multisigContract]);

    useEffect(() => {
        const interval = setInterval(() => {
            setSeconds(seconds => seconds + 1);
        }, 5000);
        return () => clearInterval(interval);
    },[seconds]);

return (
        <Table color='orange'>
            <Table.Header>
                <Table.Row>
                    <Table.HeaderCell width={1}></Table.HeaderCell>
                    <Table.HeaderCell width={1}></Table.HeaderCell>
                    <Table.HeaderCell width='six'></Table.HeaderCell>
                </Table.Row>
            </Table.Header>
            <Table.Body>
                <Table.Row>
                    <Table.Cell singleLine>multi sig address</Table.Cell>
                <Table.Cell>{Config.multiSigAddress}</Table.Cell>
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

export default MultiSigTable;