<script>
	import { onMount } from 'svelte';
	import Tailwindcss from './Tailwindcss.svelte';
	import Directory from './Components/Directory.svelte';
	import SlideUp from './Components/SlideUp.svelte';
	import { slideUp } from './stores';
	import { debounce } from './helpers';

	let media = [];
	let searchValue = '';

	const stringIncludes = (subject, search) => subject.toLowerCase().includes(search.toLowerCase());

	const handleSearch = (dir) => {
		let subdirIsMatched = false;

		for(const subdir of dir.files) {
			if(handleSearch(subdir)) subdirIsMatched = true;
		}

		// check if name matches search
		if(stringIncludes(dir.name, searchValue) || subdirIsMatched) {
			dir.show = true;
			return true;
		} else {
			dir.show = false;
			return false;
		}
	};

	const search = debounce(function() {
		for (const folder of media) {
			handleSearch(folder);
		}

		media = media; // used to force reactivity
	}, 300);

	const handleUpdatePath = (e) => {
		slideUp.updatePath(e.detail.path);
	}

	const prepareMedia = (r) => {
		for (const d of r) {
			d.files = prepareMedia(d.files);
			d.show = true;
		}

		return r;
	};

	onMount(async () => {
		const res = await fetch(`${process.env.apiUrl}/media`);
			
		try {
			let t = await res.json();

			for (const f of t) {
				f.files = prepareMedia(f.files);
				f.show = true;
			}

			media = t;
			console.log(media);
		} catch(e) {
			console.error('no json', e);
		}
	});
</script>

<header class="p-2 pt-3">
	<h1 class="text-3xl font-bold text-center font-mono">Subtitle offsetter</h1>
</header>

<section class="p-2">
	<!-- transition-colors duration-100 ease-in-out focus:outline-0 border border-transparent focus:bg-white focus:border-gray-300 placeholder-gray-600 rounded-lg bg-gray-200 py-2 pr-4 pl-10 block w-full appearance-none leading-normal ds-input -->
	<input
		type="text"
		class="shadow-lg bg-white rounded w-full px-4 py-2 outline-none"
		placeholder="Search for media"
		bind:value={searchValue}
		on:keyup={search}>
</section>

<main>
	<section class="p-2">
		<ul>
			{#each media as m, i (i) }
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
