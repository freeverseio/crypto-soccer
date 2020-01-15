import React, {useState, useEffect, Fragment} from 'react';
import { Table, Label, Button, Input, Container } from 'semantic-ui-react';
import gql from 'graphql-tag';
import { useMutation } from '@apollo/react-hooks';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import {
    faClock,
    faGavel,
} from '@fortawesome/free-solid-svg-icons';
import uuidv1 from 'uuid/v1';
import auctionAnalizer from './AuctionAnalizer';
import signPutAssetForSaleMTx from './marketUtils';

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

const PlayerTableRow = (props) => {
    const { player, web3 } = props;
    const [timeout, setTimeout] = useState(120);
    const [price, setPrice] = useState(50);
    const [createAuctionMutation] = useMutation(CREATE_AUCTION);

    const date = new Date();
    const nowSeconds = Math.round(date.getTime() / 1000);
    const lastAuction = player.auctionsByPlayerId.nodes[0];
    const currentAuction = (lastAuction && (lastAuction.validUntil > nowSeconds)) ? lastAuction : null;
    const bidsCount = currentAuction ? (currentAuction.bidsByAuction.totalCount) : 0;
    const timeLeft = useTimeLeft(currentAuction);

    const canBePutOnSale = auctionAnalizer.canBePutOnSale(lastAuction);
    const isWaitingPayment = auctionAnalizer.isWaitingPayment(lastAuction);
    const isWaitingWithdrawal = auctionAnalizer.isWaitingWithdrawal(lastAuction);

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
        })
            .catch(console.error);
    };

    return (
        <Table.Row key={player.playerId}>
            <Table.Cell>{player.name}</Table.Cell>
            <Table.Cell>{player.defence}</Table.Cell>
            <Table.Cell>{player.pass}</Table.Cell>
            <Table.Cell>{player.speed}</Table.Cell>
            <Table.Cell>{player.endurance}</Table.Cell>
            <Table.Cell>{player.shoot}</Table.Cell>
            <Table.Cell>
                {isWaitingPayment && <Label>Paying</Label>}
                {canBePutOnSale &&
                    <React.Fragment>
                        <Input size='mini' type="number" icon='clock' value={timeout}
                            style={{ width: 100 }}
                            onChange={event => setTimeout(event.target.value)}
                        />
                        <Input size='mini' type="number" icon='money' value={price}
                            style={{ width: 100 }}
                            onChange={event => setPrice(event.target.value)}
                        />
                        <Button size='mini' color='green' onClick={createAuction}>Sell</Button>
                    </React.Fragment>
                }
                {currentAuction &&
                    <React.Fragment>
                        <Label>
                            <FontAwesomeIcon icon={faGavel} /> {bidsCount}
                        </Label>
                        <Label>
                            <FontAwesomeIcon icon={faClock} /> {timeLeft}
                        </Label>
                    </React.Fragment>
                }
            </Table.Cell>
        </Table.Row>
    )
}

export default function PlayersTable(props) {
    const players = props.players ? props.players : [];
    const web3 = props.web3;

    return (
        <Table color='grey' inverted >
            <Table.Header>
                <Table.Row>
                    <Table.HeaderCell>Name</Table.HeaderCell>
                    <Table.HeaderCell>Defence</Table.HeaderCell>
                    <Table.HeaderCell>Pass</Table.HeaderCell>
                    <Table.HeaderCell>Speed</Table.HeaderCell>
                    <Table.HeaderCell>Endurance</Table.HeaderCell>
                    <Table.HeaderCell>Shoot</Table.HeaderCell>
                </Table.Row>
            </Table.Header>
            <Table.Body >
                { players.map(player =>  <PlayerTableRow key={player.playerId} player={player} web3={web3} /> ) }
            </Table.Body>
        </Table>
    );
}

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