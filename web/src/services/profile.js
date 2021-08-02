const { client } = require('./axios');

const getProfiles = () => client.get('/api/v1/profiles');

export const profileService = {
  getProfiles: getProfiles,
};
