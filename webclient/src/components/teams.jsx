import React, { PureComponent } from 'react';
import { Card } from 'semantic-ui-react'
import TeamCard from './team_card';

class Teams extends PureComponent {
    constructor(props){
        super(props);

        this.state = {
            count: 0,
            selectedIndex: 0
        }
    }

    render() {
        const { ethLeagueManager } = this.props;
        const { count, selectedIndex } = this.state;

        if (ethLeagueManager) {
            ethLeagueManager.countTeams()
                .then(result => this.setState({ count: result }));
        }

        let teams = [];
        for (var i = 0; i < count; i++) {
            teams.push(
                <TeamCard
                    key={i}
                    index={i}
                    ethLeagueManager={ethLeagueManager}
                    onClick={index => this.setState({selectedIndex: index})}
                />
            );
        }

        return (
            <Card.Group>
                {teams}
            </Card.Group>
        );
    }
}

export default Teams;
