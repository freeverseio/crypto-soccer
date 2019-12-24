import React from 'react';
import { Switch, Route } from 'react-router-dom';

import Academy from './academy/Academy';

const Main = (props) => (
    <main>
        <Switch>
            <Route exact path='/' render={() => <Academy {...props}/>} />
            {/* <Route exact path='/specialplayers' render={() => <SpecialPlayer {...props} />} /> */}
        </Switch>
    </main>
)

export default Main;