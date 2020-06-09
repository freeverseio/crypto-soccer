import React, { useState, useEffect } from 'react';
import { Table } from 'semantic-ui-react';
import Config from '../../Config';
import RoleCard from './RoleCard';
import IntArrayCard from './IntArrayCard';

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

    return (
        <Table.Row>
            <Table.Cell>DeActivate and Activate</Table.Cell>
            <Table.Cell>{toIntArray(activateIds)}</Table.Cell>
            <Table.Cell>
                <IntArrayCard onChange={setActivateIds}/>
            </Table.Cell>
        </Table.Row>                
    )
}

export default ActivateCard;