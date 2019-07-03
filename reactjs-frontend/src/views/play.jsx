import React, { PureComponent } from 'react';
import { Segment } from 'semantic-ui-react'
import Match from '../components/match';

class Play extends PureComponent {
    constructor(props){
        super(props);

        this.state = {};
    }

    render() {
        const { testingFacade, teams } = this.props;

        return (
            <Segment>
                <Match contract={testingFacade} teams={teams} />
            </Segment>
        )
    }
}

export default Play;