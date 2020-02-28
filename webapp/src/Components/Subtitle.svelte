<script>
  import { createEventDispatcher } from 'svelte';
  import { slideUp } from '../stores.js';

  export let subtitle;
  export let parent;

  const dispatch = createEventDispatcher();

  const updateSubtitle = () => {
    dispatch('updatePath', {
      path: subtitle.name,
    });
    slideUp.openAndUpdate(subtitle, parent);
  };
</script>

<span on:click={updateSubtitle}>Sub: {subtitle.name} ({subtitle.language})</span><br>
