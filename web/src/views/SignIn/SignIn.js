import React from 'react';
import { Link } from 'react-router-dom';
import { Loading } from 'components';
import { Button, Paper, TextField, Typography } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import { profileAction } from 'actions/profile';
import { userService } from 'services/user';
import { sha256 } from 'utils/common';
import { publicRoute } from 'routes';

const SignInForm = () => {
  const classes = useStyles();

  const [username, setUsername] = React.useState('vietanhs0817@gmail.com');
  const [usernameError, setUsernameError] = React.useState('');
  const [password, setPassword] = React.useState('admin');
  const [passwordError, setPasswordError] = React.useState('');

  const [isLoading, setIsLoading] = React.useState(false);

  const validateUsername = (value) => {
    if (value.trim() === '') {
      return 'Email cannot be empty';
    }
    return '';
  };

  const handleChangeUsername = (event) => {
    const { value } = event.target;
    setUsername(value);
    setUsernameError(validateUsername(value));
  };

  const validatePassword = (value) => {
    if (value.trim() === '') {
      return 'Password cannot be empty';
    }
    return '';
  };

  const handleChangePassword = (event) => {
    const { value } = event.target;
    setPassword(value);
    setPasswordError(validatePassword(value));
  };

  const handleClickSubmit = () => {
    const usernameError = validateUsername(username);
    setUsernameError(usernameError);
    const passwordError = validatePassword(password);
    setPasswordError(passwordError);
    if (usernameError || passwordError) return;

    setIsLoading(true);
    const body = { email: username, password: sha256(password) };
    userService
      .login(body)
      .then((profile) => {
        profileAction.login(profile);
        setTimeout(() => {
          profileAction.refresh(profile);
        }, 60 * 55 * 1000);
      })
      .finally(() => {
        setIsLoading(false);
      });
  };

  const handlePressKey = (event) => {
    if (event.key === 'Enter') handleClickSubmit();
  };

  return (
    <Paper className={classes.paper}>
      <Typography variant='h5' className={classes.header}>
        {'Sign In'}
      </Typography>
      <TextField
        variant='outlined'
        label={'Email'}
        className={classes.input}
        value={username}
        error={Boolean(usernameError)}
        onChange={handleChangeUsername}
        onKeyPress={handlePressKey}
      />
      <TextField
        variant='outlined'
        type='password'
        label={'Password'}
        className={classes.input}
        value={password}
        error={Boolean(passwordError)}
        onChange={handleChangePassword}
        onKeyPress={handlePressKey}
      />
      <Button
        variant='contained'
        color='primary'
        className={classes.button}
        startIcon={<Loading visible={isLoading} />}
        onClick={handleClickSubmit}
      >
        {'Sign In'}
      </Button>

      <Typography className={classes.link}>
        {`Don't have account?`} <Link to={publicRoute.signUp.path}>{'Sign Up'}</Link>
      </Typography>
    </Paper>
  );
};

const useStyles = makeStyles((theme) => ({
  header: {
    alignSelf: 'center',
    marginBottom: theme.spacing(2),
  },
  paper: {
    padding: 24,
    width: 420,
    display: 'flex',
    flexDirection: 'column',
  },
  input: {
    marginTop: theme.spacing(2),
    marginBottom: theme.spacing(2),
  },
  button: {
    marginTop: theme.spacing(2),
    marginBottom: theme.spacing(2),
  },
  link: {
    marginTop: theme.spacing(2),
    marginBottom: theme.spacing(2),
  },
}));

export default SignInForm;
