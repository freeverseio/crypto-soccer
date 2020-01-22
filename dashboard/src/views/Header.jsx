import React, { Component } from 'react';
import { Link, withRouter } from 'react-router-dom';
import { Menu, Icon, Header } from 'semantic-ui-react'

class TopBar extends Component {
    render() {
        const { url } = this.props;
        const location = this.props.location.pathname;

        return (
            <Menu pointing >
                <Menu.Item as={Link} to="/" active={location === '/'} >
                    <img src='/logo62.png' alt="not found"/>
                </Menu.Item>
                <Menu.Item as={Link} to="/academy" name='Academy' active={location === '/academy'} />
                <Menu.Item as={Link} to="/shop" name='Shop' active={location === '/shop'} />
                <Menu.Item as={Link} to="/teams" name='Teams' active={location === '/teams'} />
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