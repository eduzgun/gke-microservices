<script lang="ts">
	import { onMount } from 'svelte';
	import type { Philosopher, Interaction } from '$lib/types';
	import { page } from "$app/state";
	import { interactionApi } from '$lib/api/interactions'; // Import both APIs
	import { philosopherApi } from '$lib/api/philosophers';
	import "../../../app.css";
	import Comments from './comments.svelte';

	let interactions = $state<Interaction[]>([]);
	let philosopher = $state<Philosopher>();

	let loading = $state(true);
	let error = $state<string | null>(null);

	let id = $derived(page.params.id);


	onMount(async () => {
		try {
			loading = true;
			error = null;

			if (!id) {
				error = 'Philosopher ID is required';
				return;
			}

			philosopher = await philosopherApi.getById(id);

			if (!philosopher) {
				error = 'Philosopher not found';
				return;
			}

		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load philosopher';
			console.error('Page load error:', err);
		} finally {
			loading = false;
		}

		try {
			if (!id) {
				error = 'Philosopher ID is required';
				return;
			}
			
			interactions = await interactionApi.get(id);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load interactions';
			console.error('Page load error:', err);
		} finally {
			loading = false;
		}
	});
</script>


<svelte:head>
	{#if philosopher}
		<title>{philosopher.name} - Great Philosophers</title>
		<meta name="description" content="Learn about {philosopher.name}, {philosopher.birthplace}" />
	{/if}
</svelte:head>

<div class="min-h-screen bg-gray-50 py-8">
	<div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
		
		<!-- Back Navigation -->
		<div class="mb-8">
			<a 
				href="/" 
				class="inline-flex items-center text-blue-600 hover:text-blue-800 transition-colors"
			>
				<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
				</svg>
				Back to Philosophers
			</a>
		</div>

		{#if loading}
			<div class="flex justify-center items-center py-20">
				<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
			</div>

		{:else if error}
			<div class="max-w-md mx-auto bg-red-50 border border-red-200 rounded-lg p-6 text-center">
				<h3 class="text-lg font-semibold text-red-900 mb-2">Error Loading Philosopher</h3>
				<p class="text-red-700">{error}</p>
				<a href="/" class="mt-4 inline-block text-blue-600 hover:text-blue-800">
					Return to main page
				</a>
			</div>

		{:else if !philosopher}
			<div class="text-center py-20">
				<h3 class="text-xl font-semibold text-gray-900 mb-2">Philosopher Not Found</h3>
				<p class="text-gray-600 mb-4">The requested philosopher could not be found.</p>
				<a href="/" class="text-blue-600 hover:text-blue-800">
					Return to main page
				</a>
			</div>

		{:else}
			<div class="bg-white shadow-lg rounded-lg overflow-hidden">
				
				<!-- Header Section -->
				<div class="relative">
					{#if philosopher.portrait_uri}
						<div class="h-64 md:h-80 bg-gradient-to-r from-blue-900 to-purple-900">
							<img 
								src={philosopher.portrait_uri} 
								alt="{philosopher.name} portrait"
								class="w-full h-full object-cover opacity-90"
							/>
							<div class="absolute inset-0 bg-gradient-to-t from-black/50 to-transparent"></div>
						</div>
					{:else}
						<div class="h-64 md:h-80 bg-gradient-to-r from-blue-900 to-purple-900 flex items-center justify-center">
							<div class="text-white text-6xl font-bold opacity-50">
								{philosopher.name.charAt(0)}
							</div>
						</div>
					{/if}
					
					<div class="absolute bottom-0 left-0 right-0 p-6 text-white">
						<h1 class="text-4xl md:text-5xl font-bold mb-2">
							{philosopher.name}
						</h1>
						<div class="text-lg md:text-xl opacity-90">
							{philosopher.date_born} - {philosopher.date_died}
						</div>
						<div class="text-base md:text-lg opacity-80">
							{philosopher.birthplace}
						</div>
					</div>
				</div>

				<!-- Content Section -->
				<div class="p-6 md:p-8">
					
					<!-- Biography -->
					<div class="mb-8">
						<h2 class="text-2xl font-bold text-gray-900 mb-4">Biography</h2>
						<p class="text-gray-700 leading-relaxed text-lg">
							{philosopher.bio}
						</p>
					</div>

					<!-- Interests -->
					{#if philosopher.interests && philosopher.interests.length > 0}
						<div class="mb-8">
							<h2 class="text-2xl font-bold text-gray-900 mb-4">Areas of Interest</h2>
							<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
								{#each philosopher.interests as interest}
									<div class="bg-blue-50 border border-blue-200 rounded-lg p-3 text-center">
										<span class="text-blue-800 font-medium">
											{interest}
										</span>
									</div>
								{/each}
							</div>
						</div>
					{/if}

					<!-- Quick Facts -->
					<div class="mb-8">
						<h2 class="text-2xl font-bold text-gray-900 mb-4">Quick Facts</h2>
						<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
							<div class="bg-gray-50 rounded-lg p-4">
								<h3 class="font-semibold text-gray-900 mb-2">Born</h3>
								<p class="text-gray-700">{philosopher.date_born}</p>
							</div>
							<div class="bg-gray-50 rounded-lg p-4">
								<h3 class="font-semibold text-gray-900 mb-2">Died</h3>
								<p class="text-gray-700">{philosopher.date_died}</p>
							</div>
							<div class="bg-gray-50 rounded-lg p-4">
								<h3 class="font-semibold text-gray-900 mb-2">Birthplace</h3>
								<p class="text-gray-700">{philosopher.birthplace}</p>
							</div>
							<div class="bg-gray-50 rounded-lg p-4">
								<h3 class="font-semibold text-gray-900 mb-2">Key Areas</h3>
								<p class="text-gray-700">
									{philosopher.interests ? philosopher.interests.slice(0, 2).join(', ') : 'Philosophy'}
								</p>
							</div>
						</div>
					</div>

					<!-- Action Buttons -->
					<div class="border-t pt-6">
						<div class="flex flex-col sm:flex-row gap-4">
							<a 
								href="/"
								class="flex-1 bg-blue-600 text-white text-center py-3 px-6 rounded-lg hover:bg-blue-700 transition-colors"
							>
								View All Philosophers
							</a>
							<a 
								href="/philosophers/add"
								class="flex-1 border border-gray-300 text-gray-700 text-center py-3 px-6 rounded-lg hover:bg-gray-50 transition-colors"
							>
								Add New Philosopher
							</a>
						</div>
					</div>


					<!-- Comments Section -->
					<Comments 
						interactions={interactions} 
						philosopherId={philosopher.id}
					/>
				
				</div>
			</div>
		{/if}
	</div>
</div>