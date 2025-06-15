<script>
	import { onMount } from 'svelte';
	import { authStore } from '$lib/auth0';
	import { api, type APIKey, type CreateAPIKeyRequest } from '$lib/api';
	import AuthGuard from '$lib/components/AuthGuard.svelte';

	let apiKeys = [];
	let loading = false;
	let error = '';
	let success = '';
	
	// API Key creation form
	let showCreateForm = false;
	let newKeyName = '';
	let newKeyPermissions = ['read'];
	let createdKey = '';

	const permissionOptions = [
		{ value: 'read', label: 'Read' },
		{ value: 'write', label: 'Write' },
		{ value: 'delete', label: 'Delete' }
	];

	onMount(async () => {
		await loadAPIKeys();
	});

	async function loadAPIKeys() {
		try {
			loading = true;
			error = '';
			const response = await api.getAPIKeys();
			apiKeys = response.api_keys;
		} catch (err) {
			error = 'Failed to load API keys: ' + err.message;
		} finally {
			loading = false;
		}
	}

	async function createAPIKey() {
		if (!newKeyName.trim()) {
			error = 'Please enter a name for the API key';
			return;
		}

		try {
			loading = true;
			error = '';
			success = '';
			
			const response = await api.createAPIKey({
				name: newKeyName.trim(),
				permissions: newKeyPermissions
			});
			
			createdKey = response.key;
			success = 'API key created successfully! Make sure to copy it now - you won\'t be able to see it again.';
			
			// Reset form
			newKeyName = '';
			newKeyPermissions = ['read'];
			showCreateForm = false;
			
			// Reload API keys
			await loadAPIKeys();
		} catch (err) {
			error = 'Failed to create API key: ' + err.message;
		} finally {
			loading = false;
		}
	}

	async function revokeAPIKey(id, name) {
		if (!confirm(`Are you sure you want to revoke the API key "${name}"? This action cannot be undone.`)) {
			return;
		}

		try {
			loading = true;
			error = '';
			success = '';
			
			await api.revokeAPIKey(id);
			success = `API key "${name}" has been revoked successfully.`;
			
			// Reload API keys
			await loadAPIKeys();
		} catch (err) {
			error = 'Failed to revoke API key: ' + err.message;
		} finally {
			loading = false;
		}
	}

	function copyToClipboard(text) {
		navigator.clipboard.writeText(text).then(() => {
			success = 'API key copied to clipboard!';
			setTimeout(() => { success = ''; }, 3000);
		}).catch(() => {
			error = 'Failed to copy to clipboard';
		});
	}

	function formatDate(dateString) {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function togglePermission(permission) {
		if (newKeyPermissions.includes(permission)) {
			newKeyPermissions = newKeyPermissions.filter(p => p !== permission);
		} else {
			newKeyPermissions = [...newKeyPermissions, permission];
		}
	}
</script>

<svelte:head>
	<title>Profile - Squirrel Reminders</title>
</svelte:head>

<AuthGuard>
	<div class="profile-page">
		<div class="header">
			<h1>Profile & Settings</h1>
		</div>

		<!-- User Info Section -->
		{#if $authStore.user}
			<div class="section">
				<h2>User Information</h2>
				<div class="user-card">
					{#if $authStore.user.picture}
						<img src={$authStore.user.picture} alt="User avatar" class="user-avatar-large" />
					{/if}
					<div class="user-details">
						<h3>{$authStore.user.name || 'Unknown User'}</h3>
						<p class="user-email">{$authStore.user.email}</p>
						{#if $authStore.user.email_verified}
							<span class="badge verified">✓ Email Verified</span>
						{:else}
							<span class="badge unverified">Email Not Verified</span>
						{/if}
					</div>
				</div>
			</div>
		{/if}

		<!-- API Keys Section -->
		<div class="section">
			<div class="section-header">
				<h2>API Keys</h2>
				<button on:click={() => showCreateForm = !showCreateForm} class="btn btn-primary">
					{showCreateForm ? 'Cancel' : 'Create New API Key'}
				</button>
			</div>

			{#if error}
				<div class="alert alert-error">{error}</div>
			{/if}

			{#if success}
				<div class="alert alert-success">{success}</div>
			{/if}

			{#if createdKey}
				<div class="alert alert-info">
					<h4>Your new API key:</h4>
					<div class="key-display">
						<code>{createdKey}</code>
						<button on:click={() => copyToClipboard(createdKey)} class="copy-btn">Copy</button>
					</div>
					<p><strong>Important:</strong> Save this API key now. You won't be able to see it again!</p>
					<button on:click={() => createdKey = ''} class="btn btn-sm">Dismiss</button>
				</div>
			{/if}

			<!-- Create API Key Form -->
			{#if showCreateForm}
				<div class="create-form">
					<h3>Create New API Key</h3>
					<div class="form-group">
						<label for="key-name">Name:</label>
						<input 
							id="key-name"
							type="text" 
							bind:value={newKeyName} 
							placeholder="e.g., Mobile App, Dashboard Access"
							maxlength="255"
						/>
					</div>
					
					<div class="form-group">
						<label>Permissions:</label>
						<div class="permissions-grid">
							{#each permissionOptions as option}
								<label class="permission-option">
									<input 
										type="checkbox" 
										checked={newKeyPermissions.includes(option.value)}
										on:change={() => togglePermission(option.value)}
									/>
									{option.label}
								</label>
							{/each}
						</div>
					</div>

					<div class="form-actions">
						<button on:click={createAPIKey} class="btn btn-primary" disabled={loading}>
							{loading ? 'Creating...' : 'Create API Key'}
						</button>
						<button on:click={() => showCreateForm = false} class="btn btn-secondary">
							Cancel
						</button>
					</div>
				</div>
			{/if}

			<!-- API Keys List -->
			{#if loading}
				<div class="loading">Loading API keys...</div>
			{:else if apiKeys.length === 0}
				<div class="empty-state">
					<p>No API keys found. Create your first API key to access the API programmatically.</p>
				</div>
			{:else}
				<div class="api-keys-list">
					{#each apiKeys as key}
						<div class="api-key-card">
							<div class="key-info">
								<h4>{key.name}</h4>
								<div class="key-meta">
									<span class="key-prefix">Key: {key.key_prefix}••••••••</span>
									<span class="permissions">
										Permissions: {key.permissions.join(', ')}
									</span>
								</div>
								<div class="key-dates">
									<span>Created: {formatDate(key.created_at)}</span>
									{#if key.last_used_at}
										<span>Last used: {formatDate(key.last_used_at)}</span>
									{:else}
										<span>Never used</span>
									{/if}
								</div>
							</div>
							<div class="key-actions">
								<span class="status {key.is_active ? 'active' : 'inactive'}">
									{key.is_active ? 'Active' : 'Revoked'}
								</span>
								{#if key.is_active}
									<button 
										on:click={() => revokeAPIKey(key.id, key.name)} 
										class="btn btn-danger btn-sm"
										disabled={loading}
									>
										Revoke
									</button>
								{/if}
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</div>
</AuthGuard>

<style>
	.profile-page {
		max-width: 800px;
		margin: 0 auto;
		padding: 2rem;
	}

	.header h1 {
		margin: 0 0 2rem 0;
		color: #111827;
		font-size: 2rem;
		font-weight: 600;
	}

	.section {
		background: white;
		border-radius: 0.5rem;
		padding: 1.5rem;
		margin-bottom: 2rem;
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
		border: 1px solid #e5e7eb;
	}

	.section h2 {
		margin: 0 0 1rem 0;
		color: #374151;
		font-size: 1.25rem;
		font-weight: 600;
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
		border-radius: 0.375rem;
		margin-bottom: 1rem;
	}

	.alert-error {
		background-color: #fef2f2;
		color: #991b1b;
		border: 1px solid #fecaca;
	}

	.alert-success {
		background-color: #f0fdf4;
		color: #166534;
		border: 1px solid #bbf7d0;
	}

	.alert-info {
		background-color: #eff6ff;
		color: #1e40af;
		border: 1px solid #bfdbfe;
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
		margin-bottom: 1rem;
	}

	.form-group label {
		display: block;
		margin-bottom: 0.25rem;
		font-weight: 500;
		color: #374151;
	}

	.form-group input[type="text"] {
		width: 100%;
		padding: 0.5rem;
		border: 1px solid #d1d5db;
		border-radius: 0.375rem;
		font-size: 0.875rem;
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
		padding: 0.5rem 1rem;
		border: none;
		border-radius: 0.375rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s ease;
		text-decoration: none;
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
	}

	.btn-primary {
		background-color: #3b82f6;
		color: white;
	}

	.btn-primary:hover {
		background-color: #2563eb;
	}

	.btn-secondary {
		background-color: #f3f4f6;
		color: #374151;
	}

	.btn-secondary:hover {
		background-color: #e5e7eb;
	}

	.btn-danger {
		background-color: #dc2626;
		color: white;
	}

	.btn-danger:hover {
		background-color: #b91c1c;
	}

	.btn-sm {
		padding: 0.25rem 0.75rem;
		font-size: 0.875rem;
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

	@media (max-width: 640px) {
		.profile-page {
			padding: 1rem;
		}

		.section-header {
			flex-direction: column;
			align-items: stretch;
			gap: 1rem;
		}

		.user-card {
			flex-direction: column;
			text-align: center;
		}

		.api-key-card {
			flex-direction: column;
			gap: 1rem;
		}

		.key-actions {
			align-items: stretch;
		}

		.permissions-grid {
			grid-template-columns: 1fr;
		}

		.form-actions {
			flex-direction: column;
		}
	}
</style>