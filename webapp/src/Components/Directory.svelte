<script>
  import { createEventDispatcher } from 'svelte';
  import { slide } from 'svelte/transition';
  import { quartOut } from 'svelte/easing';
  import Subtitle from './Subtitle.svelte';

  export let info;

  const dispatch = createEventDispatcher();
  
  let expanded = false;

  const updatePath = (e) => {
    dispatch('updatePath', {
      path: info.name + '/' + e.detail.path,
    });
  };
</script>

<li
  class="folder bg-white hidden"
  class:show={info.show}>
  <div
    class="py-2 px-4"
    on:click={() => expanded = !expanded}>
    <span class="font-mono absolute">{expanded ? '-' : '+'}</span>
    <span class="pl-5">{info.name}</span>
  </div>

  <!-- Subtitles -->
  <ul
    class="hidden px-2 ml-2 -mt-2 mb-2 bg-white rounded rounded-t-none"
    class:expanded={expanded}
    transition:slide="{{ duration: 300, easing: quartOut }}">
    {#each Object.keys(info.subtitles) as s}
      <Subtitle
        subtitle={{ language: info.subtitles[s], name: s }}
        parent={info}
        on:updatePath={updatePath} />
    {:else}
      {#if !info.files.length}
        <li class="px-2 py-1">No subtitles found.</li>
      {/if}
    {/each}
  </ul>

  <!-- Subdirectories -->
  <ul
    class="hidden px-2"
    class:expanded={expanded}
    transition:slide="{{ duration: 300, easing: quartOut }}">
    {#each info.files as m, i (i)}
        <svelte:self 
          info={m}
          on:updatePath={updatePath} />
    {/each}
  </ul>
</li>

<style>
  .show {
    display: list-item;
  }

  .expanded {
    display: block;
  }

  .folder:not(:last-child) {
    @apply border-b border-gray-400
  }

  .folder:first-child {
    @apply rounded-t
  }

  .folder:last-child {
    @apply rounded-b
  }
</style>
