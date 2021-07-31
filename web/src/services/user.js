const { client } = require('./axios');

const login = (body) => client.post(`/api/v1/users/login`, { ...body });
const register = (body) => client.post(`/api/v1/users/register`, { ...body });
const refresh = (headers) => client.post(`/api/v1/users/refresh-token`, {}, { headers });
const logout = (headers) => client.post(`/api/v1/users/logout`, {}, {headers});

export const userService = {
  login,
  register,
  refresh,
  logout
};
