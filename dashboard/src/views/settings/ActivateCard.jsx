import React, { useState, useEffect } from 'react';
import { Table } from 'semantic-ui-react';
import Config from '../../Config';
import RoleCard from './RoleCard';
import IntArrayCard from './IntArrayCard';

const ActivateCard = ({account, proxyContract}) => {
    const [activateIds, setActivateIds] = useState("");

    return (
        <Table.Row>
            <Table.Cell>DeActivate and Activate</Table.Cell>
            <Table.Cell>{activateIds}</Table.Cell>
            <Table.Cell>
                <IntArrayCard onChange={setActivateIds}/>
            </Table.Cell>
        </Table.Row>                
    )
}

export default ActivateCard;