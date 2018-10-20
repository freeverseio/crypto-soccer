import React, { PureComponent } from 'react';
import { Segment, Divider } from 'semantic-ui-react'
import TeamPlayerTable from '../components/team_players_table';
import TeamCreator from '../components/team_creator';
import TeamList from '../components/team_list';
import Match from '../components/match';

class Main extends PureComponent {
    constructor(props){
        super(props);

        this.state = {};
    }

    render() {
        const { testingFacade, teams } = this.props;
        const { team } = this.state;

        return (
            <Segment>
                <Match contract={testingFacade} teams={teams} />
                <TeamCreator contract={testingFacade}/>
                <Divider />
                <TeamList teams={teams} onChange={team => this.setState({team})} />
                <Divider />
                <TeamPlayerTable team={team} />
            </Segment>
        )
    }
}

export default Main;