import React, { Component } from 'react';
import { Card } from 'semantic-ui-react'
import TeamCard from './team_card';

class Teams extends Component {
    constructor(props){
        super(props);

        this.state = {
            count: 0
        }
    }

    render() {
        const { ethLeagueManager } = this.props;
        const { count } = this.state;

        if (ethLeagueManager){
            ethLeagueManager.countTeams()
            .then(result => this.setState({count: result}));
        }

        let teams = [];
        for (var i = 0; i < count; i++) {
            teams.push(<TeamCard key={i}/>);
        }

        return (
            <Card.Group>
                {teams}
            </Card.Group>
        );
    }
}

export default Teams;
