import { fail } from '@sveltejs/kit';
//import type { Actions } from './$types';
import { PUBLIC_GO_API_BASE } from '$env/static/public';
import type { RequestEvent, Actions } from '@sveltejs/kit';


const API_BASE = PUBLIC_GO_API_BASE || 'http://localhost:8080';

export const actions = {
	addComment: async ({ request, params }) => {
		const data = await request.formData();
		const user_name = data.get('user_name');
		const comment_text = data.get('comment_text');
		const philosopher_id = params.id;

		// Check if philosopher_id exists
		if (!philosopher_id) {
			return fail(400, { user_name, comment_text, error: 'Philosopher ID is required' });
		}

		if (!user_name || !comment_text) {
			return fail(400, { user_name, comment_text, missing: true });
		}

		try {
			// Call your Go API here
			const response = await fetch(`${API_BASE}/api/philosophers/${philosopher_id}/comments`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					user_name: user_name.toString(),
					comment_text: comment_text.toString(),
					philosopher_id: parseInt(philosopher_id)
				})
			});

			if (!response.ok) {
				throw new Error('Failed to post comment');
			}

			return {
				success: true,
				message: 'Comment posted successfully!'
			};

		} catch (error) {
			return fail(500, {
				user_name,
				comment_text,
				error: 'Failed to post comment'
			});
		}
	}
} satisfies Actions;