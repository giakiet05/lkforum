import Home from './pages/Home.svelte';
import Login from './pages/Login.svelte';
import Register from './pages/Register.svelte';
import Profile from './pages/Profile.svelte';
import Settings from './pages/Settings.svelte';
import PostDetail from './pages/PostDetail.svelte';

const routes = {
    '/': Home,
    '/login': Login,
    '/register': Register,
    '/profile': Profile,
    '/settings': Settings,
    '/post/:id': PostDetail,
};

export default routes;
