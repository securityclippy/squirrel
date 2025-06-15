<script>
	import { enhance } from '$app/forms';
	import { invalidateAll } from '$app/navigation';
	import AuthGuard from '$lib/components/AuthGuard.svelte';
	
	export let data;
	
	// Extract data from server-side load
	$: ({ reminders, count } = data);

	const formatDateTime = (dateString) => {
		return new Date(dateString).toLocaleString();
	};

	const formatTime = (timeString) => {
		return timeString.slice(0, 5); // HH:MM format
	};

	const getDayNames = (daysArray) => {
		if (!daysArray || daysArray.length === 0) return 'No specific days';
		const dayNames = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
		return daysArray.map(day => dayNames[day]).join(', ');
	};

	const getStatusColor = (status) => {
		switch (status) {
			case 'pending': return 'bg-yellow-100 text-yellow-800';
			case 'completed': return 'bg-green-100 text-green-800';
			case 'failed': return 'bg-red-100 text-red-800';
			default: return 'bg-gray-100 text-gray-800';
		}
	};

	const getReminderTypeColor = (type) => {
		return type === 'persistent' 
			? 'bg-blue-100 text-blue-800' 
			: 'bg-purple-100 text-purple-800';
	};

	function handleDeleteSubmit(event) {
		if (!confirm('Are you sure you want to delete this reminder?')) {
			event.preventDefault();
			return false;
		}
		return true;
	}
</script>

<svelte:head>
	<title>Reminders - Squirrel Reminders</title>
</svelte:head>

