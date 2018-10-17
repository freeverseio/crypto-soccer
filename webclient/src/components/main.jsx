import React, { PureComponent } from 'react';
import { Card, Form, Segment, Divider } from 'semantic-ui-react'
import TeamCard from './team_card';
import TeamPlayerTable from './team_players_table';

class Main extends PureComponent {
    constructor(props){
        super(props);
        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);

        this.state = {
            name: '',
        }
    }

    handleSubmit(event) {
        const { ethLeagueManager } = this.props;
        ethLeagueManager.createTeam(this.state.name);
        this.setState({name: ''});
        event.preventDefault();
    }

    handleChange(event) {
        this.setState({ name: event.target.value });
    }

    teamCardList = () => {
        const { teams, ethLeagueManager } = this.props;

        return (
            teams.map((team, key) => {
                return (
                    <TeamCard
                        key={key}
                        index={key}
                        name={team}
                        ethLeagueManager={ethLeagueManager}
                        onClick={index => this.setState({ selectedIndex: index })}
                    />
                )
            })
        )
    }

    render() {
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
                    {this.teamCardList()}
                </Card.Group>
                <Divider />
                {/* <TeamPlayerTable index={selectedIndex}/> */}
            </Segment>
        )
    }
}

export default Main;