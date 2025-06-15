<script>
	import { page } from '$app/stores';
	
	const navItems = [
		{ href: '/', label: 'Home', icon: '🏠' },
		{ href: '/reminders', label: 'Reminders', icon: '📝' },
		{ href: '/profile', label: 'Profile', icon: '👤' },
		{ href: '/settings', label: 'Settings', icon: '⚙️' }
	];
</script>

<nav class="sidebar">
	<ul>
		{#each navItems as item}
			<li>
				<a 
					href={item.href} 
					class:active={$page.url.pathname === item.href}
				>
					<span class="icon">{item.icon}</span>
					<span class="label">{item.label}</span>
				</a>
			</li>
		{/each}
	</ul>
</nav>

<style>
	.sidebar {
		width: 250px;
		background: #f0fdf4;
		border-right: 1px solid #bbf7d0;
		padding: 1rem 0;
		height: 100%;
		overflow-y: auto;
	}

	.sidebar ul {
		list-style: none;
		margin: 0;
		padding: 0;
	}

	.sidebar li {
		margin: 0.5rem 0;
	}

	.sidebar a {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.75rem 1.5rem;
		color: #374151;
		text-decoration: none;
		transition: all 0.2s;
		border-radius: 0 25px 25px 0;
		margin-right: 1rem;
	}

	.sidebar a:hover {
		background: #dcfce7;
		color: #059669;
		transform: translateX(5px);
	}

	.sidebar a.active {
		background: linear-gradient(135deg, #4ade80, #22d3ee);
		color: white;
		font-weight: 600;
	}

	.sidebar a.active:hover {
		transform: none;
	}

	.icon {
		font-size: 1.2rem;
		width: 1.5rem;
		text-align: center;
	}

	.label {
		font-size: 0.95rem;
	}

	@media (max-width: 768px) {
		.sidebar {
			width: 100%;
			height: auto;
			border-right: none;
			border-bottom: 1px solid #bbf7d0;
		}
		
		.sidebar ul {
			display: flex;
			justify-content: center;
			overflow-x: auto;
		}
		
		.sidebar li {
			margin: 0 0.25rem;
			flex-shrink: 0;
		}

		.sidebar a {
			margin-right: 0;
			border-radius: 25px;
			padding: 0.5rem 1rem;
			white-space: nowrap;
		}

		.sidebar a:hover {
			transform: translateY(-2px);
		}
	}
</style>