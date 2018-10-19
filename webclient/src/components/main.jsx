import React, { PureComponent } from 'react';
import { Card, Form, Segment, Divider } from 'semantic-ui-react'
import TeamPlayerTable from './team_players_table';

class Main extends PureComponent {
    constructor(props){
        super(props);
        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);

        this.state = {
            name: ''
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