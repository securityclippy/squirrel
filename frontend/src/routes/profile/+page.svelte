<script>
	import { onMount } from 'svelte';
	import { authStore } from '$lib/auth0';
	import AuthGuard from '$lib/components/AuthGuard.svelte';

	let loading = false;
	let error = '';
	let success = '';
	
	// Profile editing
	let editMode = false;
	let editedName = '';
	let editedEmail = '';
	

	onMount(() => {
		// Initialize edit form with current user data
		if ($authStore.user) {
			editedName = $authStore.user.name || '';
			editedEmail = $authStore.user.email || '';
		}
	});

	function startEdit() {
		editMode = true;
		editedName = $authStore.user?.name || '';
		editedEmail = $authStore.user?.email || '';
	}

	function cancelEdit() {
		editMode = false;
		error = '';
	}

	async function saveProfile() {
		try {
			loading = true;
			error = '';
			
			// TODO: Implement API call to update profile
			// await api.updateUserProfile({
			//   name: editedName.trim(),
			//   email: editedEmail.trim()
			// });
			
			// Simulate API call
			await new Promise(resolve => setTimeout(resolve, 1000));
			
			editMode = false;
			success = 'Profile updated successfully!';
			setTimeout(() => { success = ''; }, 5000);
		} catch (err) {
			error = 'Failed to update profile: ' + (err.message || 'Unknown error');
		} finally {
			loading = false;
		}
	}

</script>

<svelte:head>
	<title>Profile - Squirrel Reminders</title>
</svelte:head>

