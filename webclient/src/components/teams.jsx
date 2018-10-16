import React, { Component } from 'react';
import { Card } from 'semantic-ui-react'
import TeamCard from './team_card';

class Teams extends Component {
    render() {
        return (
            <Card.Group>
                <TeamCard />
            </Card.Group>
        );
    }
}

export default Teams;
