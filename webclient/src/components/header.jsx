import React, { Component } from 'react';
import { Menu, Segment } from 'semantic-ui-react'

class Header extends Component {
    state = { activeItem: 'home' };

    handleItemClick = (e, { name }) => this.setState({ activeItem: name });

    render() {
        const { activeItem } = this.state;

        return (
            <div>
                <Menu pointing secondary>
                    <Menu.Item name='home' active={activeItem === 'home'} onClick={this.handleItemClick} />
                    <Menu.Item name='play' active={activeItem === 'play'} onClick={this.handleItemClick} />
                    <Menu.Item name='teams' active={activeItem === 'teams'} onClick={this.handleItemClick} />
                    <Menu.Item name='market' active={activeItem === 'market'} onClick={this.handleItemClick} />
                    <Menu.Item name='shop' active={activeItem === 'shop'} onClick={this.handleItemClick} />
                </Menu>
            </div>
        )
    }
}

export default Header;