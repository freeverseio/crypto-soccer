import React, { PureComponent } from 'react';
import { Segment, Divider, Button, Icon, Grid, Header, GridColumn } from 'semantic-ui-react'
import TeamPlayerTable from './team_players_table';
import TeamSelect from './team_select';
import TeamCreator from './team_creator';
import TeamList from './team_list';

class Main extends PureComponent {
    constructor(props){
        super(props);
        this.playGame = this.playGame.bind(this);

        this.state = {};
        this.teamA = -1;
        this.teamB = -1;
    }

    playGame(){
        const { testingFacade } = this.props;

        if (this.teamA < 0) return;
        if (this.teamB < 0) return;

        testingFacade.playGame(this.teamA, this.teamB)
        .then(result => this.setState({result: result}));
    }

    render() {
        const { testingFacade, teams } = this.props;
        const { team, result } = this.state;

        return (
            <Segment>
                <Segment>
                    <Grid relaxed>
                        <Grid.Row>
                            <Grid.Column width={6}>
                                <TeamSelect fluid placeholder='Select team A' teams={teams} onChange={(e, data) => this.teamA = data.value} />
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
                                <TeamSelect fluid placeholder='Select team B' teams={teams} onChange={(e, data) => this.teamB = data.value} />
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
                <TeamCreator contract={testingFacade}/>
                <Divider />
                <TeamList teams={teams} onChange={team => this.setState({team})} />
                <Divider />
                <TeamPlayerTable team={team} />
            </Segment>
        )
    }
}

export default Main;