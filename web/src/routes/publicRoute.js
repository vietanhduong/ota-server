import { SignInForm } from 'views/SignIn';
import { SignUpForm } from 'views/SignUp';

const publicRoute = {
  signIn: {
    path: '/auth/signIn',
    component: SignInForm,
  },
  signUp: {
    path: '/auth/signUp',
    component: SignUpForm,
  },
};

export default publicRoute;
