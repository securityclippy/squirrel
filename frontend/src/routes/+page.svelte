<svelte:head>
	<title>Squirrel Reminders</title>
</svelte:head>

<div class="welcome">
	<h2>Welcome to Squirrel Reminders!</h2>
	<p>Keep track of all your important tasks and never forget anything again. 🌰</p>
	
	<div class="quick-actions">
		<a href="/reminders" class="action-button primary">View Reminders</a>
		<a href="/reminders/new" class="action-button secondary">Add New Reminder</a>
	</div>
</div>

<style>
	.welcome {
		text-align: center;
		max-width: 600px;
		margin: 3rem auto;
	}

	.welcome h2 {
		color: #059669;
		font-size: 2.5rem;
		margin-bottom: 1rem;
	}

	.welcome p {
		font-size: 1.2rem;
		color: #666;
		margin-bottom: 2rem;
	}

	.quick-actions {
		display: flex;
		gap: 1rem;
		justify-content: center;
		flex-wrap: wrap;
	}

	.action-button {
		padding: 1rem 2rem;
		border-radius: 8px;
		text-decoration: none;
		font-weight: 600;
		transition: transform 0.2s, box-shadow 0.2s;
		display: inline-block;
	}

	.action-button:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(0,0,0,0.15);
	}

	.action-button.primary {
		background: linear-gradient(135deg, #4ade80, #22d3ee);
		color: white;
	}

	.action-button.secondary {
		background: white;
		color: #059669;
		border: 2px solid #059669;
	}

	.action-button.secondary:hover {
		background: #059669;
		color: white;
	}
</style>