import React, { PureComponent } from 'react';
import { Card, Form, Segment, Divider } from 'semantic-ui-react'
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

    render() {
        const { teams, ethLeagueManager } = this.props;

        const cardList = teams.map((team, key) => (
            <Card
                image='https://static.independent.co.uk/s3fs-public/thumbnails/image/2017/11/21/13/borat.jpg'
                header={team}
                meta='Team'
                description={team + " is amazing!"}
                onClick={() => console.log(team)}
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
                {/* <TeamPlayerTable index={selectedIndex}/> */}
            </Segment>
        )
    }
}

export default Main;