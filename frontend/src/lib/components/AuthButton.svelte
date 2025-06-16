<script lang="ts">
	import { authStore, login, logout } from '$lib/auth0';
	import { onMount } from 'svelte';
	
	let showDropdown = false;
	let dropdownElement: HTMLDivElement;

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

	function toggleDropdown() {
		showDropdown = !showDropdown;
	}

	function closeDropdown() {
		showDropdown = false;
	}

	// Close dropdown when clicking outside
	onMount(() => {
		function handleClickOutside(event: MouseEvent) {
			if (dropdownElement && !dropdownElement.contains(event.target as Node)) {
				showDropdown = false;
			}
		}

		document.addEventListener('click', handleClickOutside);
		return () => {
			document.removeEventListener('click', handleClickOutside);
		};
	});
</script>

{#if $authStore.isLoading}
	<button disabled class="btn btn-secondary">
		<span class="loading">Loading...</span>
	</button>
{:else if $authStore.isAuthenticated}
	<div class="user-menu" bind:this={dropdownElement}>
		<button on:click={toggleDropdown} class="user-avatar-btn">
			{#if $authStore.user?.picture}
				<img src={$authStore.user.picture} alt="User avatar" class="user-avatar" />
			{:else}
				<div class="user-avatar-fallback">
					{($authStore.user?.name || $authStore.user?.email || 'U').charAt(0).toUpperCase()}
				</div>
			{/if}
		</button>
		
		{#if showDropdown}
			<div class="dropdown-menu">
				<div class="dropdown-header">
					<div class="user-name">{$authStore.user?.name || 'User'}</div>
					<div class="user-email">{$authStore.user?.email || ''}</div>
				</div>
				<div class="dropdown-divider"></div>
				<a href="/profile" class="dropdown-item" on:click={closeDropdown}>
					<span class="dropdown-icon">👤</span>
					Profile
				</a>
				<a href="/settings" class="dropdown-item" on:click={closeDropdown}>
					<span class="dropdown-icon">⚙️</span>
					Settings
				</a>
				<a href="/statistics" class="dropdown-item" on:click={closeDropdown}>
					<span class="dropdown-icon">📊</span>
					Statistics
				</a>
				<div class="dropdown-divider"></div>
				<button on:click={() => { handleLogout(); closeDropdown(); }} class="dropdown-item logout-item">
					<span class="dropdown-icon">🚪</span>
					Logout
				</button>
			</div>
		{/if}
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
	.user-menu {
		position: relative;
		display: inline-block;
	}

	.user-avatar-btn {
		background: none;
		border: none;
		cursor: pointer;
		padding: 2px;
		border-radius: 50%;
		transition: all 0.2s ease;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.user-avatar-btn:hover {
		transform: scale(1.05);
		box-shadow: 0 2px 8px rgba(0,0,0,0.2);
	}

	.user-avatar {
		width: 40px;
		height: 40px;
		border-radius: 50%;
		object-fit: cover;
		border: 2px solid rgba(255,255,255,0.8);
	}

	.user-avatar-fallback {
		width: 40px;
		height: 40px;
		border-radius: 50%;
		background: linear-gradient(135deg, #6366f1, #8b5cf6);
		color: white;
		display: flex;
		align-items: center;
		justify-content: center;
		font-weight: 600;
		font-size: 1rem;
		border: 2px solid rgba(255,255,255,0.8);
	}

	.dropdown-menu {
		position: absolute;
		top: calc(100% + 8px);
		right: 0;
		background: white;
		border-radius: 8px;
		box-shadow: 0 4px 20px rgba(0,0,0,0.15);
		border: 1px solid #e5e7eb;
		min-width: 220px;
		z-index: 1000;
		overflow: hidden;
		animation: dropdownFadeIn 0.2s ease-out;
	}

	@keyframes dropdownFadeIn {
		from {
			opacity: 0;
			transform: translateY(-4px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	.dropdown-header {
		padding: 12px 16px;
		background: #f8fafc;
		border-bottom: 1px solid #e5e7eb;
	}

	.user-name {
		font-weight: 600;
		color: #111827;
		font-size: 0.875rem;
		margin-bottom: 2px;
	}

	.user-email {
		font-size: 0.75rem;
		color: #6b7280;
		word-break: break-word;
	}

	.dropdown-divider {
		height: 1px;
		background: #e5e7eb;
	}

	.dropdown-item {
		display: flex;
		align-items: center;
		gap: 8px;
		width: 100%;
		padding: 12px 16px;
		border: none;
		background: none;
		text-align: left;
		cursor: pointer;
		transition: background-color 0.2s ease;
		color: #374151;
		text-decoration: none;
		font-size: 0.875rem;
	}

	.dropdown-item:hover {
		background: #f3f4f6;
	}

	.dropdown-icon {
		font-size: 1rem;
		opacity: 0.7;
	}

	.logout-item {
		color: #dc2626;
	}

	.logout-item:hover {
		background: #fef2f2;
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