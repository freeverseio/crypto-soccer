import React from 'react';
import { useQuery } from '@apollo/react-hooks';
import gql from 'graphql-tag';
import { Table } from 'semantic-ui-react';
import Config from '../../Config';

const GET_PLAYSTORE_ORDERS = gql`
    {
        allPlaystoreOrders {
            nodes {
            orderId
            packageName
            productId
            purchaseToken
            playerId
            teamId
            state
            stateExtra
        }
    }
}
`;

export default () => {
    const { loading, error, data } = useQuery(GET_PLAYSTORE_ORDERS, {
         pollInterval: Config.iap_polling_ms,
    });

    if (loading) return 'Loading...';
    if (error) return `Error! ${error.message}`;

    const orders = data.allPlaystoreOrders.nodes ? data.allPlaystoreOrders.nodes : [];
    console.log(orders);

    return (
        <Table sortable celled fixed>
            <Table.Header>
                <Table.Row>
                    <Table.HeaderCell>orderId</Table.HeaderCell>
                    <Table.HeaderCell>packageName</Table.HeaderCell>
                    <Table.HeaderCell>productId</Table.HeaderCell>
                    <Table.HeaderCell>purchaseToken</Table.HeaderCell>
                    <Table.HeaderCell>playerId</Table.HeaderCell>
                    <Table.HeaderCell>teamId</Table.HeaderCell>
                    <Table.HeaderCell>state</Table.HeaderCell>
                    <Table.HeaderCell>stateExtra</Table.HeaderCell>
                </Table.Row>
            </Table.Header>
            <Table.Body>
                {orders.map(order => (
                    <Table.Row key={order.purchaseToken}>
                        <Table.Cell>{order.orderId}</Table.Cell>
                        <Table.Cell>{order.packageName}</Table.Cell>
                        <Table.Cell>{order.productId}</Table.Cell>
                        <Table.Cell>{order.purchaseToken}</Table.Cell>
                        <Table.Cell>{order.playerId}</Table.Cell>
                        <Table.Cell>{order.teamId}</Table.Cell>
                        <Table.Cell>{order.state}</Table.Cell>
                        <Table.Cell>{order.stateExtra}</Table.Cell>
                    </Table.Row>
                ))}
            </Table.Body>
        </Table>
    )
}