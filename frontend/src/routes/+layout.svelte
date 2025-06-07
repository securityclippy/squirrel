<script>
	import Header from '$lib/components/Header.svelte';
	import Sidebar from '$lib/components/Sidebar.svelte';
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