import React, { Component } from 'react';
import { Segment, Button, Icon, Grid, Header, GridColumn } from 'semantic-ui-react'
import TeamSelect from './team_select';

class Match extends Component {
    constructor(props) {
        super(props);

        this.state = {
            teamA: -1,
            teamB: -1,
        };
    }

    playGame() {
        const { contract } = this.props;
        const { teamA, teamB } = this.state;

        if (teamA < 0) return;
        if (teamB < 0) return;

        contract.playGame(teamA, teamB)
            .then(result => this.setState({ result: result }));
    }

    render() {
        const { teams } = this.props;
        const { teamA, teamB, result } = this.state;

        return (
            <Segment>
                <Grid relaxed>
                    <Grid.Row>
                        <Grid.Column width={6}>
                            <TeamSelect fluid placeholder='Select team A' teams={teams} value={teamA} onChange={(_, data) => this.setState({ teamA: data.value })} />
                        </Grid.Column>
                        <Grid.Column width={4}>
                            <Button animated fluid onClick={() => this.playGame()}>
                                <Button.Content visible>Play</Button.Content>
                                <Button.Content hidden>
                                    <Icon name='arrow right' />
                                </Button.Content>
                            </Button>
                        </Grid.Column>
                        <Grid.Column width={6}>
                            <TeamSelect fluid placeholder='Select team B' teams={teams} value={teamB} onChange={(_, data) => this.setState({ teamB: data.value })} />
                        </Grid.Column>
                    </Grid.Row>
                    <Grid.Row>
                        <GridColumn width={8}>
                            <Header textAlign='center' as="h1">{result && result[0]}</Header>
                        </GridColumn>
                        <GridColumn width={8}>
                            <Header textAlign='center' as="h1">{result && result[1]}</Header>
                        </GridColumn>
                    </Grid.Row>
                </Grid>
            </Segment>
        );
    }
}

export default Match;