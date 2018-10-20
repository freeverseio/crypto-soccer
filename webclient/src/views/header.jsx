import React, { Component } from 'react';
import { Link, withRouter } from 'react-router-dom';
import { Menu } from 'semantic-ui-react'

class Header extends Component {
    render() {
        const location = this.props.location.pathname;

        return (
            <div>
                <Menu pointing secondary>
                    <Link to='/'><Menu.Item name='home' active={location === '/'} /></Link>
                    <Link to='/play'><Menu.Item name='play' active={location === '/play'} /></Link>
                    <Link to='/teams'><Menu.Item name='teams' active={location === '/teams'} /></Link>
                    <Link to='/market'><Menu.Item name='market' active={location === '/market'} /></Link>
                    <Link to='/shop'><Menu.Item name='shop' active={location === '/shop'} /></Link>
                </Menu>
            </div>
        )
    }
}

export default withRouter(Header);