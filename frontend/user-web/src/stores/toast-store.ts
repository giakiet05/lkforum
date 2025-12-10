import { writable } from 'svelte/store';

type ToastType = 'success' | 'error' | 'warning' | 'info';

interface Toast {
  id: string;
  message: string;
  type: ToastType;
  duration?: number;
}

function createToastStore() {
  const { subscribe, update } = writable<Toast[]>([]);

  return {
    subscribe,
    show(message: string, type: ToastType = 'info', duration = 3000) {
      const id = Math.random().toString(36).substring(2, 9);
      const toast: Toast = { id, message, type, duration };
      
      update(toasts => [...toasts, toast]);

      if (duration > 0) {
        setTimeout(() => {
          this.remove(id);
        }, duration);
      }
    },
    success(message: string, duration = 3000) {
      this.show(message, 'success', duration);
    },
    error(message: string, duration = 4000) {
      this.show(message, 'error', duration);
    },
    warning(message: string, duration = 3500) {
      this.show(message, 'warning', duration);
    },
    info(message: string, duration = 3000) {
      this.show(message, 'info', duration);
    },
    remove(id: string) {
      update(toasts => toasts.filter(t => t.id !== id));
    },
    clear() {
      update(() => []);
    }
  };
}

export const toastStore = createToastStore();
export type { Toast, ToastType };
