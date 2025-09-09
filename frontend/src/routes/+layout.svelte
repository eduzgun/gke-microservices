<script lang="ts">
  import { onMount } from 'svelte';
  import { authStore } from '$lib/stores/auth';
  import favicon from '$lib/assets/favicon.svg';
  import Navbar from '$lib/components/Navbar.svelte';
  
  onMount(() => {
    authStore.init();
  });
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
  <title>Philosophers app</title>
</svelte:head>

<Navbar />
<main class="min-h-screen bg-gray-50">
  <slot />
</main>