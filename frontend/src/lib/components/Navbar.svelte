<script lang="ts">
  import { authStore } from '$lib/stores/auth';
  import { goto } from '$app/navigation';

  const authState = $derived.by(() => $authStore);

  let showDropdown = $state(false);

  const handleLogout = async () => {
    await authStore.logout();
    showDropdown = false;
    goto('/');
  };

  const closeDropdown = () => {
    showDropdown = false;
  };
</script>

<nav class="bg-white shadow-lg border-b">
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
    <div class="flex justify-between items-center h-16">
      <!-- Logo -->
      <div class="flex items-center">
        <a href="/" class="text-xl font-bold text-gray-800 hover:text-blue-600">
          Philosophy App
        </a>
      </div>

      <!-- Auth Section -->
      <div class="flex items-center space-x-4">
        {#if authState.loading}
          <!-- Loading State -->
          <div class="flex items-center text-gray-500">
            <svg class="animate-spin -ml-1 mr-2 h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span class="text-sm">Loading...</span>
          </div>
        {:else if authState.user}
          <div class="relative">
            <button 
              onclick={() => showDropdown = !showDropdown}
              class="flex items-center text-gray-700 hover:text-blue-600 focus:outline-none"
            >
              <span class="text-sm font-medium">Welcome, {authState.user.username}</span>
              <svg class="ml-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            </button>

            {#if showDropdown}
              <div 
                class="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg py-1 z-50 border"
                role="menu"
              >
                <hr class="my-1">
                <button 
                  onclick={handleLogout}
                  class="block w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50"
                >
                  Logout
                </button>
              </div>
            {/if}
          </div>
        {:else}
          <div class="flex items-center space-x-3">
            <a 
              href="/login" 
              class="text-gray-600 hover:text-blue-600 px-3 py-2 rounded-md transition-colors"
            >
              Login
            </a>
            <a 
              href="/register" 
              class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md transition-colors font-medium"
            >
              Register
            </a>
          </div>
        {/if}
      </div>
    </div>
  </div>
</nav>

<!-- Click outside dropdown -->
{#if showDropdown}
  <div class="fixed inset-0 z-40" onclick={closeDropdown}></div>
{/if}