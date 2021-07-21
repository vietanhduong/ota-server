import { createTheme, MuiThemeProvider } from '@material-ui/core';
import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';

ReactDOM.render(
  <React.StrictMode>
    <MuiThemeProvider
      theme={createTheme({
        palette: {
          primary: {
            main: '#007aff',
          },
        },
      })}
    >
      <App />
    </MuiThemeProvider>
  </React.StrictMode>,
  document.getElementById('root'),
);
