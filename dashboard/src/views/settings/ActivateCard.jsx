import React, { useState, useEffect } from 'react';
import { Table } from 'semantic-ui-react';
import Config from '../../Config';
import RoleCard from './RoleCard';
import IntArrayCard from './IntArrayCard';
import AddrCard from './AddrCard';
import { Button, Input } from 'semantic-ui-react';

function toIntArray(x) {
    var array = [];
    var res = x.split(",");
    for (var elem of res) {
        var intElem = parseInt(elem);
        if (!intElem) return [];
        array.push(intElem)
    }
    return array;
}

const ActivateCard = ({account, proxyContract}) => {
    const [activateIds, setActivateIds] = useState("");
    const [deactivateIds, setDeactivateIds] = useState("");
    const [directoryAddr, setDirectoryAddr] = useState("");

    const submitActivation = (deactivateIds, activateIds, directoryAddr) => {
        var intDeact = toIntArray(deactivateIds);
        var intAct = toIntArray(activateIds);
        proxyContract.methods.upgrade(
            intDeact, 
            intAct, 
            directoryAddr
        ).send({ from: account, gasPrice: Config.gasPrice })
        .catch(console.error);
    }

    return (
        <Table.Row>
            <Table.Cell>
                <Button size='mini' color='red' onClick={() => { 
                    submitActivation(deactivateIds, activateIds, directoryAddr) 
                }}
                >Submit</Button>
            </Table.Cell>
            <Table.Cell>ActivateIds:   <IntArrayCard onChange={setActivateIds}/></Table.Cell>
            <Table.Cell>DeactivateIds: <IntArrayCard onChange={setDeactivateIds}/></Table.Cell>
            <Table.Cell>New Dir Addr:<AddrCard onChange={setDirectoryAddr}/></Table.Cell>
        </Table.Row> 
    )
}

export default ActivateCard;