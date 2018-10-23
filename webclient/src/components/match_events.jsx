import React from 'react';
import { Segment, Button, Icon, Grid, Header, GridColumn, Item, Progress, Image } from 'semantic-ui-react'

export const AttackEvent = (props) => (
    <Item>
        <Item.Image size='tiny' src='https://images2.corriereobjects.it/methode_image/2016/05/04/Cultura/Foto%20Cultura%20-%20Trattate/italia-germania-1982_650x435%20(1)-kOeB-U43180371083434wgE-1224x916@Corriere-Web-Sezioni-593x443.jpg?v=20160505000206' />
        <Item.Content verticalAlign='middle'>
            <Item.Header>{props.min && props.min}</Item.Header>
            <Item.Description>{props.text && props.text}</Item.Description>
        </Item.Content>
    </Item>
);

export const DefendEvent = (props) => (
    <Item>
        <Item.Image size='tiny' src='https://images2.corriereobjects.it/methode_image/2016/05/04/Cultura/Foto%20Cultura%20-%20Trattate/italia-germania-1982_650x435%20(1)-kOeB-U43180371083434wgE-1224x916@Corriere-Web-Sezioni-593x443.jpg?v=20160505000206' />
        <Item.Content verticalAlign='middle'>
            <Item.Description>{props.text && props.text}</Item.Description>
        </Item.Content>
    </Item>
);

export const ShootEvent = (props) => (
    <Item>
        <Item.Image size='tiny' src='https://images2.corriereobjects.it/methode_image/2016/05/04/Cultura/Foto%20Cultura%20-%20Trattate/italia-germania-1982_650x435%20(1)-kOeB-U43180371083434wgE-1224x916@Corriere-Web-Sezioni-593x443.jpg?v=20160505000206' />
        <Item.Content verticalAlign='middle'>
            <Item.Description>{props.text && props.text}</Item.Description>
        </Item.Content>
    </Item>
);

export const GoalEvent = (props) => (
    <Item>
        <Item.Image size='tiny' src='https://images2.corriereobjects.it/methode_image/2016/05/04/Cultura/Foto%20Cultura%20-%20Trattate/italia-germania-1982_650x435%20(1)-kOeB-U43180371083434wgE-1224x916@Corriere-Web-Sezioni-593x443.jpg?v=20160505000206' />
        <Item.Content verticalAlign='middle'>
            <Item.Header>{props.text && props.text}</Item.Header>
        </Item.Content>
    </Item>
);

export const BlockedEvent = (props) => (
    <Item>
        <Item.Image size='tiny' src='https://images2.corriereobjects.it/methode_image/2016/05/04/Cultura/Foto%20Cultura%20-%20Trattate/italia-germania-1982_650x435%20(1)-kOeB-U43180371083434wgE-1224x916@Corriere-Web-Sezioni-593x443.jpg?v=20160505000206' />
        <Item.Content verticalAlign='middle'>
            <Item.Description>{props.text && props.text}</Item.Description>
        </Item.Content>
    </Item>
);