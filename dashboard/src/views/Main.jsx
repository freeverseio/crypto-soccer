import React from 'react';
import { Switch, Route } from 'react-router-dom';

import Home from './home/Home';
import Teams from './teams/Teams';
import Players from './players/Players';
import Academy from './academy/Academy';
import Shop from './shop/Shop';
import Settings from './settings/Settings';
import Iap from './iap/Iap';
import Auctions from './auctions/Auctions';
import Messages from './messages/Messages';

const Main = (props) => (
    <main>
        <Switch>
            <Route exact path='/' render={() => <Home {...props}/>} />
            <Route exact path='/academy' render={() => <Academy {...props}/>} />
            <Route exact path='/shop' render={() => <Shop {...props}/>} />
            <Route exact path='/teams' render={() => <Teams {...props}/>} />
            <Route exact path='/players' render={() => <Players {...props}/>} />
            <Route exact path='/settings' render={() => <Settings {...props}/>} />
            <Route exact path='/iap' render={() => <Iap {...props}/>} />
            <Route exact path='/auctions' render={() => <Auctions {...props}/>} />
            <Route exact path='/messages' render={() => <Messages {...props}/>} />
            {/* <Route exact path='/specialplayers' render={() => <SpecialPlayer {...props} />} /> */}
        </Switch>
    </main>
)

export default Main;