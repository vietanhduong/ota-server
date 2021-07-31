import axios from 'axios';
import { store } from 'reducers';

const client = axios.create({
  baseURL: process.env.REACT_APP_HOST,
});

client.interceptors.request.use((config) => {
  if (!config.headers.Authorization) {
    config.headers.Authorization = `token ${store.getState().profile.access_token}`;
  }
  return config;
});

client.interceptors.response.use(
  ({ data }) => data,
  ({ response }) => {
    return Promise.resolve({ status: 0, detail: response?.data });
  },
);

export { client };
