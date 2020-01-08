import React from 'react';
import { Card, Image, Grid, Divider, Button } from 'semantic-ui-react';
import gql from 'graphql-tag';
import { useMutation } from '@apollo/react-hooks';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import {
    faShoePrints,
    faShieldAlt,
    faDollarSign,
} from '@fortawesome/free-solid-svg-icons';

const DELETE_SHOP_ITEM = gql`
mutation DeleteShopItem(
    $uuid: UUID!
    ) {
        deleteShopItem(uuid: $uuid)
    }
`;

export default function PlayerCard(props) {
    const { item, options } = props;
    const [deleteShopItemMutation] = useMutation(DELETE_SHOP_ITEM);

    const deleteShopItem = async () => {
        deleteShopItemMutation({
            variables: {
                uuid: item.uuid,
            }
        });
    };

    return (
        <Card color='red'>
            <Image src={options.image} wrapped ui={false} />
            <Card.Content>
                <Card.Header>{options.text}</Card.Header>
                <Divider />
                <Card.Meta>
                    <Grid columns='equal'>
                        <Grid.Row>
                            <Grid.Column textAlign="center"><FontAwesomeIcon icon={faDollarSign} /> {item.price}</Grid.Column>
                        </Grid.Row>
                        {/* <Grid.Row>
                            <Grid.Column textAlign="center"><FontAwesomeIcon icon={faShoePrints} /> {item.price}</Grid.Column>
                            <Grid.Column textAlign="center"><FontAwesomeIcon icon={faShieldAlt} /> {item.type}</Grid.Column>
                        </Grid.Row> */}
                    </Grid>
                </Card.Meta>
                <Card.Description>
                </Card.Description>
            </Card.Content>
            <Card.Content extra>
                <Button value='Delete' color='red' fluid onClick={deleteShopItem}>Delete</Button>
            </Card.Content>
        </Card>
    )
};
