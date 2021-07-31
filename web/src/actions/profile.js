import {store} from 'reducers';
import {ActionType} from 'reducers/profile';
import {userService} from 'services/user';

const PROFILE = 'profile';

const refresh = async (profile) => {
  await userService
    .refresh({Authorization: `token ${profile.refresh_token}`})
    .then(login)
    .then((nextProfile) => {
      setTimeout(() => {
        refresh(nextProfile);
      }, 60 * 55 * 1000);
    });
};

const login = (profile) => {
  store.dispatch({
    type: ActionType.USER_LOGIN,
    data: profile,
  });
  localStorage.setItem(PROFILE, JSON.stringify(profile));
  return profile;
};

const logout = () => {
  const raw = localStorage.getItem(PROFILE);
  if (raw.length === 0) return;

  const profile = JSON.parse(raw);
  userService.logout({Authorization: `token ${profile.access_token}`}).finally(() => {
    store.dispatch({
      type: ActionType.USER_LOGOUT,
    });

    localStorage.removeItem(PROFILE);
  });
};

export const profileAction = {
  ActionType,
  login,
  logout,
  refresh,
  PROFILE
};
