<script lang="ts">
    import type { Interaction } from '$lib/types';
    import { interactionApi } from '$lib/api/interactions';

    const {
        interactions: initialInteractions,
        philosopherId,
    } = $props<{
        interactions: Interaction[] | null;
        philosopherId: number;
    }>();

    let interactions = $state<Interaction[]>(initialInteractions ?? []);

    const comments = $derived(interactions.filter(i => i.type === 'comment'));
    const likes = $derived(interactions.filter(i => i.type === 'like'));
    const likeCount = $derived(likes.length);
    //const userLiked = $derived(likes.some(like => like.username === 'current_user')); // Replace with actual user check

    let newComment = $state('');
    let loading = $state(false);
    let error = $state<string | null>(null);

    const handleSubmit = async (e: Event) => {
        e.preventDefault();
        const trimmed = newComment.trim();
        if (!trimmed) return;

        loading = true;
        error = null;

        try {
            const response = await interactionApi.addComment(philosopherId, trimmed);

            const newInteraction: Interaction = {
                id: response.id,
                username: response.username || 'You', // Backend should provide this
                type: 'comment',
                content: trimmed,
                created_at: new Date().toISOString()
            };

            interactions = [newInteraction, ...interactions];
            newComment = '';

        } catch (err) {
            error = 'Failed to post comment. Please try again.';
            console.error('🚨 Comment error:', err);
        } finally {
            loading = false;
        }
    };

    const toggleLike = async () => {
        try {
            const response = await interactionApi.toggleLike(philosopherId);

            if (response.liked) {
                const newLike: Interaction = {
                    id: Date.now(), // Temporary ID, backend should provide real one
                    username: 'current_user', // Replace with actual username
                    type: 'like',
                    content: '',
                    created_at: new Date().toISOString()
                };
                interactions = [...interactions, newLike];
            } else {
                interactions = interactions.filter(i => 
                    !(i.type === 'like' && i.username === 'current_user')
                );
            }

        } catch (err) {
            error = 'Failed to update like. Please try again.';
            console.error('🚨 Like error:', err);
        }
    };

    // Sync with parent data if it changes
    $effect(() => {
        if (initialInteractions) {
            interactions = [...initialInteractions];
        }
    });
</script>

<div class="mt-8 border-t pt-6">
    <h3 class="text-xl font-bold text-gray-900 mb-4">
        Comments ({comments.length}) 
        <button 
            onclick={toggleLike}
            class={`ml-2 text-sm px-2 py-1 rounded-full transition-colors ${
                likeCount > 0 
                    ? 'bg-red-100 text-red-600 hover:bg-red-200' 
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
            }`}
        >
            ❤️ {likeCount}
        </button>
    </h3>

    <!-- Post Comment Form -->
    <form onsubmit={handleSubmit} class="mb-6">
        <textarea
            bind:value={newComment}
            placeholder="Share your thoughts about this philosopher..."
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            rows="3"
        ></textarea>
        {#if error}
            <p class="mt-2 text-sm text-red-600">{error}</p>
        {/if}
        <button
            type="submit"
            disabled={loading || !newComment.trim()}
            class="mt-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
            {loading ? 'Posting...' : 'Post Comment'}
        </button>
    </form>

    <!-- Comments List -->
    <div class="space-y-4">
        {#if comments.length === 0}
            <p class="text-gray-500 italic">No comments yet. Be the first to comment!</p>
        {:else}
            {#each comments as comment (comment.id)}
                <div class="bg-gray-50 p-4 rounded-lg border border-gray-100">
                    <div class="flex justify-between items-start">
                        <span class="text-sm font-medium text-gray-700">{comment.username}</span>
                        <span class="text-sm text-gray-500">
                            {new Date(comment.created_at).toLocaleDateString()}
                        </span>
                    </div>
                    <p class="mt-2 text-gray-700">{comment.content}</p>
                </div>
            {/each}
        {/if}
    </div>
</div>