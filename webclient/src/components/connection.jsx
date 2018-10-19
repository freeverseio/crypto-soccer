import React, { Component } from 'react';
import { Header, Segment } from 'semantic-ui-react'

class Connection extends Component {
    render() {
        const { testingFacade } = this.props;

        return (
            <Segment clearing>
                <Header as='h2' floated='right'>
                    {testingFacade ? "connected" : "disconnected "}
                </Header>
            </Segment>
        )
    }
}

export default Connection;