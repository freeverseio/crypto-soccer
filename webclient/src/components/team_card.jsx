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
        const { ethLeagueManager, index, onClick} = this.props;
        const { name } = this.state;

        if (ethLeagueManager){
            ethLeagueManager.teamName(index)
            .then(result => this.setState({name: result}));
        }

        return (
            <Card
                onClick={() => onClick(index)}
                image='https://static.independent.co.uk/s3fs-public/thumbnails/image/2017/11/21/13/borat.jpg'
                header={name}
                meta='Team'
                description={name + " is amazing!"}
            />
        )
    }
}

export default TeamCard;