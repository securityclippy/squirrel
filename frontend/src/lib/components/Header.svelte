<header class="top-bar">
	<div class="logo-section">
		<div class="squirrel-logo">🐿️</div>
		<h1>Squirrel Reminders</h1>
	</div>
</header>

<style>
	.top-bar {
		background: linear-gradient(135deg, #4ade80, #22d3ee);
		color: white;
		padding: 1rem 2rem;
		box-shadow: 0 2px 4px rgba(0,0,0,0.1);
		z-index: 10;
	}

	.logo-section {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.squirrel-logo {
		font-size: 2rem;
		animation: bounce 2s infinite;
	}

	@keyframes bounce {
		0%, 100% { transform: translateY(0); }
		50% { transform: translateY(-5px); }
	}

	.top-bar h1 {
		margin: 0;
		font-size: 1.5rem;
		font-weight: 600;
	}
</style>