<AuthGuard>
	<div class="profile-page">
		<div class="header">
			<h1>Profile</h1>
			<p class="subtitle">Manage your account information and view statistics</p>
		</div>

		{#if error}
			<div class="alert alert-error">
				{error}
				<button class="close-btn" on:click={() => error = ''}>×</button>
			</div>
		{/if}

		{#if success}
			<div class="alert alert-success">
				{success}
				<button class="close-btn" on:click={() => success = ''}>×</button>
			</div>
		{/if}

		<!-- User Profile Section -->
		{#if $authStore.user}
			<div class="section">
				<div class="section-header">
					<h2>👤 Profile Information</h2>
					{#if !editMode}
						<button class="btn btn-secondary" on:click={startEdit}>Edit Profile</button>
					{/if}
				</div>

				<div class="profile-content">
					<div class="profile-avatar">
						{#if $authStore.user.picture}
							<img src={$authStore.user.picture} alt="User avatar" class="user-avatar-large" />
						{:else}
							<div class="avatar-fallback">
								{($authStore.user.name || $authStore.user.email || 'U').charAt(0).toUpperCase()}
							</div>
						{/if}
					</div>

					<div class="profile-details">
						{#if editMode}
							<!-- Edit Mode -->
							<div class="edit-form">
								<div class="form-group">
									<label for="edit-name">Full Name</label>
									<input 
										id="edit-name"
										type="text" 
										bind:value={editedName}
										placeholder="Enter your full name"
									/>
								</div>
								
								<div class="form-group">
									<label for="edit-email">Email Address</label>
									<input 
										id="edit-email"
										type="email" 
										bind:value={editedEmail}
										placeholder="Enter your email address"
									/>
								</div>

								<div class="form-actions">
									<button 
										class="btn btn-primary" 
										on:click={saveProfile}
										disabled={loading}
									>
										{loading ? 'Saving...' : 'Save Changes'}
									</button>
									<button class="btn btn-secondary" on:click={cancelEdit}>
										Cancel
									</button>
								</div>
							</div>
						{:else}
							<!-- View Mode -->
							<div class="profile-info">
								<div class="info-item">
									<label>Full Name</label>
									<span class="value">{$authStore.user.name || 'Not set'}</span>
								</div>
								
								<div class="info-item">
									<label>Email Address</label>
									<span class="value">{$authStore.user.email}</span>
								</div>
								
								<div class="info-item">
									<label>Email Status</label>
									{#if $authStore.user.email_verified}
										<span class="badge verified">✓ Verified</span>
									{:else}
										<span class="badge unverified">Not Verified</span>
									{/if}
								</div>
							</div>
						{/if}
					</div>
				</div>
			</div>
		{/if}


	</div>
</AuthGuard>

<style>
	.profile-page {
		max-width: 800px;
		margin: 0 auto;
		padding: 2rem;
	}

	.header {
		margin-bottom: 2rem;
	}

	.header h1 {
		margin: 0 0 0.5rem 0;
		color: #111827;
		font-size: 2rem;
		font-weight: 600;
	}

	.subtitle {
		margin: 0;
		color: #6b7280;
		font-size: 1rem;
	}

	.profile-content {
		display: flex;
		gap: 2rem;
		align-items: flex-start;
	}

	.profile-avatar {
		flex-shrink: 0;
	}

	.user-avatar-large {
		width: 120px;
		height: 120px;
		border-radius: 50%;
		object-fit: cover;
		border: 4px solid #e5e7eb;
	}

	.avatar-fallback {
		width: 120px;
		height: 120px;
		border-radius: 50%;
		background: linear-gradient(135deg, #6366f1, #8b5cf6);
		color: white;
		display: flex;
		align-items: center;
		justify-content: center;
		font-weight: 600;
		font-size: 2.5rem;
		border: 4px solid #e5e7eb;
	}

	.profile-details {
		flex: 1;
	}

	.profile-info {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.info-item {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.info-item label {
		font-weight: 600;
		color: #374151;
		font-size: 0.875rem;
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	.info-item .value {
		color: #111827;
		font-size: 1rem;
	}

	.edit-form {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}



	.section {
		background: white;
		border-radius: 12px;
		padding: 2rem;
		margin-bottom: 2rem;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
		border: 1px solid #e5e7eb;
	}

	.section h2 {
		margin: 0 0 1.5rem 0;
		color: #374151;
		font-size: 1.5rem;
		font-weight: 600;
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.section-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1rem;
	}

	.user-card {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.user-avatar-large {
		width: 64px;
		height: 64px;
		border-radius: 50%;
		object-fit: cover;
	}

	.user-details h3 {
		margin: 0 0 0.25rem 0;
		color: #111827;
		font-size: 1.25rem;
		font-weight: 600;
	}

	.user-email {
		margin: 0 0 0.5rem 0;
		color: #6b7280;
	}

	.badge {
		padding: 0.25rem 0.5rem;
		border-radius: 0.25rem;
		font-size: 0.75rem;
		font-weight: 500;
	}

	.badge.verified {
		background-color: #d1fae5;
		color: #065f46;
	}

	.badge.unverified {
		background-color: #fee2e2;
		color: #991b1b;
	}

	.alert {
		padding: 1rem;
		border-radius: 8px;
		margin-bottom: 1.5rem;
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.alert-error {
		background: #fef2f2;
		color: #991b1b;
		border: 1px solid #fecaca;
	}

	.alert-success {
		background: #f0fdf4;
		color: #166534;
		border: 1px solid #bbf7d0;
	}

	.close-btn {
		background: none;
		border: none;
		font-size: 1.5rem;
		cursor: pointer;
		opacity: 0.7;
		margin-left: 1rem;
	}

	.close-btn:hover {
		opacity: 1;
	}

	.key-display {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		margin: 0.5rem 0;
		padding: 0.5rem;
		background-color: #f3f4f6;
		border-radius: 0.25rem;
	}

	.key-display code {
		flex: 1;
		font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
		word-break: break-all;
	}

	.copy-btn {
		padding: 0.25rem 0.5rem;
		background-color: #3b82f6;
		color: white;
		border: none;
		border-radius: 0.25rem;
		cursor: pointer;
		font-size: 0.75rem;
	}

	.copy-btn:hover {
		background-color: #2563eb;
	}

	.create-form {
		border: 1px solid #e5e7eb;
		border-radius: 0.5rem;
		padding: 1.5rem;
		margin-bottom: 1.5rem;
		background-color: #f9fafb;
	}

	.create-form h3 {
		margin: 0 0 1rem 0;
		color: #374151;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		margin-bottom: 1.5rem;
	}

	.form-group label {
		font-weight: 600;
		color: #374151;
		font-size: 0.875rem;
	}

	.form-group input[type="text"],
	.form-group input[type="email"] {
		padding: 0.75rem;
		border: 1px solid #d1d5db;
		border-radius: 8px;
		font-size: 0.875rem;
		background: white;
	}

	.form-group input:focus {
		outline: none;
		border-color: #3b82f6;
		box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
	}

	.permissions-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
		gap: 0.5rem;
	}

	.permission-option {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		cursor: pointer;
	}

	.form-actions {
		display: flex;
		gap: 0.5rem;
		margin-top: 1.5rem;
	}

	.btn {
		padding: 0.75rem 1.5rem;
		border: none;
		border-radius: 8px;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s ease;
		text-decoration: none;
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.875rem;
	}

	.btn-primary {
		background: #3b82f6;
		color: white;
	}

	.btn-primary:hover {
		background: #2563eb;
	}

	.btn-secondary {
		background: #f3f4f6;
		color: #374151;
		border: 1px solid #d1d5db;
	}

	.btn-secondary:hover {
		background: #e5e7eb;
	}

	.btn-danger {
		background: #dc2626;
		color: white;
	}

	.btn-danger:hover {
		background: #b91c1c;
	}

	.btn-sm {
		padding: 0.5rem 1rem;
		font-size: 0.75rem;
	}

	.btn:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.loading {
		text-align: center;
		padding: 2rem;
		color: #6b7280;
	}

	.empty-state {
		text-align: center;
		padding: 2rem;
		color: #6b7280;
	}

	.api-keys-list {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.api-key-card {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		padding: 1rem;
		border: 1px solid #e5e7eb;
		border-radius: 0.5rem;
		background-color: white;
	}

	.key-info h4 {
		margin: 0 0 0.5rem 0;
		color: #111827;
		font-weight: 600;
	}

	.key-meta {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		margin-bottom: 0.5rem;
	}

	.key-prefix {
		font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
		font-size: 0.875rem;
		color: #6b7280;
	}

	.permissions {
		font-size: 0.875rem;
		color: #6b7280;
	}

	.key-dates {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		font-size: 0.75rem;
		color: #9ca3af;
	}

	.key-actions {
		display: flex;
		flex-direction: column;
		align-items: flex-end;
		gap: 0.5rem;
	}

	.status {
		padding: 0.25rem 0.5rem;
		border-radius: 0.25rem;
		font-size: 0.75rem;
		font-weight: 500;
	}

	.status.active {
		background-color: #d1fae5;
		color: #065f46;
	}

	.status.inactive {
		background-color: #fee2e2;
		color: #991b1b;
	}

	@media (max-width: 768px) {
		.profile-page {
			padding: 1rem;
		}

		.section {
			padding: 1.5rem;
		}

		.section-header {
			flex-direction: column;
			align-items: stretch;
			gap: 1rem;
		}

		.profile-content {
			flex-direction: column;
			align-items: center;
			text-align: center;
			gap: 1.5rem;
		}



		.form-actions {
			flex-direction: column;
		}
	}

</style>