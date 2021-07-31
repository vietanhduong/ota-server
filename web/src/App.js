import React from 'react';
import {BrowserRouter as Router, Route, Switch} from 'react-router-dom';
import {PrivateLayout, PublicLayout} from 'layouts';
import {Provider} from 'react-redux';
import {store} from 'reducers';
import {MuiThemeProvider} from '@material-ui/core';
import {appTheme} from 'containers/Theme';
import {profileAction} from 'actions/profile';
import './App.scss';
import {jwt_decode} from 'utils/common';

const App = () => {
  const [isReady, setIsReady] = React.useState(false);

  React.useEffect(() => {
    const profile = JSON.parse(localStorage.getItem('profile'));
    if (!profile) {
      setIsReady(true);
      return
    }

    const payload = jwt_decode(profile.access_token);
    const timeLeft = payload.exp * 1000 - Date.now();

    if (timeLeft > 60 * 5 * 1000) {
      profileAction.login(profile);
      setTimeout(() => {
        profileAction.refresh(profile);
      }, timeLeft - 60 * 5 * 1000);
      setIsReady(true);
    } else {
      profileAction.refresh(profile).finally(() => {
        setIsReady(true);
      });
    }
  }, []);

  return (
    <Provider store={store}>
      <MuiThemeProvider theme={appTheme}>
        <Router>
          {isReady && (
            <Switch>
              <Route path='/auth' component={PublicLayout}/>
              <Route path='/' component={PrivateLayout}/>
            </Switch>
          )}
        </Router>
      </MuiThemeProvider>
    </Provider>
  );
};

export default App;
