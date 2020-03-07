import React, { Component } from 'react';
import { Link, withRouter } from 'react-router-dom';
import { Menu, Input } from 'semantic-ui-react'

class TopBar extends Component {
    render() {
        const { url, onUrlChange } = this.props;
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
                        <Input
                            icon='database'
                            iconPosition='left'
                            placeholder='Horizon URL'
                            value={url}
                            onChange={event => onUrlChange(event.target.value)}
                        />
                    </Menu.Item>
                </Menu.Menu>
            </Menu>
        )
    }
}

export default withRouter(TopBar);