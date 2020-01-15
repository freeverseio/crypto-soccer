import React from 'react';
import { Table, Label } from 'semantic-ui-react';
import gql from 'graphql-tag';
import { useQuery } from '@apollo/react-hooks';

const ALL_AUCTIONS = gql`
{
  allAuctions(
    condition: {
      seller: "0xb8CE9ab6943e0eCED004cDe8e3bBed6568B2Fa01"
      state: "PAID"
    }
  ) {
    nodes {
      uuid
      paymentUrl
      price
      bidsByAuction(orderBy: EXTRA_PRICE_DESC, first: 1) {
        nodes {
          extraPrice
          paymentId
        }
      }
      playerByPlayerId {
        name
      }
    }
  }
}`;

const PlayerTableRow = (props) => {
    const { auction } = props;
    const uuid = auction.uuid;
    const amount = auction.price + auction.bidsByAuction.nodes[0].extraPrice;
    const name = auction.playerByPlayerId.name;

    return(
         <Table.Row>
             <Table.Cell>{uuid}</Table.Cell>
             <Table.Cell>{name}</Table.Cell>
             <Table.Cell>{amount}</Table.Cell>
            <Table.Cell>
                <Label as='a' href={auction.paymentUrl} content='withdraw' icon='euro' />
            </Table.Cell>
         </Table.Row>
    )
}

export default function WithdrawalTable(props) {
    const { loading, error, data } = useQuery(ALL_AUCTIONS, {
        pollInterval: 5000,
    });

    if (loading) return null;
    if (error) return `Error! ${error}`;

    const auctions = data.allAuctions.nodes;

    return (
        <Table color='grey' inverted >
            <Table.Header>
                <Table.Row>
                    <Table.HeaderCell>ID</Table.HeaderCell>
                    <Table.HeaderCell>name</Table.HeaderCell>
                    <Table.HeaderCell>Price</Table.HeaderCell>
                </Table.Row>
            </Table.Header>
            <Table.Body >
                { auctions.map(auction =>  <PlayerTableRow key={auction.uuid} auction={auction} /> ) }
            </Table.Body>
        </Table>
    );
}
