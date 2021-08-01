import { SignInForm } from 'views/SignIn';
import { SignUpForm } from 'views/SignUp';

const publicRoute = {
  signIn: {
    path: '/auth/sign-in',
    component: SignInForm,
  },
  signUp: {
    path: '/auth/sign-up',
    component: SignUpForm,
  },
};

export default publicRoute;
