<script lang="ts">
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import '../app.css'
  
  // Get the error from the page store
  $: error = $page.error;
  $: status = $page.status;

  const goHome = () => {
    goto('/');
  };
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-50">
  <div class="max-w-md w-full text-center px-6">
    <!-- Error code -->
    <div class="text-6xl font-bold text-gray-900 mb-4">
      {status || 404}
    </div>
    
    <!-- Error message -->
    <h1 class="text-2xl font-semibold text-gray-900 mb-2">
      {#if status === 404}
        Page Not Found
      {:else}
        Something went wrong
      {/if}
    </h1>
    
    <p class="text-gray-600 mb-8">
      {#if status === 404}
        The page you're looking for doesn't exist.
      {:else}
        {error?.message || 'An unexpected error occurred.'}
      {/if}
    </p>
    
    <!-- Action button -->
    <button
      on:click={goHome}
      class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-6 rounded-lg transition-colors duration-200"
    >
      Go Home
    </button>
  </div>
</div>