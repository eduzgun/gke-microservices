<script lang="ts">
	import { goto } from '$app/navigation';
  	import { philosopherApi } from '$lib/api/philosophers';
	import type { Philosopher } from '$lib/types';
	import "../../../app.css";

	let name = $state('');
	let dateBorn = $state('');
	let dateDied = $state('');
	let birthplace = $state('');
	let bio = $state('');
	let portraitUri = $state('');
	let interests = $state('');
	let loading = $state(false);
	let error = $state<string | null>(null);
	let success = $state(false);

	const handleSubmit = async (e: Event) => {
		e.preventDefault();
		if (!name.trim() || !bio.trim()) {
			error = 'Name and bio are required';
			return;
		}

		try {
			loading = true;
			error = null;

			// Parse interests from comma-separated string
			const interestsArray = interests
				.split(',')
				.map(interest => interest.trim())
				.filter(interest => interest.length > 0);

			const newPhilosopher: Omit<Philosopher, 'id' | 'created_at'> = {
				name: name.trim(),
				date_born: dateBorn.trim(),
				date_died: dateDied.trim(),
				birthplace: birthplace.trim(),
				bio: bio.trim(),
				portrait_uri: portraitUri.trim(),
				interests: interestsArray
			};

			// Call your actual API to save the philosopher
			const result = await philosopherApi.create(newPhilosopher);
			
			if (!result) {
				error = 'Failed to save philosopher';
				return;
			}
			
			console.log('Philosopher created with ID:', result.id);
			success = true;
			
			// Clear form after successful submission
			setTimeout(() => {
				resetForm();
				goto('/'); // Redirect to main page
			}, 2000);

		} catch (err) {
			error = 'Failed to save philosopher';
			console.error(err);
		} finally {
			loading = false;
		}
	}

	function resetForm() {
		name = '';
		dateBorn = '';
		dateDied = '';
		birthplace = '';
		bio = '';
		portraitUri = '';
		interests = '';
		success = false;
		error = null;
	}
</script>

<div class="min-h-screen bg-gray-50 py-8">
	<div class="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8">
		<div class="text-center mb-8">
			<h1 class="text-3xl font-bold text-gray-900 mb-4">Add New Philosopher</h1>
			<p class="text-gray-600">
				Share the wisdom of a great thinker with the world
			</p>
		</div>

		<div class="bg-white shadow-md rounded-lg p-6">
			{#if success}
				<div class="bg-green-50 border border-green-200 rounded-lg p-4 mb-6">
					<div class="flex items-center">
						<div class="flex-shrink-0">
							<svg class="h-5 w-5 text-green-400" viewBox="0 0 20 20" fill="currentColor">
								<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
							</svg>
						</div>
						<div class="ml-3">
							<h3 class="text-sm font-medium text-green-800">
								Philosopher added successfully!
							</h3>
							<p class="text-sm text-green-700 mt-1">
								Redirecting to main page...
							</p>
						</div>
					</div>
				</div>
			{/if}

			{#if error}
				<div class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
					<h3 class="text-sm font-medium text-red-800">{error}</h3>
				</div>
			{/if}

			<form onsubmit={handleSubmit} class="space-y-6">
				<!-- Name -->
				<div>
					<label for="name" class="block text-sm font-medium text-gray-700 mb-2">
						Name *
					</label>
					<input
						type="text"
						id="name"
						bind:value={name}
						class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
						placeholder="e.g., Socrates"
						required
					/>
				</div>

				<!-- Date Born -->
				<div>
					<label for="dateBorn" class="block text-sm font-medium text-gray-700 mb-2">
						Date Born
					</label>
					<input
						type="text"
						id="dateBorn"
						bind:value={dateBorn}
						class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
						placeholder="e.g., 470 BCE"
					/>
				</div>

				<!-- Date Died -->
				<div>
					<label for="dateDied" class="block text-sm font-medium text-gray-700 mb-2">
						Date Died
					</label>
					<input
						type="text"
						id="dateDied"
						bind:value={dateDied}
						class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
						placeholder="e.g., 399 BCE"
					/>
				</div>

				<!-- Birthplace -->
				<div>
					<label for="birthplace" class="block text-sm font-medium text-gray-700 mb-2">
						Birthplace
					</label>
					<input
						type="text"
						id="birthplace"
						bind:value={birthplace}
						class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
						placeholder="e.g., Athens, Greece"
					/>
				</div>

				<!-- Portrait URI -->
				<div>
					<label for="portraitUri" class="block text-sm font-medium text-gray-700 mb-2">
						Portrait URL
					</label>
					<input
						type="url"
						id="portraitUri"
						bind:value={portraitUri}
						class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
						placeholder="https://example.com/portrait.jpg"
					/>
				</div>

				<!-- Bio -->
				<div>
					<label for="bio" class="block text-sm font-medium text-gray-700 mb-2">
						Biography *
					</label>
					<textarea
						id="bio"
						bind:value={bio}
						rows="4"
						class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
						placeholder="Write a brief biography..."
						required
					></textarea>
				</div>

				<!-- Interests -->
				<div>
					<label for="interests" class="block text-sm font-medium text-gray-700 mb-2">
						Interests
					</label>
					<input
						type="text"
						id="interests"
						bind:value={interests}
						class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
						placeholder="Ethics, Logic, Metaphysics (comma-separated)"
					/>
					<p class="mt-1 text-sm text-gray-500">
						Separate multiple interests with commas
					</p>
				</div>

				<!-- Action Buttons -->
				<div class="flex gap-4 pt-6">
					<button
						type="submit"
						disabled={loading}
						class="flex-1 bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
					>
						{#if loading}
							<span class="flex items-center justify-center">
								<svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
									<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
									<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
								</svg>
								Adding...
							</span>
						{:else}
							Add Philosopher
						{/if}
					</button>
					
					<button
						type="button"
						onclick={resetForm}
						class="px-4 py-2 border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
					>
						Reset
					</button>

					<a
						href="/"
						class="px-4 py-2 border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 text-center"
					>
						Cancel
					</a>
				</div>
			</form>
		</div>
	</div>
</div>