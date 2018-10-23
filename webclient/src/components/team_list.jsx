import React from 'react';
import { Item, Accordion, Icon } from 'semantic-ui-react'
import TeamPlayerTable from '../components/team_players_table';

export default props => {
    const { teams } = props;

    const cardList = teams.map(team => (
        <Item key={team.index}>
            <Item.Image size='tiny' src='https://upload.wikimedia.org/wikipedia/it/0/07/Fc_barcelona.png' />
            <Item.Content>
                <Item.Header as='a'>{team.name}</Item.Header>
                <Item.Meta>TODO Description</Item.Meta>
                <Item.Description>
                    <Accordion>
                        <Accordion.Title active={true} index={0} >
          <Icon name='dropdown' />
          What is a dog?
        </Accordion.Title>
                        <Accordion.Content>
                            <TeamPlayerTable team={team} />
                        </Accordion.Content>
                    </Accordion>
                </Item.Description>
            </Item.Content>
        </Item>
    ));

    return (
        <Item.Group>
             {cardList}
        </Item.Group>
    );
}