import sha256 from 'js-sha256';
import jwt_decode from 'jwt-decode';
import {profileAction} from "../actions/profile";

const getHost = () => {
  return process.env.REACT_APP_HOST || window.location.origin;
}

const getExchangeCode = () => {
  console.log("call");
  const raw = localStorage.getItem(profileAction.PROFILE);
  if (raw.length === 0) return "";
  const profile = JSON.parse(raw);
  return profile.exchange_code;
}

export const getDownloadUrl = ({profile_id}) =>
  `itms-services://?action=download-manifest&amp;` +
  `url=${getHost()}/api/v1/profiles/ios/${profile_id}/manifest.plist?code=${getExchangeCode()}`;

export {sha256, jwt_decode};
