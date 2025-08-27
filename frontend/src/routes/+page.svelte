<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from "$app/state";
	import { loadPhilosopher } from '$lib/api';
	import type { Philosopher } from '$lib/types';
	import { enhance } from '$app/forms';
	import "../../../app.css";

	// Types for comments and likes
	interface Comment {
		id: number;
		user_name: string;
		comment_text: string;
		created_at: string;
		philosopher_id: number;
	}

	interface LikeData {
		id: number;
		philosopher_id: number;
		like_count: number;
		user_has_liked: boolean;
	}

	// Existing state
	let philosopher = $state<Philosopher | null>(null);
	let loading = $state(true);
	let error = $state<string | null>(null);

	// New state for likes and comments
	let likeData = $state<LikeData | null>(null);
	let comments = $state<Comment[]>([]);
	let loadingComments = $state(false);
	let submittingComment = $state(false);
	let submittingLike = $state(false);
	
	// Get form data from SvelteKit form actions
	let { data, form }: { data: any, form: any } = $props();

	let id = $derived(page.params.id);

	onMount(async () => {
		try {
			loading = true;
			error = null;

			if (!id) {
				error = 'Philosopher ID is required';
				return;
			}

			// Load philosopher data
			const loadedPhilosopher = await loadPhilosopher(id);
			
			if (!loadedPhilosopher) {
				error = 'Philosopher not found';
				return;
			}
			
			philosopher = loadedPhilosopher;

			// Load likes and comments in parallel
			await Promise.all([
				loadLikes(),
				loadComments()
			]);

		} catch (err) {
			error = 'Failed to load philosopher details';
			console.error(err);
		} finally {
			loading = false;
		}
	});

	async function loadLikes() {
		if (!id) return;
		
		try {
			// Replace with your actual API call
			const response = await fetch(`/api/philosophers/${id}/likes`);
			if (response.ok) {
				likeData = await response.json();
			}
		} catch (err) {
			console.error('Failed to load likes:', err);
		}
	}

	async function loadComments() {
		if (!id) return;
		
		try {
			loadingComments = true;
			// Replace with your actual API call
			const response = await fetch(`/api/philosophers/${id}/comments`);
			if (response.ok) {
				comments = await response.json();
			}
		} catch (err) {
			console.error('Failed to load comments:', err);
		} finally {
			loadingComments = false;
		}
	}

	async function toggleLike() {
		if (!id || submittingLike) return;

		try {
			submittingLike = true;
			
			const response = await fetch(`/api/philosophers/${id}/like`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' }
			});

			if (response.ok) {
				likeData = await response.json();
			}
		} catch (err) {
			console.error('Failed to toggle like:', err);
		} finally {
			submittingLike = false;
		}
	}



	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}
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
			<div class="bg-white shadow-lg rounded-lg overflow-hidden mb-8">
				
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

					<!-- Like Section -->
					{#if likeData}
						<div class="mb-8 border-t pt-6">
							<div class="flex items-center gap-4">
								<button
									onclick={toggleLike}
									disabled={submittingLike}
									class="flex items-center gap-2 px-4 py-2 rounded-lg transition-all {likeData.user_has_liked 
										? 'bg-red-50 text-red-600 border border-red-200 hover:bg-red-100' 
										: 'bg-gray-50 text-gray-600 border border-gray-200 hover:bg-gray-100'} 
										disabled:opacity-50"
								>
									<svg class="w-5 h-5" fill={likeData.user_has_liked ? "currentColor" : "none"} stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
									</svg>
									{likeData.user_has_liked ? 'Liked' : 'Like'}
								</button>
								<span class="text-gray-600">
									{likeData.like_count} {likeData.like_count === 1 ? 'like' : 'likes'}
								</span>
							</div>
						</div>
					{/if}

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
								href="/add"
								class="flex-1 border border-gray-300 text-gray-700 text-center py-3 px-6 rounded-lg hover:bg-gray-50 transition-colors"
							>
								Add New Philosopher
							</a>
						</div>
					</div>
				</div>
			</div>

			<!-- Comments Section -->
			<div class="bg-white shadow-lg rounded-lg overflow-hidden">
				<div class="p-6 md:p-8">
					<h2 class="text-2xl font-bold text-gray-900 mb-6">
						Comments ({comments.length})
					</h2>

					<!-- Add Comment Form -->
					<div class="mb-8 bg-gray-50 rounded-lg p-6">
						<h3 class="text-lg font-semibold text-gray-900 mb-4">Share your thoughts</h3>
						
						<!-- SvelteKit 5 Form Actions -->
						<form method="POST" action="?/addComment" use:enhance>
							<input type="hidden" name="philosopher_id" value={id} />
							
							<div class="mb-4">
								<label for="user_name" class="block text-sm font-medium text-gray-700 mb-2">
									Your Name
								</label>
								<input
									type="text"
									name="user_name"
									placeholder="Enter your name"
									class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
									class:border-red-500={form?.missing}
									value={form?.user_name ?? ''}
									required
								/>
								{#if form?.missing}
									<p class="mt-1 text-red-500 text-sm">Name is required</p>
								{/if}
							</div>

							<div class="mb-4">
								<label for="comment_text" class="block text-sm font-medium text-gray-700 mb-2">
									Comment
								</label>
								<textarea
									name="comment_text"
									placeholder="What are your thoughts?"
									rows="4"
									class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 resize-none"
									class:border-red-500={form?.missing}
									value={form?.comment_text ?? ''}
									required
								></textarea>
								{#if form?.missing}
									<p class="mt-1 text-red-500 text-sm">Comment is required</p>
								{/if}
							</div>

							<button
								type="submit"
								class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
							>
								Post Comment
							</button>
						</form>
					</div>

					<!-- Success Message -->
					{#if form?.success}
						<div class="mb-6 p-4 bg-green-50 text-green-700 rounded-lg border border-green-200">
							{form.message}
						</div>
					{/if}

					<!-- Comments List -->
					{#if loadingComments}
						<div class="flex justify-center py-8">
							<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
						</div>
					{:else if comments.length === 0}
						<div class="text-center py-8 text-gray-500">
							<p>No comments yet. Be the first to share your thoughts!</p>
						</div>
					{:else}
						<div class="space-y-6">
							{#each comments as comment}
								<div class="border border-gray-200 rounded-lg p-4">
									<div class="flex items-start justify-between mb-3">
										<div class="flex items-center gap-3">
											<div class="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
												<span class="text-blue-600 font-medium text-sm">
													{comment.user_name.charAt(0).toUpperCase()}
												</span>
											</div>
											<div>
												<h4 class="font-medium text-gray-900">{comment.user_name}</h4>
												<time class="text-sm text-gray-500">
													{formatDate(comment.created_at)}
												</time>
											</div>
										</div>
									</div>
									<p class="text-gray-700 leading-relaxed pl-11">
										{comment.comment_text}
									</p>
								</div>
							{/each}
						</div>
					{/if}
				</div>
			</div>
		{/if}
	</div>
</div>