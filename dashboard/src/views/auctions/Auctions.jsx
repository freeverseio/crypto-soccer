import React from 'react';
import { useQuery } from '@apollo/react-hooks';
import gql from 'graphql-tag';
import { Table } from 'semantic-ui-react';
import Config from '../../Config';

const GET_AUCTIONS = gql`
{
    allAuctions {
        nodes {
            id
            playerId
            price
            validUntil
            seller
            state
            stateExtra
        }
    }
}
`;

export default () => {
    const { loading, error, data } = useQuery(GET_AUCTIONS, {
         pollInterval: Config.auctions_polling_ms,
    });

    if (loading) return 'Loading...';
    if (error) return `Error! ${error.message}`;

    const orders = data.allAuctions.nodes ? data.allAuctions.nodes : [];
    console.log(orders);

    return (
        <Table sortable celled fixed>
            <Table.Header>
                <Table.Row>
                    <Table.HeaderCell>id</Table.HeaderCell>
                    <Table.HeaderCell>playerId</Table.HeaderCell>
                    <Table.HeaderCell>price</Table.HeaderCell>
                    <Table.HeaderCell>validUntil</Table.HeaderCell>
                    <Table.HeaderCell>seller</Table.HeaderCell>
                    <Table.HeaderCell>state</Table.HeaderCell>
                    <Table.HeaderCell>stateExtra</Table.HeaderCell>
                </Table.Row>
            </Table.Header>
            <Table.Body>
                {orders.map(order => (
                    <Table.Row key={order.id}>
                        <Table.Cell>{order.id}</Table.Cell>
                        <Table.Cell>{order.playerId}</Table.Cell>
                        <Table.Cell>{order.price}</Table.Cell>
                        <Table.Cell>{order.validUntil}</Table.Cell>
                        <Table.Cell>{order.seller}</Table.Cell>
                        <Table.Cell>{order.state}</Table.Cell>
                        <Table.Cell>{order.stateExtra}</Table.Cell>
                    </Table.Row>
                ))}
            </Table.Body>
        </Table>
    )
}