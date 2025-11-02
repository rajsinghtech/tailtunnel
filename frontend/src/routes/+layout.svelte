<script lang="ts">
	import '../app.css';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';

	let sidebarOpen = $state(false);

	const navItems = [
		{ href: '/', label: 'TailCanary', icon: '🐦' },
		{ href: '/machines', label: 'SSH Machines', icon: '🖥️' }
	];

	interface DiagnosticsInfo {
		hostName: string;
		tailscaleIP: string;
		dnsName: string;
		os: string;
		online: boolean;
		exitNodeId?: string;
		magicDnsSuffix: string;
		natType?: string;
		portMapProtocol?: string;
	}

	function getNATColor(natType: string | undefined): string {
		if (!natType) return '';
		if (natType === 'No NAT') return 'text-green-600 dark:text-green-400';
		if (natType === 'EZ NAT') return 'text-yellow-600 dark:text-yellow-400';
		if (natType === 'Hard NAT') return 'text-orange-600 dark:text-orange-400';
		return '';
	}

	let diagnostics = $state<DiagnosticsInfo | null>(null);
	let diagnosticsError = $state('');

	async function loadDiagnostics() {
		try {
			const response = await fetch('/api/diagnostics');
			if (!response.ok) {
				throw new Error('Failed to load diagnostics');
			}
			diagnostics = await response.json();
		} catch (e) {
			diagnosticsError = e instanceof Error ? e.message : 'Failed to load diagnostics';
		}
	}

	onMount(() => {
		loadDiagnostics();
		// Refresh diagnostics every 60 seconds
		const interval = setInterval(loadDiagnostics, 60000);
		return () => clearInterval(interval);
	});
</script>

<div class="min-h-screen bg-background flex">
	<!-- Sidebar -->
	<aside class="hidden md:flex md:flex-col md:w-64 md:h-screen border-r bg-card md:sticky md:top-0">
		<div class="p-6 border-b flex-shrink-0">
			<div class="flex items-center gap-3">
				<img src="/logo.svg" alt="TailTunnel" class="w-10 h-10" />
				<div>
					<h1 class="text-xl font-bold">TailTunnel</h1>
					<p class="text-xs text-muted-foreground">Tailscale Toolkit</p>
				</div>
			</div>
		</div>

		<nav class="flex-1 p-4 overflow-y-auto">
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

		<!-- Diagnostics Section -->
		<div class="border-t p-4 flex-shrink-0">
			<h3 class="text-xs font-semibold text-muted-foreground mb-3">DEVICE INFO</h3>
			{#if diagnosticsError}
				<p class="text-xs text-red-500">{diagnosticsError}</p>
			{:else if diagnostics}
				<div class="space-y-2 text-xs">
					<div>
						<span class="text-muted-foreground">Host:</span>
						<span class="ml-1 font-medium">{diagnostics.hostName}</span>
					</div>
					<div>
						<span class="text-muted-foreground">IP:</span>
						<span class="ml-1 font-mono text-xs">{diagnostics.tailscaleIP}</span>
					</div>
					<div>
						<span class="text-muted-foreground">OS:</span>
						<span class="ml-1">{diagnostics.os}</span>
					</div>
					{#if diagnostics.natType}
						<div>
							<span class="text-muted-foreground">NAT:</span>
							<span class="ml-1 font-semibold {getNATColor(diagnostics.natType)}">{diagnostics.natType}</span>
						</div>
					{/if}
					<div>
						<span class="text-muted-foreground">Status:</span>
						<span class="ml-1 {diagnostics.online ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'}">{diagnostics.online ? 'Online' : 'Offline'}</span>
					</div>
				</div>
			{:else}
				<p class="text-xs text-muted-foreground">Loading...</p>
			{/if}
		</div>
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
		<aside class="md:hidden fixed left-0 top-[57px] bottom-0 z-50 w-64 bg-card border-r flex flex-col">
			<nav class="flex-1 p-4">
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

			<!-- Diagnostics Section (Mobile) -->
			<div class="border-t p-4">
				<h3 class="text-xs font-semibold text-muted-foreground mb-3">DEVICE INFO</h3>
				{#if diagnosticsError}
					<p class="text-xs text-red-500">{diagnosticsError}</p>
				{:else if diagnostics}
					<div class="space-y-2 text-xs">
						<div>
							<span class="text-muted-foreground">Host:</span>
							<span class="ml-1 font-medium">{diagnostics.hostName}</span>
						</div>
						<div>
							<span class="text-muted-foreground">IP:</span>
							<span class="ml-1 font-mono text-xs">{diagnostics.tailscaleIP}</span>
						</div>
						<div>
							<span class="text-muted-foreground">OS:</span>
							<span class="ml-1">{diagnostics.os}</span>
						</div>
						{#if diagnostics.natType}
							<div>
								<span class="text-muted-foreground">NAT:</span>
								<span class="ml-1 font-semibold {getNATColor(diagnostics.natType)}">{diagnostics.natType}</span>
							</div>
						{/if}
						<div>
							<span class="text-muted-foreground">Status:</span>
							<span class="ml-1 {diagnostics.online ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'}">{diagnostics.online ? 'Online' : 'Offline'}</span>
						</div>
					</div>
				{:else}
					<p class="text-xs text-muted-foreground">Loading...</p>
				{/if}
			</div>
		</aside>
	{/if}

	<!-- Main content -->
	<main class="flex-1 md:pt-0 pt-[57px]">
		<slot />
	</main>
</div>
