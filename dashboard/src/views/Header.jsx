import React, { Component } from 'react';
import { Link, withRouter } from 'react-router-dom';
import { Menu, Icon, Header } from 'semantic-ui-react'

class TopBar extends Component {
    render() {
        const { url } = this.props;
        const location = this.props.location.pathname;

        return (
            <Menu pointing secondary>
                <Link to='/'><Menu.Item name='academy' active={location === '/'} /></Link>

                <Menu.Menu position='right'>
                    <Menu.Item>
                        <Icon name='database' />
                        <Header as='h5' floated='right'>
                            <Header.Content>{url}</Header.Content>
                        </Header>
                    </Menu.Item>

                </Menu.Menu>
            </Menu>
        )
    }
}

export default withRouter(TopBar);