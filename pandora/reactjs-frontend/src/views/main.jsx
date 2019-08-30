import React from 'react';
import { Switch, Route } from 'react-router-dom';

import Home from './home';
import Play from './play';
import Teams from './teams';
import Market from './market';
import Shop from './shop';

const Main = (props) => (
    <main>
        <Switch>
            <Route exact path='/' render={() => <Home {...props}/>} />
            <Route exact path='/play' render={() => <Play {...props} />} />
            <Route exact path='/teams' render={() => <Teams {...props} />} />
            <Route exact path='/market' render={() => <Market {...props} />} />
            <Route exact path='/shop' render={() => <Shop {...props} />} />
        </Switch>
    </main>
)

export default Main;