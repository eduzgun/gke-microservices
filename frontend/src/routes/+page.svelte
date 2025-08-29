<script lang="ts">
	import { onMount } from 'svelte';
	import { philosopherApi } from '$lib/api/philosophers';
	import type { Philosopher } from '$lib/types';
	import "../app.css";

	let philosophers = $state<Philosopher[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);

	onMount(async () => {
		try {
			loading = true;
			error = null;
			philosophers = await philosopherApi.getAll();
		} catch (err) {
			error = 'Failed to load philosophers';
			console.error(err);
		} finally {
			loading = false;
		}
	});

	// Function to truncate bio text
	function truncateBio(bio: string, maxLength: number = 120): string {
		if (bio.length <= maxLength) return bio;
		return bio.slice(0, maxLength).trim() + '...';
	}
</script>

<div class="min-h-screen bg-gray-50 py-8">
	<div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
		<!-- Header -->
		<div class="text-center mb-12">
			<h1 class="text-4xl font-bold text-gray-900 mb-4">Great Philosophers</h1>
			<p class="text-lg text-gray-600 max-w-2xl mx-auto">
				Explore the minds that shaped human thought throughout history
			</p>
			<!-- Add New Philosopher Button -->
			<div class="mt-6">
				<a 
					href="/philosophers/add" 
					class="inline-flex items-center px-6 py-3 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 transition-colors"
				>
					<svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
					</svg>
					Add New Philosopher
				</a>
			</div>
		</div>

		{#if loading}
			<div class="flex justify-center items-center py-20">
				<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
			</div>

		{:else if error}
			<div class="max-w-md mx-auto bg-red-50 border border-red-200 rounded-lg p-6 text-center">
				<h3 class="text-lg font-semibold text-red-900 mb-2">Error Loading Data</h3>
				<p class="text-red-700">{error}</p>
			</div>

		{:else if philosophers.length === 0}
			<div class="text-center py-20">
				<h3 class="text-xl font-semibold text-gray-900 mb-2">No Philosophers Found</h3>
				<p class="text-gray-600 mb-4">Check back later for more content.</p>
				<a 
					href="/philosophers/add" 
					class="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
				>
					Add the first philosopher
				</a>
			</div>

		{:else}
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
				{#each philosophers as philosopher}
					<a 
						href="/philosophers/{philosopher.id}" 
						class="block bg-white rounded-lg shadow-md hover:shadow-xl transition-all duration-300 overflow-hidden transform hover:-translate-y-1 group"
					>
						<!-- Portrait -->
						{#if philosopher.portrait_uri}
							<div class="h-48 bg-gray-200 overflow-hidden">
								<img
									src={philosopher.portrait_uri}
									alt="{philosopher.name} portrait"
									class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
									loading="lazy"
								/>
							</div>
						{:else}
							<div class="h-48 bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center">
								<div class="text-white text-6xl font-bold opacity-80">
									{philosopher.name.charAt(0)}
								</div>
							</div>
						{/if}
						
						<!-- Content -->
						<div class="p-6">
							<!-- Name -->
							<h2 class="text-xl font-bold text-gray-900 mb-3 group-hover:text-blue-600 transition-colors">
								{philosopher.name}
							</h2>
							
							<!-- Dates -->
							<div class="flex items-center text-sm text-gray-600 mb-2">
								<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
								</svg>
								{philosopher.date_born} - {philosopher.date_died}
							</div>
							
							<!-- Birthplace -->
							<div class="flex items-center text-sm text-gray-500 mb-3">
								<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
								</svg>
								{philosopher.birthplace}
							</div>
							
							<!-- Bio Preview -->
							<p class="text-gray-700 text-sm mb-4 leading-relaxed">
								{truncateBio(philosopher.bio)}
							</p>
							
							<!-- Interests Preview -->
							{#if philosopher.interests && philosopher.interests.length > 0}
								<div class="mb-4">
									<div class="flex flex-wrap gap-1">
										{#each philosopher.interests.slice(0, 2) as interest}
											<span class="inline-block bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded-full">
												{interest}
											</span>
										{/each}
										{#if philosopher.interests.length > 2}
											<span class="inline-block text-gray-500 text-xs px-2 py-1">
												+{philosopher.interests.length - 2} more
											</span>
										{/if}
									</div>
								</div>
							{/if}

							<!-- Read More Indicator -->
							<div class="flex items-center text-blue-600 text-sm font-medium group-hover:text-blue-800 transition-colors">
								<span>Learn more</span>
								<svg class="w-4 h-4 ml-1 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
								</svg>
							</div>
						</div>
					</a>
				{/each}
			</div>
		{/if}
	</div>
</div>