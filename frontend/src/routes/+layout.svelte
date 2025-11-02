<script lang="ts">
	import '../app.css';
	import { page } from '$app/stores';

	let sidebarOpen = $state(false);

	const navItems = [
		{ href: '/', label: 'TailCanary', icon: '🐦' },
		{ href: '/machines', label: 'SSH Machines', icon: '🖥️' }
	];
</script>

<div class="min-h-screen bg-background flex">
	<!-- Sidebar -->
	<aside class="hidden md:flex md:flex-col md:w-64 border-r bg-card">
		<div class="p-6 border-b">
			<div class="flex items-center gap-3">
				<img src="/logo.svg" alt="TailTunnel" class="w-10 h-10" />
				<div>
					<h1 class="text-xl font-bold">TailTunnel</h1>
					<p class="text-xs text-muted-foreground">Tailscale Toolkit</p>
				</div>
			</div>
		</div>

		<nav class="flex-1 p-4">
			<ul class="space-y-2">
				{#each navItems as item}
					<li>
						<a
							href={item.href}
							class="flex items-center gap-3 px-4 py-3 rounded-md transition-colors {$page.url.pathname === item.href ? 'bg-primary text-primary-foreground' : 'hover:bg-accent'}"
						>
							<span class="text-xl">{item.icon}</span>
							<span class="font-medium">{item.label}</span>
						</a>
					</li>
				{/each}
			</ul>
		</nav>
	</aside>

	<!-- Mobile header -->
	<div class="md:hidden fixed top-0 left-0 right-0 z-50 bg-card border-b">
		<div class="flex items-center justify-between p-4">
			<div class="flex items-center gap-2">
				<img src="/logo.svg" alt="TailTunnel" class="w-8 h-8" />
				<h1 class="text-lg font-bold">TailTunnel</h1>
			</div>
			<button
				onclick={() => sidebarOpen = !sidebarOpen}
				class="p-2 hover:bg-accent rounded-md"
				aria-label="Toggle menu"
			>
				<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
					<line x1="3" y1="12" x2="21" y2="12"></line>
					<line x1="3" y1="6" x2="21" y2="6"></line>
					<line x1="3" y1="18" x2="21" y2="18"></line>
				</svg>
			</button>
		</div>
	</div>

	<!-- Mobile sidebar -->
	{#if sidebarOpen}
		<div class="md:hidden fixed inset-0 z-40 bg-background/80 backdrop-blur-sm" onclick={() => sidebarOpen = false}></div>
		<aside class="md:hidden fixed left-0 top-[57px] bottom-0 z-50 w-64 bg-card border-r">
			<nav class="p-4">
				<ul class="space-y-2">
					{#each navItems as item}
						<li>
							<a
								href={item.href}
								onclick={() => sidebarOpen = false}
								class="flex items-center gap-3 px-4 py-3 rounded-md transition-colors {$page.url.pathname === item.href ? 'bg-primary text-primary-foreground' : 'hover:bg-accent'}"
							>
								<span class="text-xl">{item.icon}</span>
								<span class="font-medium">{item.label}</span>
							</a>
						</li>
					{/each}
				</ul>
			</nav>
		</aside>
	{/if}

	<!-- Main content -->
	<main class="flex-1 md:pt-0 pt-[57px]">
		<slot />
	</main>
</div>
