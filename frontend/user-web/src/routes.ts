import Home from './pages/Home.svelte';
import Popular from './pages/Popular.svelte';
import Explore from './pages/Explore.svelte';
import All from './pages/All.svelte';
import Login from './pages/Login.svelte';
import Register from './pages/Register.svelte';
import Profile from './pages/Profile.svelte';
import Settings from './pages/Settings.svelte';
import PostDetail from './pages/PostDetail.svelte';
import EditPost from './pages/EditPost.svelte';
import Community from './pages/Community.svelte';
import ManageCommunities from './pages/ManageCommunities.svelte';
import ModTools from './pages/ModTools.svelte';
import CommunitySettings from './pages/CommunitySettings.svelte';
import Messages from './pages/Messages.svelte';
import GoogleCallback from './pages/GoogleCallback.svelte';
import GoogleSetup from './pages/GoogleSetup.svelte';
import AuthError from './pages/AuthError.svelte';
import ForgotPassword from './pages/ForgotPassword.svelte';
import Search from './pages/Search.svelte';

const routes = {
    '/': Home,
    '/popular': Popular,
    '/explore': Explore,
    '/all': All,
    '/search': Search,
    '/login': Login,
    '/register': Register,
    '/forgot-password': ForgotPassword,
    '/profile': Profile,
    '/profile/:username': Profile,
    '/settings': Settings,
    '/post/:slugId': PostDetail,
    '/post/:slugId/edit': EditPost,
    '/lk/:name': Community,
    '/lk/:name/mod': ModTools,
    '/lk/:name/settings': CommunitySettings,
    '/communities/manage': ManageCommunities,
    '/messages': Messages,
    '/auth/callback': GoogleCallback,
    '/auth/google-setup': GoogleSetup,
    '/auth/error': AuthError,
};

export default routes;
