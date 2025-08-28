<script lang="ts">
  import { goto } from '$app/navigation';
  import { authApi } from '$lib/api/auth';

  import '../../app.css';
  
  let email = $state('');
  let password = $state('');
  let error = $state<string | null>(null);
  let loading = $state(false);

  const isValid = $derived({
    email: /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email),
    password: password.length >= 6,
  });

  $effect(() => {
    if (email || password) error = null;
  });

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    if (!isValid.email || !isValid.password) return;

    loading = true;
    error = null;

    try {
      await authApi.login({ email, password });
      goto('/profile');
    } catch (err: any) {
      error = err.message || 'Login failed';
    } finally {
      loading = false;
    }
  };
</script>

<main class="flex items-center justify-center min-h-screen bg-gray-50">
  <form onsubmit={handleSubmit} class="w-full max-w-sm p-6 bg-white rounded shadow">
    <h2 class="text-2xl font-bold mb-4">Log In</h2>

    {#if error}
      <div class="mb-4 p-2 text-sm text-red-600 bg-red-100 rounded">
        {error}
      </div>
    {/if}

    <div class="mb-4">
      <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
      <input
        id="email"
        type="email"
        class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring focus:ring-blue-200"
        bind:value={email}
        required
      />
      {#if email && !isValid.email}
        <p class="mt-1 text-sm text-red-500">Enter a valid email</p>
      {/if}
    </div>

    <div class="mb-6">
      <label for="password" class="block text-sm font-medium text-gray-700">Password</label>
      <input
        id="password"
        type="password"
        class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring focus:ring-blue-200"
        bind:value={password}
        required
      />
      {#if password && !isValid.password}
        <p class="mt-1 text-sm text-red-500">Password must be at least 6 characters</p>
      {/if}
    </div>

    <button
      type="submit"
      disabled={loading || !isValid.email || !isValid.password}
      class="w-full py-2 px-4 bg-blue-600 text-white font-medium rounded hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
    >
      {loading ? 'Logging in...' : 'Log In'}
    </button>
  </form>
</main>

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    font-family: system-ui, sans-serif;
  }
</style>