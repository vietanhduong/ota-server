import { createStore, combineReducers } from 'redux';
import profile from './profile';

export const store = createStore(
  combineReducers({
    profile,
  }),
  {},
);
