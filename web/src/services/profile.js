const { client } = require('./axios');

const getProfile = () => client.get('/api/v1/profiles');

export const profileService = {
  getProfile,
};