<AuthGuard>
	<div class="reminders-page">
		<div class="header">
			<h1>My Reminders</h1>
			<a href="/reminders/new" class="add-button">
				<span class="icon">➕</span>
				Add New Reminder
			</a>
		</div>

	{#if reminders.length === 0}
		<div class="empty-state">
			<div class="empty-icon">📝</div>
			<h3>No reminders yet</h3>
			<p>Create your first reminder to get started!</p>
			<a href="/reminders/new" class="primary-button">Create Reminder</a>
		</div>
	{:else}
		<div class="reminders-grid">
			{#each reminders as reminder (reminder.id)}
				<div class="reminder-card">
					<div class="card-header">
						<h3 class="reminder-title">{reminder.title}</h3>
						<div class="card-actions">
							<a 
								href="/reminders/{reminder.id}/edit"
								class="edit-button"
								title="Edit reminder"
							>
								✏️
							</a>
							<form method="POST" action="?/delete" use:enhance on:submit={handleDeleteSubmit} style="display: inline;">
								<input type="hidden" name="id" value={reminder.id} />
								<button 
									type="submit"
									class="delete-button"
									title="Delete reminder"
								>
									🗑️
								</button>
							</form>
						</div>
					</div>

					{#if reminder.description}
						<p class="reminder-description">{reminder.description}</p>
					{/if}

					<div class="reminder-details">
						<div class="detail-row">
							<span class="label">Type:</span>
							<span class="badge {getReminderTypeColor(reminder.reminder_type)}">
								{reminder.reminder_type}
							</span>
						</div>

						<div class="detail-row">
							<span class="label">Status:</span>
							<span class="badge {getStatusColor(reminder.status)}">
								{reminder.status}
							</span>
						</div>

						<div class="detail-row">
							<span class="label">Scheduled:</span>
							<span>{formatDateTime(reminder.scheduled_at)}</span>
						</div>

						{#if reminder.reminder_type === 'persistent'}
							<div class="detail-row">
								<span class="label">Time:</span>
								<span>{formatTime(reminder.scheduled_time)}</span>
							</div>

							<div class="detail-row">
								<span class="label">Days:</span>
								<span>{getDayNames(reminder.scheduled_days_of_week)}</span>
							</div>
						{/if}

						<div class="detail-row">
							<span class="label">Channels:</span>
							<div class="channels">
								{#each reminder.notification_channels as channel}
									<span class="channel-badge">{channel}</span>
								{/each}
							</div>
						</div>

						<div class="detail-row">
							<span class="label">Window:</span>
							<span>±{reminder.delivery_window_minutes} minutes</span>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
	</div>
</AuthGuard>

<style>
	.reminders-page {
		max-width: 1200px;
		margin: 0 auto;
		padding: 2rem;
	}

	.header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 2rem;
		flex-wrap: wrap;
		gap: 1rem;
	}

	.header h1 {
		color: #059669;
		font-size: 2rem;
		margin: 0;
	}

	.add-button {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		background: linear-gradient(135deg, #4ade80, #22d3ee);
		color: white;
		padding: 0.75rem 1.5rem;
		border-radius: 8px;
		text-decoration: none;
		font-weight: 600;
		transition: transform 0.2s, box-shadow 0.2s;
	}

	.add-button:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(0,0,0,0.15);
	}

	.loading {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 3rem;
		color: #666;
	}

	.spinner {
		width: 40px;
		height: 40px;
		border: 4px solid #e5e7eb;
		border-left: 4px solid #4ade80;
		border-radius: 50%;
		animation: spin 1s linear infinite;
		margin-bottom: 1rem;
	}

	@keyframes spin {
		0% { transform: rotate(0deg); }
		100% { transform: rotate(360deg); }
	}

	.error {
		background: #fef2f2;
		border: 1px solid #fecaca;
		border-radius: 8px;
		padding: 2rem;
		text-align: center;
		color: #dc2626;
	}

	.retry-button {
		background: #dc2626;
		color: white;
		border: none;
		padding: 0.5rem 1rem;
		border-radius: 4px;
		cursor: pointer;
		margin-top: 1rem;
	}

	.empty-state {
		text-align: center;
		padding: 4rem 2rem;
		color: #666;
	}

	.empty-icon {
		font-size: 4rem;
		margin-bottom: 1rem;
	}

	.primary-button {
		display: inline-block;
		background: linear-gradient(135deg, #4ade80, #22d3ee);
		color: white;
		padding: 1rem 2rem;
		border-radius: 8px;
		text-decoration: none;
		font-weight: 600;
		margin-top: 1rem;
		transition: transform 0.2s, box-shadow 0.2s;
	}

	.primary-button:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(0,0,0,0.15);
	}

	.reminders-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
		gap: 1.5rem;
	}

	.reminder-card {
		background: white;
		border: 1px solid #e5e7eb;
		border-radius: 12px;
		padding: 1.5rem;
		box-shadow: 0 2px 4px rgba(0,0,0,0.05);
		transition: transform 0.2s, box-shadow 0.2s;
	}

	.reminder-card:hover {
		transform: translateY(-2px);
		box-shadow: 0 8px 25px rgba(0,0,0,0.1);
	}

	.card-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		margin-bottom: 1rem;
	}

	.reminder-title {
		color: #059669;
		font-size: 1.25rem;
		margin: 0;
		flex: 1;
		line-height: 1.3;
	}

	.card-actions {
		display: flex;
		gap: 0.5rem;
		margin-left: 1rem;
	}

	.edit-button,
	.delete-button {
		background: none;
		border: none;
		cursor: pointer;
		padding: 0.25rem;
		border-radius: 4px;
		transition: background-color 0.2s;
		font-size: 1.1rem;
		text-decoration: none;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.edit-button:hover {
		background: #f0f9ff;
	}

	.delete-button:hover {
		background: #fef2f2;
	}

	.reminder-description {
		color: #666;
		margin-bottom: 1rem;
		line-height: 1.5;
	}

	.reminder-details {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.detail-row {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		flex-wrap: wrap;
	}

	.label {
		font-weight: 600;
		color: #374151;
		min-width: 80px;
	}

	.badge {
		padding: 0.25rem 0.75rem;
		border-radius: 12px;
		font-size: 0.875rem;
		font-weight: 500;
		text-transform: capitalize;
	}

	.channels {
		display: flex;
		gap: 0.5rem;
		flex-wrap: wrap;
	}

	.channel-badge {
		background: #f0fdf4;
		color: #059669;
		padding: 0.25rem 0.5rem;
		border-radius: 8px;
		font-size: 0.75rem;
		font-weight: 500;
		border: 1px solid #bbf7d0;
	}

	@media (max-width: 640px) {
		.reminders-grid {
			grid-template-columns: 1fr;
		}
		
		.header {
			flex-direction: column;
			align-items: stretch;
		}
		
		.add-button {
			justify-content: center;
		}
	}
</style>