const { client } = require('./axios');

const login = (body) => client.post(`/api/v1/users/login`, { ...body });
const register = (body) => client.post(`/api/v1/users/registerlogin`, { ...body });
const refresh = (headers) => client.post(`/api/v1/users/refresh-token`, {}, { headers });

export const userService = {
  login,
  register,
  refresh,
};
