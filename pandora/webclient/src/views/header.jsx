import React, { Component } from 'react';
import { Link, withRouter } from 'react-router-dom';
import { Menu, Icon, Header } from 'semantic-ui-react'

class TopBar extends Component {
    render() {
        const { url } = this.props;
        const location = this.props.location.pathname;

        return (
            <Menu pointing secondary>
                <Link to='/'><Menu.Item name='home' active={location === '/'} /></Link>
                <Link to='/play'><Menu.Item name='play' active={location === '/play'} /></Link>
                <Link to='/teams'><Menu.Item name='teams' active={location === '/teams'} /></Link>
                <Link to='/market'><Menu.Item name='market' active={location === '/market'} /></Link>
                <Link to='/shop'><Menu.Item name='shop' active={location === '/shop'} /></Link>

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