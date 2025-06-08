import { serverApi } from '$lib/server-api.js';
import { error, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
	default: async ({ request }) => {
		// Initialize variables to ensure they're available in catch block
		let title = '';
		let description = '';
		let reminderType: 'one-time' | 'recurring' = 'one-time';
		let scheduledAt = '';
		let scheduledTime = '09:00';
		let deliveryWindowMinutes = 15;
		let deliveryMethod = 'email';
		let deliveryAddress = '';
		let notificationChannels: string[] = [];
		let scheduledDaysOfWeek: number[] = [];
		let isPersistent = false;
		let reminderIntervalMinutes: number | undefined;

		try {
			const data = await request.formData();
			
			// Extract form data
			title = data.get('title') as string;
			description = data.get('description') as string;
			reminderType = data.get('reminder_type') as 'one-time' | 'recurring';
			scheduledAt = data.get('scheduled_at') as string;
			scheduledTime = data.get('scheduled_time') as string;
			deliveryWindowMinutes = parseInt(data.get('delivery_window_minutes') as string) || 15;
			deliveryMethod = data.get('delivery_method') as string;
			deliveryAddress = data.get('delivery_address') as string;
			isPersistent = data.get('is_persistent') === 'on';
			const intervalStr = data.get('reminder_interval_minutes') as string;
			reminderIntervalMinutes = intervalStr ? parseInt(intervalStr) : undefined;
			
			// Parse notification channels
			notificationChannels = [];
			if (data.get('notification_channels_email')) notificationChannels.push('email');
			if (data.get('notification_channels_sms')) notificationChannels.push('sms');
			if (data.get('notification_channels_call')) notificationChannels.push('call');
			
			// Parse scheduled days of week for recurring reminders
			scheduledDaysOfWeek = [];
			if (reminderType === 'recurring') {
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
				is_persistent: isPersistent,
				reminder_interval_minutes: reminderIntervalMinutes || 30
			};
			
			// Validation
			if (!title?.trim()) {
				return { error: 'Title is required', data: formData };
			}
			
			if (!scheduledAt && reminderType === 'one-time') {
				return { error: 'Scheduled date/time is required for one-time reminders', data: formData };
			}
			
			if (reminderType === 'recurring' && scheduledDaysOfWeek.length === 0) {
				return { error: 'Please select at least one day of the week for recurring reminders', data: formData };
			}
			
			if (notificationChannels.length === 0) {
				return { error: 'Please select at least one notification channel', data: formData };
			}
			
			if (!deliveryAddress?.trim()) {
				return { error: 'Delivery address is required', data: formData };
			}
			
			// Prepare the data for submission
			let finalScheduledAt = scheduledAt;
			
			// For recurring reminders, set scheduled_at to today + scheduled_time
			if (reminderType === 'recurring') {
				const today = new Date();
				const [hours, minutes] = scheduledTime.split(':');
				today.setHours(parseInt(hours), parseInt(minutes), 0, 0);
				finalScheduledAt = today.toISOString();
			} else if (scheduledAt) {
				// Ensure one-time reminders have proper ISO string format
				finalScheduledAt = new Date(scheduledAt).toISOString();
			}
			
			const reminderData = {
				title: title.trim(),
				description: description?.trim() || undefined,
				reminder_type: reminderType,
				scheduled_at: finalScheduledAt,
				scheduled_time: reminderType === 'recurring' ? scheduledTime : '00:00',
				scheduled_days_of_week: scheduledDaysOfWeek.length > 0 ? scheduledDaysOfWeek : null,
				notification_channels: notificationChannels,
				delivery_window_minutes: deliveryWindowMinutes,
				delivery_method: deliveryMethod,
				delivery_address: deliveryAddress.trim(),
				is_persistent: isPersistent,
				reminder_interval_minutes: isPersistent ? reminderIntervalMinutes : undefined
			};
			
			console.log('Server: Creating reminder with data:', JSON.stringify(reminderData, null, 2));
			console.log('Server: Form data extracted - title:', title, 'reminderType:', reminderType, 'scheduledAt:', scheduledAt);
			console.log('Server: Notification channels:', notificationChannels);
			console.log('Server: Scheduled days of week:', scheduledDaysOfWeek);
			
			await serverApi.createReminder(reminderData);
			console.log('Server: Successfully created reminder');
			
			redirect(303, '/reminders');
		} catch (err) {
			
			console.error('Server: Error creating reminder:', err);
			
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
				is_persistent: isPersistent,
				reminder_interval_minutes: reminderIntervalMinutes || 30
			};
			
			return { 
				error: err instanceof Error ? err.message : 'Failed to create reminder',
				data: formData
			};
		}
	}
};