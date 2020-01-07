import React, { Component } from 'react';
import { Link, withRouter } from 'react-router-dom';
import { Menu, Icon, Header } from 'semantic-ui-react'

class TopBar extends Component {
    render() {
        const { url } = this.props;
        const location = this.props.location.pathname;

        return (
            <Menu pointing secondary>
                <Link to='/'><Menu.Item name='Home' active={location === '/'} /></Link>
                <Link to='/academy'><Menu.Item name='academy' active={location === '/academy'} /></Link>
                <Link to='/shop'><Menu.Item name='shop' active={location === '/shop'} /></Link>
                <Link to='/teams'><Menu.Item name='Teams' active={location === '/teams'} /></Link>
                <Link to='/players'><Menu.Item name='Players' active={location === '/players'} /></Link>

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