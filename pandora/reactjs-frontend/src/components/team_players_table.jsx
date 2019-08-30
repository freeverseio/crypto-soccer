import React, { PureComponent } from 'react';
import { Label, Table, Accordion, Icon } from 'semantic-ui-react'

// const roleNames = ['Goalkeeper','Defender','Midfielder','Forward'];

class TeamPlayersTable extends PureComponent {
    state = { active: false };

    render() {
        const { team } = this.props;
        const { active } = this.state;

        const players = () => {
            console.log(team)
            if (!team)
                return <Table.Row />

            return (
                team.map(player => (
                    <Table.Row key={player.teamId}>
                        <Table.Cell>
                            <Label ribbon>{player.id}</Label>
                        </Table.Cell>
                        <Table.Cell>TODO</Table.Cell>
                        <Table.Cell>{player.defence}</Table.Cell>
                        <Table.Cell>{player.speed}</Table.Cell>
                        <Table.Cell>{player.pass}</Table.Cell>
                        <Table.Cell>{player.shoot}</Table.Cell>
                        <Table.Cell>{player.endurance}</Table.Cell>
                        <Table.Cell>TODO</Table.Cell>
                    </Table.Row>
                ))
            )
        }

        return (
            <Accordion>
                <Accordion.Title active={active} onClick={() => this.setState({active: !active})} >
                    <Icon name='dropdown' />
                    Players
                </Accordion.Title>
                <Accordion.Content active={active}>
                    <Table compact>
                        <Table.Header>
                            <Table.Row>
                                <Table.HeaderCell></Table.HeaderCell>
                                <Table.HeaderCell>Age</Table.HeaderCell>
                                <Table.HeaderCell>Defense</Table.HeaderCell>
                                <Table.HeaderCell>Speed</Table.HeaderCell>
                                <Table.HeaderCell>Pass</Table.HeaderCell>
                                <Table.HeaderCell>Shoot</Table.HeaderCell>
                                <Table.HeaderCell>Endurance</Table.HeaderCell>
                                <Table.HeaderCell>Role</Table.HeaderCell>
                            </Table.Row>
                        </Table.Header>
                        <Table.Body>{players()}</Table.Body>
                    </Table>
                </Accordion.Content>
            </Accordion>
        )
    }
}

export default TeamPlayersTable;