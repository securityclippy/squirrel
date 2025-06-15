<script lang="ts">
	import { authStore, login, logout } from '$lib/auth0';

	async function handleLogin() {
		try {
			await login();
		} catch (error) {
			console.error('Login failed:', error);
		}
	}

	async function handleLogout() {
		try {
			await logout();
		} catch (error) {
			console.error('Logout failed:', error);
		}
	}
</script>

{#if $authStore.isLoading}
	<button disabled class="btn btn-secondary">
		<span class="loading">Loading...</span>
	</button>
{:else if $authStore.isAuthenticated}
	<div class="auth-menu">
		<span class="user-info">
			{#if $authStore.user?.picture}
				<img src={$authStore.user.picture} alt="User avatar" class="user-avatar" />
			{/if}
			<span class="user-name">{$authStore.user?.name || $authStore.user?.email || 'User'}</span>
		</span>
		<button on:click={handleLogout} class="btn btn-outline">Logout</button>
	</div>
{:else}
	<button on:click={handleLogin} class="btn btn-primary">Login</button>
{/if}

{#if $authStore.error}
	<div class="error">
		Authentication error: {$authStore.error.message}
	</div>
{/if}

<style>
	.auth-menu {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.user-info {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.user-avatar {
		width: 32px;
		height: 32px;
		border-radius: 50%;
		object-fit: cover;
	}

	.user-name {
		font-weight: 500;
		color: var(--text-primary, #333);
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
		background-color: var(--primary-color, #3b82f6);
		color: white;
	}

	.btn-primary:hover {
		background-color: var(--primary-hover, #2563eb);
	}

	.btn-outline {
		background-color: transparent;
		color: var(--text-primary, #333);
		border: 1px solid var(--border-color, #d1d5db);
	}

	.btn-outline:hover {
		background-color: var(--gray-50, #f9fafb);
	}

	.btn-secondary {
		background-color: var(--gray-200, #e5e7eb);
		color: var(--text-secondary, #6b7280);
	}

	.btn:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.loading {
		animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
	}

	.error {
		color: var(--error-color, #dc2626);
		font-size: 0.875rem;
		margin-top: 0.5rem;
		padding: 0.5rem;
		background-color: var(--error-bg, #fef2f2);
		border: 1px solid var(--error-border, #fecaca);
		border-radius: 0.375rem;
	}

	@keyframes pulse {
		0%, 100% {
			opacity: 1;
		}
		50% {
			opacity: 0.5;
		}
	}
</style>