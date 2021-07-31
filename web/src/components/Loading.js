import React from 'react';
import { CircularProgress } from '@material-ui/core';

const Loading = ({ visible, size = 20, icon, ...props }) =>
  visible ? <CircularProgress {...props} size={size} color='inherit' /> : icon ?? null;

export default Loading;
