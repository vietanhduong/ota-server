import sha256 from 'js-sha256';
import jwt_decode from 'jwt-decode';

export const getDownloadUrl = ({ profile_id }) =>
  `itms-services://?action=download-manifest&amp;` +
  `url=${process.env.REACT_APP_HOST}/api/v1/profiles/ios/${profile_id}/manifest.plist`;

export { sha256, jwt_decode };
