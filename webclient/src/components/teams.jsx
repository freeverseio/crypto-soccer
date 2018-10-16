import React, { PureComponent } from 'react';
import { Card, Divider } from 'semantic-ui-react'
import TeamCard from './team_card';
import TeamPlayerTable from './team_players_table';

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
                    onClick={index => this.setState({ selectedIndex: index })}
                />
            );
        }

        return (
            <React.Fragment>
                <Card.Group>
                    {teams}
                </Card.Group>
                <Divider/>
                <TeamPlayerTable index={selectedIndex}/>
            </React.Fragment>
        );
    }
}

export default Teams;
