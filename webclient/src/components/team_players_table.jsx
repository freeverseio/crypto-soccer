import React, { PureComponent } from 'react';
import { Label, Table } from 'semantic-ui-react'

const roleNames = ['Goalkeeper','Defender','Midfielder','Forward'];

class TeamPlayersTable extends PureComponent {
    render() {
        const { team } = this.props;

        const players = () => {
            if (!team)
                return <Table.Cell />

            return (
                team.players.map(player => (
                    <Table.Row key={player.index}>
                        <Table.Cell>
                            <Label ribbon>{player.index}</Label>
                        </Table.Cell>
                        <Table.Cell>{player.skills[0]}</Table.Cell>
                        <Table.Cell>{player.skills[1]}</Table.Cell>
                        <Table.Cell>{player.skills[2]}</Table.Cell>
                        <Table.Cell>{player.skills[3]}</Table.Cell>
                        <Table.Cell>{player.skills[4]}</Table.Cell>
                        <Table.Cell>{player.skills[5]}</Table.Cell>
                        <Table.Cell>{roleNames[player.skills[6]]}</Table.Cell>
                    </Table.Row>
                ))
            )
        }

        return (
            <Table celled>
                <Table.Header>
                    <Table.Row>
                        <Table.HeaderCell>Team {team && team.name}</Table.HeaderCell>
                        <Table.HeaderCell>Age</Table.HeaderCell>
                        <Table.HeaderCell>Defense</Table.HeaderCell>
                        <Table.HeaderCell>Speed</Table.HeaderCell>
                        <Table.HeaderCell>Pass</Table.HeaderCell>
                        <Table.HeaderCell>Shoot</Table.HeaderCell>
                        <Table.HeaderCell>Endurance</Table.HeaderCell>
                        <Table.HeaderCell>Role</Table.HeaderCell>
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