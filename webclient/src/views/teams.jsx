import React, { PureComponent } from 'react';
import { Segment, Divider } from 'semantic-ui-react'
import TeamCreator from '../components/team_creator';
import TeamList from '../components/team_list';

class Teams extends PureComponent {
    render() {
        const { testingFacade, teams } = this.props;

        return (
            <Segment>
                <TeamCreator contract={testingFacade}/>
                <Divider />
                <TeamList teams={teams} />
            </Segment>
        )
    }
}

export default Teams;