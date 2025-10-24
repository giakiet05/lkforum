import Home from './pages/Home.svelte';
import Login from './pages/Login.svelte';
import Register from './pages/Register.svelte';
import Profile from './pages/Profile.svelte';
import Settings from './pages/Settings.svelte';
import PostDetail from './pages/PostDetail.svelte';
import Community from './pages/Community.svelte';
import ManageCommunities from './pages/ManageCommunities.svelte';
import ModTools from './pages/ModTools.svelte';
import Messages from './pages/Messages.svelte';

const routes = {
    '/': Home,
    '/login': Login,
    '/register': Register,
    '/profile': Profile,
    '/settings': Settings,
    '/post/:id': PostDetail,
    '/lk/:name': Community,
    '/lk/:name/mod': ModTools,
    '/communities/manage': ManageCommunities,
    '/messages': Messages,
};

export default routes;
