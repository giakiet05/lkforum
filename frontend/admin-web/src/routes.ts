import Login from "./pages/Login.svelte";
import Dashboard from "./pages/Dashboard.svelte";
import Users from "./pages/Users.svelte";
import Communities from "./pages/Communities.svelte";
import Reports from "./pages/Reports.svelte";

export const routes = {
  "/": Dashboard,
  "/login": Login,
  "/dashboard": Dashboard,
  "/users": Users,
  "/communities": Communities,
  "/reports": Reports,
};
