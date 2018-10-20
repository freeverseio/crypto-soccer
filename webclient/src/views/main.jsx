import React from 'react';
import { Switch, Route } from 'react-router-dom';

import Play from './play';

const Main = (props) => (
    <main>
        <Switch>
            <Route exact path='/' render={() => <Play {...props} />} />
            <Route exact path='/play' render={() => <p>Play</p>} />
        </Switch>
    </main>
)

export default Main;