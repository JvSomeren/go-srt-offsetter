<script>
	import { onMount } from 'svelte';
	import Tailwindcss from './Tailwindcss.svelte';
	import Directory from './Components/Directory.svelte';
	import SlideUp from './Components/SlideUp.svelte';
  import { slideUp } from './stores';

	let media = [];

	const handleUpdatePath = (e) => {
		slideUp.updatePath(e.detail.path);
	}

	onMount(async () => {
		const res = await fetch(`${process.env.apiUrl}/media`);
			
		try {
			media = await res.json();
		} catch(e) {
			console.error('no json');
		}
	});
</script>

<header class="p-2 pt-3">
	<h1 class="text-3xl font-bold text-center font-mono">Subtitle offsetter</h1>
</header>

<main>
	<section>
		<ul>
			{#each media as m, i }
				<Directory
					info={m}
					on:updatePath={handleUpdatePath} />
			{/each}
		</ul>
	</section>
</main>

<SlideUp />

<style>
	
</style>
