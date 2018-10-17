import React, { Component } from 'react';
import { Header, Segment } from 'semantic-ui-react'

class Connection extends Component {
    render() {
        const { provider, ethLeagueManager } = this.props;

        return (
            <Segment clearing>
                <Header as='h2' floated='right'>
                    {ethLeagueManager ? "connected" : "disconnected "}
                </Header>
                <Header as='h2' floated='left'>
                    {provider.host}
                </Header>
            </Segment>
        )
    }
}

export default Connection;