<script lang="ts">
	import { onMount } from 'svelte';
	import { loadPhilosophers } from '../lib/api';
	import type { Philosopher } from '../lib/types';
    import "../app.css";

	let philosophers = $state<Philosopher[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);

	onMount(async () => {
		try {
			loading = true;
			error = null;
			philosophers = await loadPhilosophers();
		} catch (err) {
			error = 'Failed to load philosophers';
			console.error(err);
		} finally {
			loading = false;
		}
	});

    function formatPhilosopherDate(dateStr: string) {
        if (!dateStr) return '?';
        
        // If it contains BC or AD, return as-is
        if (dateStr.includes('BC') || dateStr.includes('AD')) {
            return dateStr;
        }
        
        // Try to parse as regular date
        const date = new Date(dateStr);
        if (!isNaN(date.getTime())) {
            return date.getFullYear().toString();
        }
        
        // Fallback to original string
        return dateStr;
        }
</script>

<div class="min-h-screen bg-gray-50 py-8">
	<div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
		<div class="text-center mb-12">
			<h1 class="text-4xl font-bold text-gray-900 mb-4">Great Philosophers</h1>
			<p class="text-lg text-gray-600 max-w-2xl mx-auto">
				Explore the minds that shaped human thought throughout history
			</p>
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
				<p class="text-gray-600">Check back later for more content.</p>
			</div>
		{:else}
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
				{#each philosophers as philosopher}
					<div class="bg-white rounded-lg shadow-md hover:shadow-lg transition-shadow duration-300 overflow-hidden">
						{#if philosopher.portrait_uri}
							<div class="h-48 bg-gray-200">
								<img 
									src={philosopher.portrait_uri} 
									alt="{philosopher.name} portrait"
									class="w-full h-full object-cover"
									loading="lazy"
								/>
							</div>
						{/if}
						
						<div class="p-6">
							<h2 class="text-xl font-bold text-gray-900 mb-3">
								{philosopher.name}
							</h2>
							
							<div class="flex items-center text-sm text-gray-600 mb-3">
								{philosopher.date_born} - 
								{philosopher.date_died}
							</div>
							
							<div class="flex items-center text-xs text-gray-500 mb-3">
								{philosopher.birthplace}
							</div>
							
							<p class="text-gray-700 text-sm mb-4 line-clamp-3">
								{philosopher.bio}
							</p>
							
							{#if philosopher.interests && philosopher.interests.length > 0}
								<div class="mb-4">
									<h4 class="text-sm font-semibold text-gray-900 mb-2">Interests:</h4>
									<div class="flex flex-wrap gap-1">
										{#each philosopher.interests.slice(0, 3) as interest}
											<span class="inline-block bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded-full">
												{interest}
											</span>
										{/each}
										{#if philosopher.interests.length > 3}
											<span class="inline-block text-gray-500 text-xs px-2 py-1">
												+{philosopher.interests.length - 3} more
											</span>
										{/if}
									</div>
								</div>
							{/if}
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>

<style>
	.line-clamp-3 {
		display: -webkit-box;
		-webkit-line-clamp: 3;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
</style>