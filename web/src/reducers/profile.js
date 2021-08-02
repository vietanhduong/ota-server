export const ActionType = {
  USER_LOGIN: 'USER_LOGIN',
  USER_LOGOUT: 'USER_LOGOUT',
};

const profile = (state = {}, { type, data }) => {
  switch (type) {
    case ActionType.USER_LOGIN:
      return { ...data, isLoggedIn: true };
    case ActionType.USER_LOGOUT:
      return {};
    default:
      return { ...state };
  }
};

export default profile;
