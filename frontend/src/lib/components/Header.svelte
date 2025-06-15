<script>
	import AuthButton from './AuthButton.svelte';
</script>

<header class="top-bar">
	<div class="header-content">
		<div class="title-section">
			<p class="flavor-text">reminders for people who are easily distracted by...</p>
		</div>
		<div class="header-right">
			<div class="logo-section">
				<div class="squirrel-logo">🐿️</div>
				<div class="exclamation">!</div>
			</div>
			<div class="auth-section">
				<AuthButton />
			</div>
		</div>
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

	.header-content {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.title-section {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.header-right {
		display: flex;
		align-items: center;
		gap: 2rem;
	}

	.logo-section {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.auth-section {
		display: flex;
		align-items: center;
	}

	.flavor-text {
		margin: 0;
		font-size: 1.25rem;
		color: black;
		font-style: italic;
	}

	.squirrel-logo {
		font-size: 3rem;
		animation: bounce 2s infinite;
	}

	.exclamation {
		font-size: 3rem;
		font-weight: bold;
		color: #3b82f6;
		text-shadow: 0 0 10px #3b82f6, 0 0 20px #3b82f6, 0 0 30px #3b82f6;
		animation: bounce 2s infinite;
	}

	@keyframes bounce {
		0%, 100% { transform: translateY(0); }
		50% { transform: translateY(-5px); }
	}

	@keyframes pulse {
		0%, 100% { transform: scale(1); }
		50% { transform: scale(1.2); }
	}

	@media (max-width: 768px) {
		.top-bar {
			padding: 1rem;
		}
		
		.flavor-text {
			font-size: 1rem;
		}
		
		.squirrel-logo {
			font-size: 2.5rem;
		}

		.exclamation {
			font-size: 2.5rem;
		}
	}
</style>