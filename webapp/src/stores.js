import { writable } from 'svelte/store';

function createSlideUp() {
  const { subscribe, set, update } = writable({
    subtitle: {},
    parent: {},
    path: '',
    open: false,
  });

  return {
    subscribe,
    updatePath: (subtitle) => update(u => ({ ...u, subtitle })),
    updatePath: (parent) => update(u => ({ ...u, parent })),
    updatePath: (path) => update(u => ({ ...u, path })),
    open: () => update(u => ({ ...u, open: true })),
    openAndUpdate: (subtitle, parent) => update(u => ({
      ...u,
      subtitle,
      parent,
      open: true,
    })),
    close: () => update(u => ({ ...u, open: false })),
    toggle: () => update(u => ({ ...u, open: !u.open })),
  }
}

export const slideUp = createSlideUp();
