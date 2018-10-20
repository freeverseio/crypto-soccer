import React, { Component } from 'react';
import { Segment } from 'semantic-ui-react'

class Header extends Component {
    render() {
        const { testingFacade } = this.props;

        return (
            <Segment clearing>
                    {testingFacade ? "connected" : "disconnected "}
            </Segment>
        )
    }
}

export default Header;