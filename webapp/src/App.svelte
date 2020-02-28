<script>
	import { onMount } from 'svelte';
	import Directory from './Components/Directory.svelte';
	import SlideUp from './Components/SlideUp.svelte';
  import { slideUp } from './stores.js';

	let media = [];

	const handleUpdatePath = (e) => {
		slideUp.updatePath(e.detail.path);

		// const data = {
		// 	subtitle: event.detail.path,
		// 	offset: event.detail.offset,
		// };
		// const res = fetch(`${process.env.apiUrl}/subtitle`, {
		// 	method: 'PUT',
		// 	headers: {
		// 		'Content-Type': 'application/json',
		// 	},
		// 	body: JSON.stringify(data),
		// })
		// 	.then(res => res.json())
		// 	.then(res => console.log(res));
	}

	onMount(async () => {
		const res = await fetch(`${process.env.apiUrl}/media`);
		media = await res.json();
	})
</script>

<main>
	<h1>Subtitle offsetter</h1>

	<section>
		<ul>
			{#each media as m, i }
				<Directory
					info={m}
					on:updatePath={handleUpdatePath} />
			{/each}
		</ul>
	</section>

	<SlideUp />
</main>

<style>

</style>
