import { error } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';

export const load: PageServerLoad = async ({ cookies, request }) => {
	try {
		console.log('Server: Loading reminders...');
		
		// For now, make a direct call to the backend without authentication
		// since we disabled auth middleware temporarily
		const backendUrl = 'http://backend:8080/api/reminders';
		console.log('Server: Fetching from', backendUrl);
		
		const response = await fetch(backendUrl, {
			headers: {
				'Content-Type': 'application/json',
			}
		});
		
		if (!response.ok) {
			throw new Error(`Backend responded with ${response.status}: ${response.statusText}`);
		}
		
		const data = await response.json();
		console.log('Server: Loaded', data.count, 'reminders');
		
		return {
			reminders: data.reminders,
			count: data.count
		};
	} catch (err) {
		console.error('Server: Error loading reminders:', err);
		throw error(500, {
			message: 'Failed to load reminders',
			details: err instanceof Error ? err.message : 'Unknown error'
		});
	}
};

export const actions: Actions = {
	delete: async ({ request }) => {
		try {
			const data = await request.formData();
			const id = data.get('id');
			
			if (!id || typeof id !== 'string') {
				throw error(400, 'Invalid reminder ID');
			}
			
			console.log('Server: Deleting reminder', id);
			await serverApi.deleteReminder(parseInt(id));
			console.log('Server: Successfully deleted reminder', id);
			
			return { success: true };
		} catch (err) {
			console.error('Server: Error deleting reminder:', err);
			throw error(500, {
				message: 'Failed to delete reminder',
				details: err instanceof Error ? err.message : 'Unknown error'
			});
		}
	}
};