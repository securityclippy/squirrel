import type { PageServerLoad } from './$types';
import { serverApi, type UserStatistics } from '$lib/server-api';

export const load: PageServerLoad = async () => {
	try {
		// For now, use a default user ID since auth is disabled
		// TODO: Get actual user ID from auth context when re-enabled
		const statistics = await serverApi.getUserStatistics(1);
		
		return {
			statistics,
			error: null
		};
	} catch (error) {
		console.error('Failed to fetch statistics:', error);
		
		// Return fallback data so the page still loads
		const fallbackStats: UserStatistics = {
			overview: {
				total_reminders: 0,
				completed_reminders: 0,
				active_reminders: 0,
				member_since: new Date().toISOString()
			},
			weekly: {
				this_week: 0,
				last_week: 0,
				change_percent: 0
			},
			monthly: {
				this_month: 0,
				last_month: 0,
				change_percent: 0
			},
			categories: []
		};

		return {
			statistics: fallbackStats,
			error: error instanceof Error ? error.message : 'Failed to load statistics'
		};
	}
};