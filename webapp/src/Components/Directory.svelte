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

<li on:click={() => expanded = !expanded}>{info.name}</li>
{#if expanded}
  <ul transition:slide="{{ duration: 300, easing: quartOut }}">
    {#each info.files as m}
        <svelte:self 
          info={m}
          on:updatePath={updatePath} />
    {/each}
    {#each Object.keys(info.subtitles) as s}
      <Subtitle
        subtitle={{ language: s, name: info.subtitles[s] }}
        parent={info}
        on:updatePath={updatePath} />
    {/each}
  </ul>
{/if}
