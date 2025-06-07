<script>
	import { enhance } from '$app/forms';
	
	export let data;
	export let form;
	
	let loading = false;
	
	// Use form data if available (from validation errors), otherwise use loaded reminder data
	$: formData = form?.data || {
		title: data.reminder.title,
		description: data.reminder.description || '',
		reminder_type: data.reminder.reminder_type,
		scheduled_at: data.reminder.reminder_type === 'one-time' 
			? new Date(data.reminder.scheduled_at).toISOString().slice(0, 16)
			: '',
		scheduled_time: data.reminder.scheduled_time.slice(0, 5),
		scheduled_days_of_week: data.reminder.scheduled_days_of_week || [],
		notification_channels: data.reminder.notification_channels,
		delivery_window_minutes: data.reminder.delivery_window_minutes,
		delivery_method: data.reminder.delivery_method,
		delivery_address: data.reminder.delivery_address,
		is_active: data.reminder.is_active
	};
	
	// Extract reminder type for reactivity
	$: reminderType = formData.reminder_type;
	
	// Display server-side validation errors
	$: error = form?.error;

	const daysOfWeek = [
		{ value: 0, label: 'Sunday' },
		{ value: 1, label: 'Monday' },
		{ value: 2, label: 'Tuesday' },
		{ value: 3, label: 'Wednesday' },
		{ value: 4, label: 'Thursday' },
		{ value: 5, label: 'Friday' },
		{ value: 6, label: 'Saturday' }
	];

	function handleFormSubmit() {
		loading = true;
		return async ({ update }) => {
			loading = false;
			await update();
		};
	}

	// Reactive field visibility based on reminder type
	$: isOneTime = reminderType === 'one-time';
	
	// Helper function to check if notification channel is selected
	function isChannelSelected(channel) {
		return formData.notification_channels.includes(channel);
	}
	
	// Helper function to check if day is selected
	function isDaySelected(day) {
		return formData.scheduled_days_of_week.includes(day);
	}
</script>

<svelte:head>
	<title>Edit Reminder - Squirrel Reminders</title>
</svelte:head>

