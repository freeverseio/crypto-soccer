import React, { useState } from 'react';
import { Input, Item, Button, List } from 'semantic-ui-react';
import gql from 'graphql-tag';
import { useMutation } from '@apollo/react-hooks';
import signPutAssetForSaleMTx from './marketUtils';
const uuidv1 = require('uuid/v1');

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

export default function AcademyPlayer(props) {
    const [price, setPrice] = useState(50);
    const [timeout, setTimeout] = useState(3600);
    const [createAuction] = useMutation(CREATE_AUCTION);
    const [deleteAcademyPlayer] = useMutation(DELETE_PLAYER);

    const { player, web3 } = props;
    console.log(web3)

    return (
        <Item>
            <Item.Content>
                <Item.Header>{player.name}</Item.Header>
                <Item.Meta>id: {player.playerId}</Item.Meta>
                <Item.Description>
                    <List>
                        <List.Item>
                            <List.Icon name='users' />
                            <List.Content>{player.shoot}</List.Content>
                        </List.Item>
                    </List>
                </Item.Description>
                <Item.Extra>
                    <Input label='Price' type='number' value={price} onChange={event => setPrice(event.target.value)} />
                    <Input label='Timeout' type='number' value={timeout} onChange={event => setTimeout(event.target.value)} />

                    <Button floated='right' basic color='green' onClick={async () => {
                        const rnd = Math.floor(Math.random() * 1000000);
                        const now = new Date();
                        const validUntil = (Math.round(now.getTime() / 1000) + timeout).toString();
                        const sellerAccount = await web3.eth.accounts.privateKeyToAccount("0x348ce564d427a3311b6536bbcff9390d69395b06ed6c486954e971d960fe8709");
                        const currencyId = 1;
                        const signature = await signPutAssetForSaleMTx(web3, currencyId, price, rnd, validUntil, player.playerId, sellerAccount);
                        const seller = sellerAccount.address;
                        createAuction({
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
                    }}>
                        Sell
                                        </Button>
                    <Button floated='right' value='Delete' basic color='red' onClick={() => {
                        deleteAcademyPlayer({
                            variables: {
                                playerId: player.playerId,
                            }
                        })
                    }
                    }>Delete</Button>
                </Item.Extra>
            </Item.Content>
        </Item>
    )
};
