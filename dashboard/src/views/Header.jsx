import React from 'react';
import { Link, withRouter } from 'react-router-dom';
import { Menu, Icon, Button } from 'semantic-ui-react'



const TopBar = (props) => {
    const { account, location } = props;

    return (
        <Menu>
            <Menu.Item as={Link} to="/" active={location === '/'} >
                <img src='/logo62.png' alt="not found" />
            </Menu.Item>
            <Menu.Item as={Link} to="/academy" name='Academy' active={location === '/academy'} />
            <Menu.Item as={Link} to="/shop" name='Shop' active={location === '/shop'} />
            <Menu.Item as={Link} to="/teams" name='Teams' active={location === '/teams'} />
            <Menu.Item position='right' as={Link} to='/settings' active={location === '/settings'}>
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