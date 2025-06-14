// @ts-nocheck
import { serverApi } from '$lib/server-api.js';
import { error, redirect } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';

export const load = async ({ params }: Parameters<PageServerLoad>[0]) => {
	try {
		const reminderId = parseInt(params.reminder_id);
		
		if (isNaN(reminderId)) {
			throw error(400, 'Invalid reminder ID');
		}
		
		console.log('Server: Loading reminder', reminderId, 'for editing...');
		const reminder = await serverApi.getReminder(reminderId);
		console.log('Server: Loaded reminder for editing:', reminder.title);
		
		return {
			reminder
		};
	} catch (err) {
		console.error('Server: Error loading reminder for editing:', err);
		throw error(500, {
			message: 'Failed to load reminder',
			details: err instanceof Error ? err.message : 'Unknown error'
		});
	}
};

export const actions = {
	default: async ({ request, params }: import('./$types').RequestEvent) => {
		// Initialize variables to ensure they're available in catch block
		let title = '';
		let description = '';
		let reminderType: 'one-time' | 'persistent' = 'one-time';
		let scheduledAt = '';
		let scheduledTime = '09:00';
		let deliveryWindowMinutes = 15;
		let deliveryMethod = 'email';
		let deliveryAddress = '';
		let isActive = true;
		let notificationChannels: string[] = [];
		let scheduledDaysOfWeek: number[] = [];

		try {
			const reminderId = parseInt(params.reminder_id);
			
			if (isNaN(reminderId)) {
				throw error(400, 'Invalid reminder ID');
			}
			
			const data = await request.formData();
			
			// Extract form data
			title = data.get('title') as string;
			description = data.get('description') as string;
			reminderType = data.get('reminder_type') as 'one-time' | 'persistent';
			scheduledAt = data.get('scheduled_at') as string;
			scheduledTime = data.get('scheduled_time') as string;
			deliveryWindowMinutes = parseInt(data.get('delivery_window_minutes') as string) || 15;
			deliveryMethod = data.get('delivery_method') as string;
			deliveryAddress = data.get('delivery_address') as string;
			isActive = data.get('is_active') === 'on';
			
			// Parse notification channels
			notificationChannels = [];
			if (data.get('notification_channels_email')) notificationChannels.push('email');
			if (data.get('notification_channels_sms')) notificationChannels.push('sms');
			if (data.get('notification_channels_call')) notificationChannels.push('call');
			
			// Parse scheduled days of week for persistent reminders
			scheduledDaysOfWeek = [];
			if (reminderType === 'persistent') {
				for (let i = 0; i < 7; i++) {
					if (data.get(`scheduled_days_of_week_${i}`)) {
						scheduledDaysOfWeek.push(i);
					}
				}
			}
			
			// Prepare form data object for potential error return
			const formData = {
				title: title || '',
				description: description || '',
				reminder_type: reminderType,
				scheduled_at: scheduledAt || '',
				scheduled_time: scheduledTime || '09:00',
				scheduled_days_of_week: scheduledDaysOfWeek,
				notification_channels: notificationChannels,
				delivery_window_minutes: deliveryWindowMinutes || 15,
				delivery_method: deliveryMethod || 'email',
				delivery_address: deliveryAddress || '',
				is_active: isActive
			};
			
			// Validation
			if (!title?.trim()) {
				return { error: 'Title is required', data: formData };
			}
			
			if (!scheduledAt && reminderType === 'one-time') {
				return { error: 'Scheduled date/time is required for one-time reminders', data: formData };
			}
			
			if (reminderType === 'persistent' && scheduledDaysOfWeek.length === 0) {
				return { error: 'Please select at least one day of the week for persistent reminders', data: formData };
			}
			
			if (notificationChannels.length === 0) {
				return { error: 'Please select at least one notification channel', data: formData };
			}
			
			if (!deliveryAddress?.trim()) {
				return { error: 'Delivery address is required', data: formData };
			}
			
			// Prepare the data for submission
			let finalScheduledAt = scheduledAt;
			
			// For persistent reminders, set scheduled_at if needed
			if (reminderType === 'persistent' && scheduledTime) {
				const today = new Date();
				const [hours, minutes] = scheduledTime.split(':');
				today.setHours(parseInt(hours), parseInt(minutes), 0, 0);
				finalScheduledAt = today.toISOString();
			}
			
			const updateData = {
				title: title.trim(),
				description: description?.trim() || undefined,
				reminder_type: reminderType,
				scheduled_at: finalScheduledAt,
				scheduled_time: reminderType === 'persistent' ? scheduledTime : '00:00',
				scheduled_days_of_week: scheduledDaysOfWeek.length > 0 ? scheduledDaysOfWeek : null,
				notification_channels: notificationChannels,
				delivery_window_minutes: deliveryWindowMinutes,
				delivery_method: deliveryMethod,
				delivery_address: deliveryAddress.trim(),
				is_active: isActive
			};
			
			console.log('Server: Updating reminder', reminderId, 'with data:', updateData);
			await serverApi.updateReminder(reminderId, updateData);
			console.log('Server: Successfully updated reminder', reminderId);
			
			throw redirect(303, '/reminders');
		} catch (err) {
			if (err instanceof Response) {
				// This is a redirect, re-throw it
				throw err;
			}
			
			console.error('Server: Error updating reminder:', err);
			
			// Return form data so user input is preserved
			const formData = {
				title: title || '',
				description: description || '',
				reminder_type: reminderType,
				scheduled_at: scheduledAt || '',
				scheduled_time: scheduledTime || '09:00',
				scheduled_days_of_week: scheduledDaysOfWeek,
				notification_channels: notificationChannels,
				delivery_window_minutes: deliveryWindowMinutes || 15,
				delivery_method: deliveryMethod || 'email',
				delivery_address: deliveryAddress || '',
				is_active: isActive
			};
			
			return { 
				error: err instanceof Error ? err.message : 'Failed to update reminder',
				data: formData
			};
		}
	}
};;null as any as Actions;