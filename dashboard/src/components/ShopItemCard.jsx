import React, { useState, useEffect } from 'react';
import { Card, Image, Grid, Divider, Button, Input, List } from 'semantic-ui-react';
import gql from 'graphql-tag';
import { useMutation } from '@apollo/react-hooks';
import signPutAssetForSaleMTx from './marketUtils';
import uuidv1 from 'uuid/v1';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import {
    faRunning,
    faBolt,
    faBurn,
    faHeart,
    faShoePrints,
    faShieldAlt,
    faClock,
    faGavel,
} from '@fortawesome/free-solid-svg-icons'

const DELETE_PLAYER = gql`
mutation DeleteAcademyPlayer(
    $playerId: String!
    ) {
        deleteSpecialPlayer(
            playerId: $playerId
        )
    }
`;

const CREATE_AUCTION = gql`
mutation CreateAuction(
    $uuid: UUID!
    $playerId: String!
    $currencyId: Int!
    $price: Int!
    $rnd: Int!
    $validUntil: String!
    $signature: String!
    $seller: String!
    ) {
        createAuction(
            input: {
            uuid: $uuid
            playerId: $playerId
            currencyId: $currencyId
            price: $price
            rnd: $rnd
            validUntil: $validUntil
            signature: $signature
            seller: $seller
            }
        )
    }
`;

const ItemTypeImage = (type) => {
    if type === 0 {
        return <Image src='/kick.png' wrapped ui={false} />

    }
}

export default function PlayerCard(props) {
    const { item} = props;
    const [price, setPrice] = useState(50);
    const [timeout, setTimeout] = useState(120);

    
    // const deletePlayer = async () => {
    //     deletePlayerMutation({
    //         variables: {
    //             playerId: player.playerId,
    //         }
    //     });
    // };
    console.log(item)

    return (
        <Card color='red'>
            <Image src='player.jpg' wrapped ui={false} />
            <Card.Content>
                {/* <Card.Header>{player.name}</Card.Header> */}
                <Divider />
                <Card.Meta>
                    <Grid columns='equal'>
                        {/* <Grid.Row>
                            <Grid.Column textAlign="center"><FontAwesomeIcon icon={faBolt} /> {player.potential}</Grid.Column>
                            <Grid.Column textAlign="center"><FontAwesomeIcon icon={faBurn} /> {player.shoot}</Grid.Column>
                            <Grid.Column textAlign="center"><FontAwesomeIcon icon={faHeart} /> {player.endurance}</Grid.Column>
                        </Grid.Row>
                        <Grid.Row>
                            <Grid.Column textAlign="center"><FontAwesomeIcon icon={faRunning} /> {player.speed}</Grid.Column>
                            <Grid.Column textAlign="center"><FontAwesomeIcon icon={faShoePrints} /> {player.pass}</Grid.Column>
                            <Grid.Column textAlign="center"><FontAwesomeIcon icon={faShieldAlt} /> {player.defence}</Grid.Column>
                        </Grid.Row> */}
                    </Grid>
                </Card.Meta>
                <Card.Description>
                </Card.Description>
            </Card.Content>
            {/* <Card.Content extra>
                {currentAuction &&
                    <Grid columns='equal'>
                        <Grid.Row>
                            <Grid.Column textAlign="center"><FontAwesomeIcon icon={faClock} /> {timeLeft}</Grid.Column>
                            <Grid.Column textAlign="center"><FontAwesomeIcon icon={faGavel} /> {bidsCount}</Grid.Column>
                        </Grid.Row>
                    </Grid>
                }
                {!currentAuction &&
                    <List>
                        <List.Item>
                            <Input
                                label={{ basic: true, content: 'sec' }}
                                type="number"
                                labelPosition='right'
                                fluid
                                icon='clock'
                                iconPosition='left'
                                value={timeout}
                                onChange={event => setTimeout(event.target.value)}
                            />
                        </List.Item>
                        <List.Item>
                            <Input
                                label={{ basic: true, content: '$' }}
                                type="number"
                                labelPosition='right'
                                fluid
                                icon='money'
                                iconPosition='left'
                                value={price}
                                onChange={event => setPrice(event.target.value)}
                            />
                        </List.Item>
                    </List>
                }
                {!currentAuction && <Button floated='right' basic color='green' onClick={createAuction}>Sell</Button>}
                {!currentAuction && <Button floated='right' value='Delete' basic color='red' onClick={deletePlayer}>Delete</Button>}
            </Card.Content> */}
        </Card>
    )
};

function useTimeLeft(currentAuction) {
    const calculateTimeLeft = (currentAuction) => {
        if (!currentAuction) return "";

        const difference = +new Date(currentAuction.validUntil * 1000) - +new Date();
        let timeLeft = "";
        if (difference > 0) {
            const days = Math.floor(difference / (1000 * 60 * 60 * 24));
            const hours = Math.floor((difference / (1000 * 60 * 60)) % 24);
            const minutes = Math.floor((difference / 1000 / 60) % 60);
            const seconds = Math.floor((difference / 1000) % 60);

            if (days > 0) { timeLeft += days + "d"; }
            if (hours > 0) { timeLeft += " " + hours + "h"; }
            if (minutes > 0) { timeLeft += " " + minutes + "m"; }
            if (seconds > 0) { timeLeft += " " + seconds + "s"; }

        }
        return timeLeft;
    }

    const [timeLeft, setTimeLeft] = useState(calculateTimeLeft(currentAuction));
    useEffect(() => {
        const timerID = setInterval(() => {
            setTimeLeft(calculateTimeLeft(currentAuction));
        }, 1000);
        return () => {
            clearInterval(timerID);
        }
    }, [currentAuction]);
    return timeLeft;
}