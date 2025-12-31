import { mount } from 'svelte'
import './app.css'
import App from './App.svelte'
// Intercept all fetch calls to debug logout spam
const originalFetch = window.fetch;
window.fetch = function(...args) {
  const url = args[0]?.toString() || '';
  if (url.includes('/logout')) {
    console.error('[FETCH INTERCEPTOR] Logout call detected!', url);
    console.trace('[FETCH INTERCEPTOR] Stack trace:');
  }
  return originalFetch.apply(this, args);
};
const app = mount(App, {
  target: document.getElementById('app')!,
})

export default app
