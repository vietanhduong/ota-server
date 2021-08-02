import { createTheme, ThemeProvider } from '@material-ui/core/styles';
import { common } from '@material-ui/core/colors';

export const appTheme = createTheme({
  props: {
    MuiTextField: {
      InputLabelProps: { shrink: true },
      inputProps: { autoSave: 'false' },
    },
    MuiButton: {},
    MuiChip: {
      variant: 'outlined',
    },
    MuiAvatar: {
      variant: 'rounded',
    },
    MuiTypography: {
      component: 'div',
    },
    MuiInputBase: {
      style: {
        backgroundColor: common.white,
      },
    },
  },
  typography: {
    subtitle1: {
      fontWeight: 500,
      lineHeight: 1.5,
    },
    subtitle2: {
      fontWeight: 500,
      lineHeight: 1.43,
    },
    button: {
      textTransform: 'none',
    },
  },
  overrides: {
    MuiTab: {
      root: {
        textTransform: 'none',
      },
    },
    MuiButton: {
      root: {
        textTransform: 'none',
      },
    },
    MuiTooltip: {
      tooltip: {
        fontSize: '0.725rem',
      },
    },
  },
  palette: {
    primary: {
      main: '#007aff',
    },
  },
});

const Theme = ({ children }) => {
  return <ThemeProvider theme={appTheme}>{children}</ThemeProvider>;
};

export default Theme;
