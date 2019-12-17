import React from 'react';
import { Switch, Route } from 'react-router-dom';

import Players from './Players';

const Main = (props) => (
    <main>
        <Switch>
            <Route exact path='/' render={() => <Players {...props}/>} />
            {/* <Route exact path='/specialplayers' render={() => <SpecialPlayer {...props} />} /> */}
        </Switch>
    </main>
)

export default Main;