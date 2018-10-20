import React from 'react';
import { Card } from 'semantic-ui-react'

export default props => {
    const { teams, onChange } = props;

    const cardList = teams.map(team => (
        <Card
            key={team.index}
            image='https://upload.wikimedia.org/wikipedia/it/0/07/Fc_barcelona.png'
            header={team.name}
            meta='Team'
            description={team.name + " is amazing!"}
            onClick={() => onChange && onChange(team)}
        />
    ));

    return (
        <Card.Group>
            {cardList}
        </Card.Group>
    );
}