<script lang="ts">
	import { authStore, login } from '$lib/auth0';
	import { onMount } from 'svelte';

	export let fallback: boolean = true;
	export let showLoginButton: boolean = true;

	async function handleLogin() {
		try {
			await login();
		} catch (error) {
			console.error('Login failed:', error);
		}
	}
</script>

{#if $authStore.isLoading}
	<div class="auth-loading">
		<div class="spinner"></div>
		<p>Loading authentication...</p>
	</div>
{:else if $authStore.isAuthenticated}
	<!-- User is authenticated, show the protected content -->
	<slot />
{:else if fallback}
	<!-- User is not authenticated, show fallback content -->
	<div class="auth-required">
		<div class="auth-card">
			<h2>Authentication Required</h2>
			<p>You need to be logged in to access this page.</p>
			
			{#if showLoginButton}
				<button on:click={handleLogin} class="btn btn-primary">
					Login to Continue
				</button>
			{/if}
			
			{#if $authStore.error}
				<div class="error">
					{$authStore.error.message}
				</div>
			{/if}
		</div>
	</div>
{/if}

<style>
	.auth-loading {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		min-height: 200px;
		gap: 1rem;
		color: var(--text-secondary, #6b7280);
	}

	.spinner {
		width: 2rem;
		height: 2rem;
		border: 2px solid var(--gray-200, #e5e7eb);
		border-top: 2px solid var(--primary-color, #3b82f6);
		border-radius: 50%;
		animation: spin 1s linear infinite;
	}

	.auth-required {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 400px;
		padding: 2rem;
	}

	.auth-card {
		background: white;
		border-radius: 0.5rem;
		padding: 2rem;
		box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
		border: 1px solid var(--border-color, #d1d5db);
		text-align: center;
		max-width: 400px;
		width: 100%;
	}

	.auth-card h2 {
		margin: 0 0 1rem 0;
		color: var(--text-primary, #111827);
		font-size: 1.5rem;
		font-weight: 600;
	}

	.auth-card p {
		margin: 0 0 1.5rem 0;
		color: var(--text-secondary, #6b7280);
		line-height: 1.5;
	}

	.btn {
		padding: 0.75rem 1.5rem;
		border: none;
		border-radius: 0.375rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s ease;
		text-decoration: none;
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 1rem;
	}

	.btn-primary {
		background-color: var(--primary-color, #3b82f6);
		color: white;
	}

	.btn-primary:hover {
		background-color: var(--primary-hover, #2563eb);
		transform: translateY(-1px);
	}

	.error {
		color: var(--error-color, #dc2626);
		font-size: 0.875rem;
		margin-top: 1rem;
		padding: 0.5rem;
		background-color: var(--error-bg, #fef2f2);
		border: 1px solid var(--error-border, #fecaca);
		border-radius: 0.375rem;
	}

	@keyframes spin {
		0% {
			transform: rotate(0deg);
		}
		100% {
			transform: rotate(360deg);
		}
	}
</style>