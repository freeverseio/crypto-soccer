import React, { PureComponent } from 'react';
import { Button, Card, Image } from 'semantic-ui-react'

class TeamCard extends PureComponent {
    constructor(props){
        super(props);

        this.state = {
            name: ''
        }
    }

    render() {
        const { ethLeagueManager, index } = this.props;
        const { name } = this.state;

        if (ethLeagueManager){
            ethLeagueManager.teamName(index)
            .then(result => this.setState({name: result}));
        }

        return (
            <Card>
                <Card.Content>
                    <Image floated='right' size='mini' src='https://react.semantic-ui.com/images/avatar/large/steve.jpg' />
                    <Card.Header>{name}</Card.Header>
                    <Card.Meta>Amazing Team!</Card.Meta>
                    <Card.Description>
                        {name} is will win!
                    </Card.Description>
                </Card.Content>
                <Card.Content extra>
                    <Button basic color='red'>
                        Delete
                </Button>
                </Card.Content>
            </Card>
        )
    }
}

export default TeamCard;