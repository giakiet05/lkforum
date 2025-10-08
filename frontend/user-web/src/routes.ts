import Home from './pages/Home.svelte';
import Login from './pages/Login.svelte';
import Register from './pages/Register.svelte';
const routes = {
    '/': Home,
    '/login': Login,
    '/register': Register,
};

export default routes;
