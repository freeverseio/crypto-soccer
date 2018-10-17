import React, { PureComponent } from 'react';
import { Label, Table } from 'semantic-ui-react'

class TeamPlayersTable extends PureComponent {
    render() {
        const { team } = this.props;

        const players = () => {
            if (!team)
                return <div />

            return (
                team.players.map(player => (
                    <Table.Row>
                        <Table.Cell>
                            <Label ribbon>{player.player}</Label>
                        </Table.Cell>
                        <Table.Cell>Cell</Table.Cell>
                        <Table.Cell>Cell</Table.Cell>
                    </Table.Row>
                ))
            )
        }

        return (
            <Table celled>
                <Table.Header>
                    <Table.Row>
                        <Table.HeaderCell>Team {team && team.name}</Table.HeaderCell>
                        <Table.HeaderCell>Header</Table.HeaderCell>
                        <Table.HeaderCell>Header</Table.HeaderCell>
                    </Table.Row>
                </Table.Header>

                <Table.Body>
                    {players()}
               </Table.Body>
            </Table>
        )
    }
}

export default TeamPlayersTable;