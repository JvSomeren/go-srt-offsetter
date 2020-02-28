<script>
  import { slide, fade } from 'svelte/transition';
  import { quartOut } from 'svelte/easing';
  import { slideUp } from '../stores.js';

	const closeSlideUp = () => {
		slideUp.close();
	};

	const onSwipedown = (e) => {
		console.log(e);
	};
</script>

{#if $slideUp.open}
  <div
    class="backdrop"
    transition:fade="{{ duration: 250, easing: quartOut }}"
    on:click={closeSlideUp}
		on:swipedown={closeSlideUp}></div>
  <aside
    class="slide-up"
    transition:slide="{{ duration: 300, easing: quartOut }}"
		on:swipedown={closeSlideUp}>
    <h2>{$slideUp.parent.name}</h2>

    <div class="controls">
      <button>-100</button>
      <button>-50</button>
      <button>+50</button>
      <button>+100</button>
    </div>
  </aside>
{/if}

<style>
	.backdrop {
		position: fixed;
		top: 0;
		bottom: 0;
		left: 0;
		right: 0;

		background-color: rgba(0,0,0,0.2)
	}

	.slide-up {
		position: fixed;
		bottom: 0;
		left: 0;
		right: 0;

		text-align: center;

		--slide-up-border-radius: 12px;
		border-top-left-radius: var(--slide-up-border-radius);
		border-top-right-radius: var(--slide-up-border-radius);
		background-color: #fff;

		box-shadow: 0px -7px 21px -8px rgba(0,0,0,0.37);
	}

	.controls {
		display: flex;
		justify-content: space-around;
		width: 80%;
		margin: 0 auto;
	}
</style>