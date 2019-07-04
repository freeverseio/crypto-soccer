import React from 'react';
import { Item } from 'semantic-ui-react'

export const AttackEvent = (props) => (
    <Item>
        <Item.Image size='tiny' src='svg/attack.svg' />
        <Item.Content verticalAlign='middle'>
            <Item.Header>{props.min && props.min}</Item.Header>
            <Item.Description>{props.text && props.text}</Item.Description>
        </Item.Content>
    </Item>
);

export const DefendEvent = (props) => (
    <Item>
        <Item.Image size='tiny' src='svg/defend.svg' />
        <Item.Content verticalAlign='middle'>
            <Item.Description>{props.text && props.text}</Item.Description>
        </Item.Content>
    </Item>
);

export const ShootEvent = (props) => (
    <Item>
        <Item.Image size='tiny' src='svg/shoot.svg' />
        <Item.Content verticalAlign='middle'>
            <Item.Description>{props.text && props.text}</Item.Description>
        </Item.Content>
    </Item>
);

export const GoalEvent = (props) => (
    <Item>
        <Item.Image size='tiny' src='svg/goal.svg' />
        <Item.Content verticalAlign='middle'>
            <Item.Header>{props.text && props.text}</Item.Header>
        </Item.Content>
    </Item>
);

export const BlockedEvent = (props) => (
    <Item>
        <Item.Image size='tiny' src='svg/catch.svg' />
        <Item.Content verticalAlign='middle'>
            <Item.Description>{props.text && props.text}</Item.Description>
        </Item.Content>
    </Item>
);