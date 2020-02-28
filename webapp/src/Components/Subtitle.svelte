<script>
  import { createEventDispatcher } from 'svelte';
  import { slideUp } from '../stores';
  import { languages } from '../helpers';

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

<span on:click={updateSubtitle}>Sub: {languages[subtitle.language]}</span><br>
