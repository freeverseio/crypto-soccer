import React from 'react';
import { Link, withRouter } from 'react-router-dom';
import { Menu, Icon, Button } from 'semantic-ui-react'

const TopBar = ({account, location}) => {
    const pathname = location.pathname;

    return (
        <Menu pointing secondary>
            <Menu.Item as={Link} to="/" name='home' active={pathname === '/'} />
            <Menu.Item as={Link} to="/academy" name='Academy' active={pathname === '/academy'} />
            {/* <Menu.Item as={Link} to="/shop" name='Shop' active={pathname === '/shop'} /> */}
            <Menu.Item as={Link} to="/teams" name='Teams' active={pathname === '/teams'} />
            <Menu.Item position='right' as={Link} to='/settings' active={pathname === '/settings'}>
                <Icon name='settings' />
            </Menu.Item>
            {window.ethereum &&
                <Menu.Item>
                    <Button icon
                        color={account ? 'green' : 'grey'}
                        onClick={window.ethereum.enable}
                        disabled={account ? true : false}>
                        <Icon name='ethereum' />
                    </Button>
                </Menu.Item>
            }
        </Menu>
    )
}

export default withRouter(TopBar);