import React, {useState} from 'react';
import { Table, Label, Button, Input } from 'semantic-ui-react';
import gql from 'graphql-tag';
import { useMutation } from '@apollo/react-hooks';
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

    const lastAuction = player.auctionsByPlayerId.nodes[0];
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
                            onChange={event => setTimeout(event.target.value)}
                        />
                        <Input size='mini' type="number" icon='money' value={price}
                            onChange={event => setPrice(event.target.value)}
                        />
                        <Button size='mini' color='green' onClick={createAuction}>Sell</Button>
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
        <Table color='grey' inverted selectable >
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