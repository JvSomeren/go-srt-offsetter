<script>
	import { onDestroy } from 'svelte';
  import { slide, fade } from 'svelte/transition';
	import { quartOut } from 'svelte/easing';
	import { slideUp } from '../stores';
	import { debounce, languages } from '../helpers';
	
	let offsetContainer = {};

	const closeSlideUp = () => {
		slideUp.close();
	};

	const updateOffset = debounce(function(path) {
		if(offsetContainer[path].deltaOffset !== 0) {
			const data = {
				subtitle: path,
				offset: offsetContainer[path].deltaOffset,
			};
			const res = fetch(`${process.env.apiUrl}/subtitle`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify(data),
			})
				.then(res => res.json())
				.catch(error => console.error(error));

			offsetContainer[path].deltaOffset = 0;
		}
	}, 500);

	const onClick = (offset) => {
		offsetContainer[$slideUp.path].deltaOffset += offset;
		offsetContainer[$slideUp.path].currentOffset += offset;

		updateOffset($slideUp.path);
	};

	const unsubscribeSlideUp = slideUp.subscribe((v) => {
		if(!offsetContainer.hasOwnProperty($slideUp.path)) {
			offsetContainer[$slideUp.path] = {
				deltaOffset: 0,
				currentOffset: 0,
			};
		}
	});

	onDestroy(async () => {
		unsubscribeSlideUp();
	});
</script>

{#if $slideUp.open}
  <div
    class="fixed-full bg-black opacity-25"
    transition:fade="{{ duration: 250, easing: quartOut }}"
    on:click={closeSlideUp}
		on:swipedown={closeSlideUp}></div>
  <aside
    class="fixed-bottom bg-white rounded-t-large shadow-t-lg pt-3 pb-5 px-3"
    transition:slide="{{ duration: 300, easing: quartOut }}"
		on:swipedown={closeSlideUp}>
    <h2 class="text-2xl font-medium truncate pl-2">
			{$slideUp.parent.name}
		</h2>
		<h3 class="text-sm rounded-full bg-gray-200 text-gray-700 font-semibold mx-auto inline-block px-3 py-1 lowercase">{languages[$slideUp.subtitle.language]}</h3>

		<div class="px-3 py-4">
			<p>Current offset: {offsetContainer[$slideUp.path].currentOffset}ms</p>
		</div>

    <div class="flex justify-around mx-auto mt-2 text-xs">
      <button class="btn btn-red" on:click={() => onClick(-100)}>-100ms</button>
      <button class="btn btn-red" on:click={() => onClick(-50)}>-50ms</button>
      <button class="btn btn-green" on:click={() => onClick(50)}>+50ms</button>
      <button class="btn btn-green" on:click={() => onClick(100)}>+100ms</button>
    </div>
  </aside>
{/if}

<style>
	.fixed-bottom {
		@apply fixed bottom-0 left-0 right-0
	}

	.fixed-full {
		@apply fixed-bottom top-0
	}

	.btn {
		@apply py-1 px-3 rounded outline-none
	}

	.btn-red {
		@apply bg-red-500 text-white
	}

	.btn-red:hover {
		@apply bg-red-600
	}

	.btn-green {
		@apply bg-green-500 text-white
	}

	.btn-green:hover {
		@apply bg-green-600
	}
</style>