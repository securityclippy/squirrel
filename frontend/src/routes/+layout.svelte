<script>
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';
	import { initializeAuth0 } from '$lib/auth0';
	import Header from '$lib/components/Header.svelte';
	import Sidebar from '$lib/components/Sidebar.svelte';

	onMount(async () => {
		if (browser) {
			try {
				await initializeAuth0();
			} catch (error) {
				console.error('Failed to initialize Auth0:', error);
			}
		}
	});
</script>

<div class="app">
	<Header />

	<div class="main-layout">
		<Sidebar />
		
		<main class="content">
			<slot />
		</main>
	</div>
</div>

<style>
	.app {
		display: flex;
		flex-direction: column;
		height: 100vh;
	}

	.main-layout {
		display: flex;
		flex: 1;
		overflow: hidden;
	}

	.content {
		flex: 1;
		padding: 2rem;
		overflow-y: auto;
		background: #ffffff;
	}

	@media (max-width: 768px) {
		.main-layout {
			flex-direction: column;
		}
	}
</style>