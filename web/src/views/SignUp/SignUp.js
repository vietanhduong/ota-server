import React from 'react';
import { Link } from 'react-router-dom';
import { Loading } from 'components';
import { Button, Paper, TextField, Typography } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
// import { profileAction } from 'actions/profile';
import { userService } from 'services/user';
// import { validator } from 'utils/validator';
import { publicRoute } from 'routes';

const SignUpForm = () => {
  const classes = useStyles();

  const [username, setUsername] = React.useState('');
  const [usernameError, setUsernameError] = React.useState('');
  const [password, setPassword] = React.useState('');
  const [passwordError, setPasswordError] = React.useState('');

  const [email, setEmail] = React.useState('');
  const [emailError, setEmailError] = React.useState('');
  const [phone, setPhone] = React.useState('');
  const [phoneError, setPhoneError] = React.useState('');

  const [isLoading, setIsLoading] = React.useState(false);

  const validateUsername = (value) => {
    if (value.trim() === '') {
      return 'Username cannot be empty';
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

  const validateEmail = (value) => {
    if (value.trim() === '') {
      return 'Email cannot be empty';
    }
    return '';
  };

  const handleChangeEmail = (event) => {
    const { value } = event.target;
    setEmail(value);
    setEmailError(validateEmail(value));
  };

  const validatePhone = (value) => {
    if (value.trim() === '') {
      return '';
    }
    return '';
  };

  const handleChangePhone = (event) => {
    const { value } = event.target;
    setPhone(value);
    setPhoneError(validatePhone(value));
  };

  const handleClickSubmit = () => {
    const usernameError = validateUsername(username);
    setUsernameError(usernameError);
    const passwordError = validatePassword(password);
    setPasswordError(passwordError);
    const emailError = validateEmail(email);
    setEmailError(emailError);
    const phoneError = validatePhone(phone);
    setPhoneError(phoneError);
    if (usernameError || passwordError || emailError || phoneError) return;

    setIsLoading(true);
    const body = {
      username,
      password,
      email,
      phone,
    };
    userService.register(body).finally(() => {
      setIsLoading(false);
    });
  };

  const handlePressKey = (event) => {
    if (event.key === 'Enter') handleClickSubmit();
  };

  return (
    <Paper className={classes.paper}>
      <Typography variant='h5' className={classes.header}>
        {'Sign Up'}
      </Typography>
      <TextField
        variant='outlined'
        label={'Username'}
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

      <TextField
        variant='outlined'
        label={'Email'}
        className={classes.input}
        value={email}
        error={Boolean(emailError)}
        onChange={handleChangeEmail}
        onKeyPress={handlePressKey}
      />
      <TextField
        variant='outlined'
        label={'Phone number'}
        className={classes.input}
        value={phone}
        error={Boolean(phoneError)}
        onChange={handleChangePhone}
        onKeyPress={handlePressKey}
      />

      <Button
        variant='contained'
        color='primary'
        className={classes.button}
        startIcon={<Loading visible={isLoading} />}
        onClick={handleClickSubmit}
      >
        {'Sign Up'}
      </Button>

      <Typography className={classes.link}>
        {`Already have an account?`} <Link to={publicRoute.signIn.path}>{'Sign In'}</Link>
      </Typography>
    </Paper>
  );
};

const useStyles = makeStyles((theme) => ({
  header: {
    alignSelf: 'center',
    marginTop: theme.spacing(2),
    marginBottom: theme.spacing(2),
  },
  title: {
    alignSelf: 'center',
    marginBottom: theme.spacing(2),
    textTransform: 'uppercase',
  },
  paper: {
    padding: 24,
    width: 420,
    display: 'inline-flex',
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

export default SignUpForm;