<div class="new-reminder-page">
	<div class="header">
		<a href="/reminders" class="back-button">← Back to Reminders</a>
		<h1>Edit Reminder</h1>
	</div>

	<form method="POST" use:enhance={handleFormSubmit} class="reminder-form">
		{#if error}
			<div class="error-message">
				{error}
			</div>
		{/if}

		<div class="form-section">
			<h2>Basic Information</h2>
			
			<div class="form-group">
				<label for="title">Title *</label>
				<input
					id="title"
					name="title"
					type="text"
					value={formData.title}
					placeholder="What should I remind you about?"
					maxlength="500"
					required
				/>
			</div>

			<div class="form-group">
				<label for="description">Description</label>
				<textarea
					id="description"
					name="description"
					value={formData.description}
					placeholder="Additional details about this reminder"
					rows="3"
				></textarea>
			</div>

			<div class="form-group">
				<label for="reminder_type">Reminder Type *</label>
				<select id="reminder_type" name="reminder_type" value={formData.reminder_type}>
					<option value="one-time">One-time</option>
					<option value="persistent">Recurring</option>
				</select>
			</div>
		</div>

		<div class="form-section">
			<h2>Scheduling</h2>

			{#if isOneTime}
				<div class="form-group">
					<label for="scheduled_at">Date & Time *</label>
					<input
						id="scheduled_at"
						name="scheduled_at"
						type="datetime-local"
						value={formData.scheduled_at}
					/>
				</div>
				<!-- Hidden input to provide default scheduled_time for one-time reminders -->
				<input type="hidden" name="scheduled_time" value="00:00" />
			{:else}
				<div class="form-group">
					<label for="scheduled_time">Time *</label>
					<input
						id="scheduled_time"
						name="scheduled_time"
						type="time"
						value={formData.scheduled_time}
					/>
				</div>

				<div class="form-group">
					<label>Days of the Week *</label>
					<div class="days-grid">
						{#each daysOfWeek as day}
							<label class="day-checkbox">
								<input
									type="checkbox"
									name="scheduled_days_of_week_{day.value}"
									value="{day.value}"
									checked={isDaySelected(day.value)}
								/>
								<span class="checkmark"></span>
								{day.label}
							</label>
						{/each}
					</div>
				</div>
			{/if}

			<div class="form-group">
				<label for="delivery_window">Delivery Window (±minutes)</label>
				<input
					id="delivery_window"
					name="delivery_window_minutes"
					type="number"
					value={formData.delivery_window_minutes}
					min="1"
					max="1440"
					step="1"
				/>
				<small>How flexible should the delivery timing be? (1-1440 minutes)</small>
			</div>
		</div>

		<div class="form-section">
			<h2>Notification Settings</h2>

			<div class="form-group">
				<label>Notification Channels *</label>
				<div class="channels-grid">
					<label class="channel-checkbox">
						<input
							type="checkbox"
							name="notification_channels_email"
							value="email"
							checked={isChannelSelected('email')}
						/>
						<span class="checkmark"></span>
						📧 Email
					</label>
					<label class="channel-checkbox">
						<input
							type="checkbox"
							name="notification_channels_sms"
							value="sms"
							checked={isChannelSelected('sms')}
						/>
						<span class="checkmark"></span>
						📱 SMS
					</label>
					<label class="channel-checkbox">
						<input
							type="checkbox"
							name="notification_channels_call"
							value="call"
							checked={isChannelSelected('call')}
						/>
						<span class="checkmark"></span>
						📞 Phone Call
					</label>
				</div>
			</div>

			<div class="form-group">
				<label for="delivery_method">Primary Delivery Method *</label>
				<select id="delivery_method" name="delivery_method" value={formData.delivery_method}>
					<option value="email">Email</option>
					<option value="sms">SMS</option>
					<option value="call">Phone Call</option>
				</select>
			</div>

			<div class="form-group">
				<label for="delivery_address">
					{#if formData.delivery_method === 'email'}
						Email Address *
					{:else if formData.delivery_method === 'sms'}
						Phone Number *
					{:else}
						Phone Number *
					{/if}
				</label>
				<input
					id="delivery_address"
					name="delivery_address"
					type={formData.delivery_method === 'email' ? 'email' : 'tel'}
					value={formData.delivery_address}
					placeholder={formData.delivery_method === 'email' ? 'you@example.com' : '+1234567890'}
					required
				/>
			</div>
		</div>

		<div class="form-section">
			<h2>Status</h2>
			
			<div class="form-group">
				<label class="checkbox-label">
					<input
						type="checkbox"
						name="is_active"
						checked={formData.is_active}
					/>
					<span class="checkmark"></span>
					Active (reminder will be sent)
				</label>
			</div>
		</div>

		<div class="form-actions">
			<a href="/reminders" class="cancel-button">Cancel</a>
			<button 
				type="submit" 
				class="submit-button"
				disabled={loading}
			>
				{#if loading}
					<span class="spinner"></span>
					Updating...
				{:else}
					Update Reminder
				{/if}
			</button>
		</div>
	</form>
</div>

<style>
	.new-reminder-page {
		max-width: 800px;
		margin: 0 auto;
		padding: 2rem;
	}

	.header {
		margin-bottom: 2rem;
	}

	.back-button {
		display: inline-flex;
		align-items: center;
		color: #059669;
		text-decoration: none;
		margin-bottom: 1rem;
		font-weight: 500;
		transition: color 0.2s;
	}

	.back-button:hover {
		color: #047857;
	}

	.header h1 {
		color: #059669;
		font-size: 2rem;
		margin: 0;
	}

	.reminder-form {
		background: white;
		border-radius: 12px;
		box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
		overflow: hidden;
	}

	.error-message {
		background: #fef2f2;
		border: 1px solid #fecaca;
		color: #dc2626;
		padding: 1rem;
		margin-bottom: 1.5rem;
		border-radius: 8px;
	}

	.form-section {
		padding: 2rem;
		border-bottom: 1px solid #e5e7eb;
	}

	.form-section:last-child {
		border-bottom: none;
	}

	.form-section h2 {
		color: #374151;
		font-size: 1.25rem;
		margin: 0 0 1.5rem 0;
		font-weight: 600;
	}

	.form-group {
		margin-bottom: 1.5rem;
	}

	.form-group:last-child {
		margin-bottom: 0;
	}

	label {
		display: block;
		font-weight: 600;
		color: #374151;
		margin-bottom: 0.5rem;
	}

	input[type="text"],
	input[type="email"],
	input[type="tel"],
	input[type="datetime-local"],
	input[type="time"],
	input[type="number"],
	textarea,
	select {
		width: 100%;
		padding: 0.75rem;
		border: 1px solid #d1d5db;
		border-radius: 8px;
		font-size: 1rem;
		transition: border-color 0.2s, box-shadow 0.2s;
	}

	input:focus,
	textarea:focus,
	select:focus {
		outline: none;
		border-color: #4ade80;
		box-shadow: 0 0 0 3px rgba(74, 222, 128, 0.1);
	}

	textarea {
		resize: vertical;
		min-height: 80px;
	}

	small {
		display: block;
		color: #6b7280;
		font-size: 0.875rem;
		margin-top: 0.25rem;
	}

	.days-grid,
	.channels-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
		gap: 0.75rem;
		margin-top: 0.5rem;
	}

	.day-checkbox,
	.channel-checkbox {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-weight: normal;
		cursor: pointer;
		padding: 0.75rem;
		border: 1px solid #e5e7eb;
		border-radius: 8px;
		transition: background-color 0.2s, border-color 0.2s;
	}

	.day-checkbox:hover,
	.channel-checkbox:hover {
		background: #f9fafb;
		border-color: #d1d5db;
	}

	.day-checkbox input,
	.channel-checkbox input {
		width: auto;
		margin: 0;
	}

	.checkmark {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 20px;
		height: 20px;
		border: 2px solid #d1d5db;
		border-radius: 4px;
		transition: all 0.2s;
	}

	input[type="checkbox"]:checked + .checkmark {
		background: #4ade80;
		border-color: #4ade80;
	}

	input[type="checkbox"]:checked + .checkmark::after {
		content: "✓";
		color: white;
		font-size: 0.75rem;
		font-weight: bold;
	}

	input[type="checkbox"] {
		display: none;
	}

	.checkbox-label {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		font-weight: normal;
		cursor: pointer;
		padding: 0.75rem;
		border: 1px solid #e5e7eb;
		border-radius: 8px;
		transition: background-color 0.2s, border-color 0.2s;
	}

	.checkbox-label:hover {
		background: #f9fafb;
		border-color: #d1d5db;
	}

	.form-actions {
		padding: 2rem;
		display: flex;
		gap: 1rem;
		justify-content: flex-end;
		background: #f9fafb;
	}

	.cancel-button {
		padding: 0.75rem 1.5rem;
		border: 1px solid #d1d5db;
		border-radius: 8px;
		color: #374151;
		text-decoration: none;
		font-weight: 500;
		transition: background-color 0.2s;
		display: flex;
		align-items: center;
	}

	.cancel-button:hover {
		background: #f3f4f6;
	}

	.submit-button {
		background: linear-gradient(135deg, #4ade80, #22d3ee);
		color: white;
		border: none;
		padding: 0.75rem 2rem;
		border-radius: 8px;
		font-weight: 600;
		cursor: pointer;
		transition: transform 0.2s, box-shadow 0.2s;
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.submit-button:hover:not(:disabled) {
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(0,0,0,0.15);
	}

	.submit-button:disabled {
		opacity: 0.7;
		cursor: not-allowed;
	}

	.spinner {
		width: 16px;
		height: 16px;
		border: 2px solid rgba(255,255,255,0.3);
		border-left: 2px solid white;
		border-radius: 50%;
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		0% { transform: rotate(0deg); }
		100% { transform: rotate(360deg); }
	}

	@media (max-width: 640px) {
		.new-reminder-page {
			padding: 1rem;
		}

		.form-section {
			padding: 1.5rem;
		}

		.form-actions {
			flex-direction: column;
			padding: 1.5rem;
		}

		.days-grid,
		.channels-grid {
			grid-template-columns: 1fr;
		}
	}
</style>