import React, { PureComponent } from 'react';
import { Card, Form, Segment, Divider, Button, Icon, Grid } from 'semantic-ui-react'
import TeamPlayerTable from './team_players_table';
import TeamSelect from './team_select';

class Main extends PureComponent {
    constructor(props){
        super(props);
        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.playGame = this.playGame.bind(this);

        this.state = {
            name: '',
        }
    }

    handleSubmit(event) {
        const { testingFacade } = this.props;
        testingFacade.createTeam(this.state.name);
        this.setState({name: ''});
        event.preventDefault();
    }

    handleChange(event) {
        this.setState({ name: event.target.value });
    }

    playGame(){
        const { testingFacade } = this.props;

        if (!this.teamA) return;
        if (!this.teamB) return;

        testingFacade.playGame(this.teamA, this.teamB)
        .then(console.log);
    }

    render() {
        const { teams } = this.props;
        const { team } = this.state;

        const cardList = teams.map(team => (
            <Card
                key={team.index}
                image='https://upload.wikimedia.org/wikipedia/it/0/07/Fc_barcelona.png'
                header={team.name}
                meta='Team'
                description={team.name + " is amazing!"}
                onClick={() => this.setState({ team })}
            />
        ));

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
                                <TeamSelect fluid placeholder='Select team B' teams={teams} onChange={(e, data) => this.teamB = data.value}/>
                            </Grid.Column>
                        </Grid.Row>
                    </Grid>
                </Segment>
                <Form onSubmit={this.handleSubmit}>
                    <Form.Group widths='equal'>
                        <Form.Input placeholder='Name' name='name' value={this.state.name} onChange={this.handleChange} />
                        <Form.Button content='Create Team' />
                    </Form.Group>
                </Form>
                <Divider />
                <Card.Group>
                    {cardList}
                </Card.Group>
                <Divider />
                <TeamPlayerTable team={team} />
            </Segment>
        )
    }
}

export default Main;