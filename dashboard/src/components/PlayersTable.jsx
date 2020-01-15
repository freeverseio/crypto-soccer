import React from 'react';
import { Table } from 'semantic-ui-react'

export default function PlayersTable(props) {
    const players = props.players ? props.players : [];

    return (
        <Table color='grey' inverted selectable >
            <Table.Header>
                <Table.Row>
                    <Table.HeaderCell>Name</Table.HeaderCell>
                    <Table.HeaderCell>Defence</Table.HeaderCell>
                    <Table.HeaderCell>Pass</Table.HeaderCell>
                    <Table.HeaderCell>Speed</Table.HeaderCell>
                    <Table.HeaderCell>Endurance</Table.HeaderCell>
                </Table.Row>
            </Table.Header>
            <Table.Body >
                {
                    players.map(player => (
                        <Table.Row key={player.playerId}>
                            <Table.Cell>{player.name}</Table.Cell>
                            <Table.Cell>{player.defence}</Table.Cell>
                            <Table.Cell>{player.pass}</Table.Cell>
                            <Table.Cell>{player.speed}</Table.Cell>
                            <Table.Cell>{player.endurance}</Table.Cell>
                        </Table.Row>
                    ))
                }
            </Table.Body>
        </Table>
    );
}