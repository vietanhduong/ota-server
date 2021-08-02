import React from 'react';
import { Redirect, Route, Switch, useHistory } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { privateRoute, publicRoute } from 'routes';

const PublicLayout = () => {
  const history = useHistory();
  const { isLoggedIn } = useSelector(({ profile }) => profile);

  React.useEffect(() => {
    if (isLoggedIn) {
      history.replace(privateRoute.home.path);
    }
  }, [history, isLoggedIn]);

  return (
    <div className='App Auth-Layout'>
      <Switch>
        {Object.values(publicRoute).map(({ path, component }) => (
          <Route exact key={path} path={path} component={component} />
        ))}
        <Redirect from='/' to={publicRoute.signIn.path} />
      </Switch>
    </div>
  );
};

export default PublicLayout;
