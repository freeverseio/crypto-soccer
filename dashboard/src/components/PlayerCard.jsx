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

export default function PlayerCard(props) {
    const { player, web3 } = props;
    const [price, setPrice] = useState(50);
    const [timeout, setTimeout] = useState(120);
    const [createAuctionMutation] = useMutation(CREATE_AUCTION);
    const [deletePlayerMutation] = useMutation(DELETE_PLAYER);

    const date = new Date();
    const nowSeconds = Math.round(date.getTime() / 1000);
    const lastAuction = player.auctionsByPlayerId.nodes[0];
    const currentAuction = (lastAuction && (lastAuction.validUntil > nowSeconds)) ? lastAuction : null;
    const bidsCount = currentAuction ? (currentAuction.bidsByAuction.totalCount) : 0;
    const timeLeft = useTimeLeft(currentAuction);
    

    const createAuction = async () => {
        const rnd = Math.floor(Math.random() * 1000000);
        const now = new Date();
        const validUntil = (Math.round(now.getTime() / 1000) + Number(timeout)).toString();
        const sellerAccount = await web3.eth.accounts.privateKeyToAccount("0x348ce564d427a3311b6536bbcff9390d69395b06ed6c486954e971d960fe8709");
        const currencyId = 1;
        const signature = await signPutAssetForSaleMTx(web3, currencyId, price, rnd, validUntil, player.playerId, sellerAccount);
        const seller = sellerAccount.address;
        createAuctionMutation({
            variables: {
                uuid: uuidv1(),
                playerId: player.playerId,
                currencyId: currencyId,
                price: Number(price),
                rnd: Number(rnd),
                validUntil: validUntil,
                signature: signature.signature,
                seller: seller,
            }
        });
    };

    const deletePlayer = async () => {
        deletePlayerMutation({
            variables: {
                playerId: player.playerId,
            }
        });
    };

    return (
        <Card>
            <Image src='player.jpg' wrapped ui={false} />
            <Card.Content>
                <Card.Header>{player.name}</Card.Header>
                <Divider />
                <Card.Meta>
                    <Grid columns='equal'>
                        <Grid.Row>
                            <Grid.Column textAlign="center"><FontAwesomeIcon icon={faBolt} />{player.potential}</Grid.Column>
                            <Grid.Column textAlign="center"><FontAwesomeIcon icon={faBurn} />{player.shoot}</Grid.Column>
                            <Grid.Column textAlign="center"><FontAwesomeIcon icon={faHeart} />{player.endurance}</Grid.Column>
                            <Grid.Column textAlign="center"><FontAwesomeIcon icon={faRunning} />{player.speed}</Grid.Column>
                            <Grid.Column textAlign="center"><FontAwesomeIcon icon={faShoePrints} />{player.pass}</Grid.Column>
                            <Grid.Column textAlign="center"><FontAwesomeIcon icon={faShieldAlt} />{player.defence}</Grid.Column>
                        </Grid.Row>
                    </Grid>
                </Card.Meta>
                <Card.Description>
                </Card.Description>
            </Card.Content>
            <Card.Content extra>
                {currentAuction &&
                    <Grid columns='equal'>
                        <Grid.Row>
                            <Grid.Column textAlign="center"><FontAwesomeIcon icon={faClock} /> {timeLeft}</Grid.Column>
                            <Grid.Column textAlign="center"><FontAwesomeIcon icon={faGavel} /> {bidsCount}</Grid.Column>
                            {/* <Grid.Column textAlign="center"><FontAwesomeIcon icon={faMoneye} /> TODO</Grid.Column> */}
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
            </Card.Content>
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
    });
    return timeLeft;
}