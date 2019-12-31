import React from 'react';
import { Switch, Route } from 'react-router-dom';

import Home from './home/Home';
import Teams from './teams/Teams';
import Players from './players/Players';
import Academy from './academy/Academy';

const Main = (props) => (
    <main>
        <Switch>
            <Route exact path='/' render={() => <Home {...props}/>} />
            <Route exact path='/academy' render={() => <Academy {...props}/>} />
            <Route exact path='/teams' render={() => <Teams {...props}/>} />
            <Route exact path='/players' render={() => <Players {...props}/>} />
            {/* <Route exact path='/specialplayers' render={() => <SpecialPlayer {...props} />} /> */}
        </Switch>
    </main>
)

export default Main;