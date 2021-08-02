import React from 'react';
import { Redirect, Route, Switch, useHistory } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { privateRoute, publicRoute } from 'routes';

const PrivateLayout = () => {
  const history = useHistory();
  const { isLoggedIn } = useSelector(({ profile }) => profile);

  React.useEffect(() => {
    if (!isLoggedIn) {
      history.replace(publicRoute.signIn.path);
    }
  }, [history, isLoggedIn]);

  return (
    <div className='App Private-Layout'>
      <Switch>
        {Object.values(privateRoute).map(({ path, component }) => (
          <Route exact key={path} path={path} component={component} />
        ))}
        <Redirect from='/' to={privateRoute.home.path} />
      </Switch>
    </div>
  );
};

export default PrivateLayout;
