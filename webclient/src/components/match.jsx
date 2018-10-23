import React, { Component } from 'react';
import { Segment, Button, Icon, Grid, Header, GridColumn, Item, Progress, Image } from 'semantic-ui-react'
import TeamSelect from './team_select';
import { AttackEvent, DefendEvent, GoalEvent, BlockedEvent, ShootEvent } from './match_events';

class Match extends Component {
    constructor(props) {
        super(props);

        this.state = {
            teamA: -1,
            teamB: -1,
            result: [],
            events: [],
            totalEvents: 0,
            playing: false
        };
    }

    playGame() {
        const { contract } = this.props;
        const { teamA, teamB } = this.state;

        if (teamA < 0) return;
        if (teamB < 0) return;

        contract.playGame(teamA, teamB)
            .then(summary => {
                this.setState({ 
                    playing: true, 
                    totalEvents: summary.length,
                    result: [0,0],
                    events: [],
                 })
                const delta = 500;
                for (let i = 1; i <= summary.length; i++) {
                    setTimeout(() => {
                        const slice = summary.slice(0, i);

                        this.setState({
                            events: slice
                        })
                        if (i === summary.length)
                            this.setState({playing: false})
                    }, delta * i);
                }
            });
    }

    parseEvent = (key, event) => {
        if (event.type === "attack")
            return (<AttackEvent min={event.min} text={event.text} />);

        if (event.type === "defended")
            return (<DefendEvent text={event.text} />);

        if (event.type === "shot")
            return (<ShootEvent text={event.text} />);

        if (event.type === "gool")
            return (<GoalEvent text={event.text} />);

        if (event.type === "blocked")
            return (<BlockedEvent text={event.text} />);
    }

    render() {
        const { teams } = this.props;
        const { teamA, teamB, result, events, totalEvents, playing } = this.state;

        return (
            <React.Fragment>
                <Segment>
                    <Grid relaxed>
                        <Grid.Row>
                            <Grid.Column width={6}>
                                <TeamSelect fluid placeholder='Select team A' teams={teams} value={teamA} onChange={(_, data) => this.setState({ teamA: data.value })} />
                            </Grid.Column>
                            <Grid.Column width={4}>
                                <Button animated fluid loading={playing} disabled={playing} onClick={() => this.playGame()}>
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
                                <Header textAlign='center' as="h1">{result[0]}</Header>
                            </GridColumn>
                            <GridColumn width={8}>
                                <Header textAlign='center' as="h1">{result[1]}</Header>
                            </GridColumn>
                        </Grid.Row>

                    </Grid>
                </Segment>
                <Segment>
                    {playing && <Progress percent={100 * events.length / totalEvents} success />}
                    <Grid relaxed>
                        <Grid.Row>
                            <GridColumn width={4}>
                                <Item.Group divided>
                                    {events.slice(0).reverse().map((event, key) => (
                                        this.parseEvent(key, event)
                                    ))}
                                </Item.Group>
                            </GridColumn>
                            <GridColumn width={12}>
                                <Image src="http://www.codethislab.com/wp-content/uploads/2016/07/7c989a40659221.5787b557aa8e5.jpg" />
                            </GridColumn>
                        </Grid.Row>
                    </Grid>
                </Segment>
            </React.Fragment>
        );
    }
}

export default Match;