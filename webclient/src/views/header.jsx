import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import { Menu } from 'semantic-ui-react'

class Header extends Component {
    state = { activeItem: 'home' };

    handleItemClick = (e, { name }) => this.setState({ activeItem: name });

    render() {
        const { activeItem } = this.state;

        return (
            <div>
                <Menu pointing secondary>
                    <Link to='/'><Menu.Item name='home' active={activeItem === 'home'} onClick={this.handleItemClick} /></Link>
                    <Link to='/play'><Menu.Item name='play' active={activeItem === 'play'} onClick={this.handleItemClick} /></Link>
                    <Link to='/teams'><Menu.Item name='teams' active={activeItem === 'teams'} onClick={this.handleItemClick} /></Link>
                    <Link to='/market'><Menu.Item name='market' active={activeItem === 'market'} onClick={this.handleItemClick} /></Link>
                    <Link to='/shop'><Menu.Item name='shop' active={activeItem === 'shop'} onClick={this.handleItemClick} /></Link>
                </Menu>
            </div>
        )
    }
}

export default Header;