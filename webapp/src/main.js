import { gestures } from '@composi/gestures';
import App from './App.svelte';

gestures();

const app = new App({
	target: document.body,
});

export default app;
