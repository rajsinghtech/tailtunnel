<script lang="ts">
	import { onMount } from 'svelte';
	import CanaryPeerCard from '$lib/components/CanaryPeerCard.svelte';
	import LatencyChart from '$lib/components/LatencyChart.svelte';
	import type { PeerInfo, PingResult, PeersResponse, PingAllResponse } from '$lib/types/canary';

	let peers = $state<PeerInfo[]>([]);
	let pingResults = $state<Map<string, PingResult>>(new Map());
	let pingHistory = $state<Map<string, Array<{ timestamp: Date; latency: number; connectionType: string }>>>(new Map());
	let loading = $state(false);
	let pinging = $state(false);
	let error = $state('');
	let lastUpdate = $state<Date | null>(null);
	let autoRefresh = $state(true);
	let refreshInterval: number;
	let pollInterval: number;
	let searchQuery = $state('');
	let showChart = $state(true);
	const MAX_HISTORY_POINTS = 50;

	async function loadPeers() {
		loading = true;
		error = '';
		try {
			const response = await fetch('/api/canary/peers');
			if (!response.ok) {
				throw new Error(`Failed to load peers: ${response.statusText}`);
			}
			const data: PeersResponse = await response.json();
			peers = data.peers;
			lastUpdate = new Date(data.timestamp);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load peers';
		} finally {
			loading = false;
		}
	}

	async function pingAll() {
		if (pinging) return;

		pinging = true;
		error = '';
		try {
			const response = await fetch('/api/canary/ping-all', {
				method: 'POST'
			});
			if (!response.ok) {
				throw new Error(`Failed to ping: ${response.statusText}`);
			}
			const data: PingAllResponse = await response.json();

			const newResults = new Map<string, PingResult>();
			const timestamp = new Date(data.timestamp);

			for (const result of data.results) {
				newResults.set(result.ip, result);

				// Update history for successful pings
				if (result.success) {
					const peer = peers.find(p => p.ip === result.ip);
					if (peer) {
						const hostKey = peer.hostName;
						const history = pingHistory.get(hostKey) || [];

						history.push({
							timestamp,
							latency: result.latencyMs,
							connectionType: result.connectionType
						});

						// Keep only last MAX_HISTORY_POINTS
						if (history.length > MAX_HISTORY_POINTS) {
							history.shift();
						}

						pingHistory.set(hostKey, history);
					}
				}
			}

			pingResults = newResults;
			pingHistory = new Map(pingHistory); // Trigger reactivity
			lastUpdate = timestamp;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to ping peers';
		} finally {
			pinging = false;
		}
	}

	async function continuousPing() {
		if (!autoRefresh) return;
		await pingAll();
	}

	async function refresh() {
		await loadPeers();
		await pingAll();
	}

	function toggleAutoRefresh() {
		autoRefresh = !autoRefresh;
		if (autoRefresh) {
			startPolling();
		} else {
			stopPolling();
		}
	}

	function startPolling() {
		stopPolling();
		continuousPing();
		pollInterval = window.setInterval(continuousPing, 10000);
		refreshInterval = window.setInterval(loadPeers, 60000);
	}

	function stopPolling() {
		if (pollInterval) {
			clearInterval(pollInterval);
		}
		if (refreshInterval) {
			clearInterval(refreshInterval);
		}
	}

	const filteredPeers = $derived(() => {
		if (!searchQuery.trim()) return peers;

		const query = searchQuery.toLowerCase().trim();
		return peers.filter(peer => {
			const pingResult = pingResults.get(peer.ip);
			const connectionType = pingResult?.connectionType || '';
			const onlineStatus = peer.online ? 'online' : 'offline';

			return peer.hostName.toLowerCase().includes(query) ||
				peer.dnsName.toLowerCase().includes(query) ||
				peer.ip.toLowerCase().includes(query) ||
				peer.os.toLowerCase().includes(query) ||
				peer.userLogin?.toLowerCase().includes(query) ||
				peer.userDisplay?.toLowerCase().includes(query) ||
				peer.tags?.some(tag => tag.toLowerCase().includes(query)) ||
				connectionType.toLowerCase().includes(query) ||
				onlineStatus.includes(query);
		});
	});

	const filteredPingHistory = $derived(() => {
		if (!searchQuery.trim()) return pingHistory;

		const query = searchQuery.toLowerCase().trim();
		const filteredMap = new Map<string, Array<{ timestamp: Date; latency: number; connectionType: string }>>();

		for (const peer of filteredPeers()) {
			const history = pingHistory.get(peer.hostName);
			if (history && history.length > 0) {
				filteredMap.set(peer.hostName, history);
			}
		}

		return filteredMap;
	});

	function getTimeSince(date: Date | null): string {
		if (!date) return 'Never';
		const seconds = Math.floor((Date.now() - date.getTime()) / 1000);
		if (seconds < 60) return `${seconds}s ago`;
		const minutes = Math.floor(seconds / 60);
		if (minutes < 60) return `${minutes}m ago`;
		const hours = Math.floor(minutes / 60);
		return `${hours}h ago`;
	}

	onMount(() => {
		loadPeers().then(() => {
			if (autoRefresh) {
				startPolling();
			}
		});

		return () => {
			stopPolling();
		};
	});
</script>

<svelte:head>
	<title>TailCanary - Network Diagnostics - TailTunnel</title>
</svelte:head>

<div class="min-h-screen bg-background">
	<div class="border-b bg-card">
		<div class="container mx-auto p-4">
			<div class="flex items-center justify-between mb-4">
				<div>
					<a href="/" class="text-sm text-muted-foreground hover:text-foreground">
						← Back to machines
					</a>
					<h1 class="text-2xl font-bold mt-1">TailCanary</h1>
					<p class="text-sm text-muted-foreground">Network Diagnostics</p>
				</div>
				<div class="flex items-center gap-2">
					{#if lastUpdate}
						<span class="text-sm text-muted-foreground">
							Last updated: {getTimeSince(lastUpdate)}
						</span>
					{/if}
				</div>
			</div>

			<div class="flex items-center gap-2">
				<button
					onclick={refresh}
					disabled={loading || pinging}
					class="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed"
				>
					{#if loading || pinging}
						Refreshing...
					{:else}
						Refresh Now
					{/if}
				</button>

				<button
					onclick={toggleAutoRefresh}
					class={autoRefresh
						? "rounded-md bg-green-600 dark:bg-green-700 px-4 py-2 text-sm font-medium text-white hover:bg-green-700 dark:hover:bg-green-800"
						: "rounded-md border border-input bg-background px-4 py-2 text-sm font-medium hover:bg-accent"
					}
				>
					{#if autoRefresh}
						● Live (10s)
					{:else}
						○ Paused
					{/if}
				</button>

				{#if pinging}
					<span class="text-sm text-muted-foreground animate-pulse">
						Pinging {peers.length} peers...
					</span>
				{/if}
			</div>
		</div>
	</div>

	<div class="container mx-auto p-4">
		{#if error}
			<div class="rounded-lg bg-red-50 dark:bg-red-950 p-4 text-red-800 dark:text-red-200 mb-4">
				{error}
			</div>
		{/if}

		{#if !error && peers.length > 0 && filteredPingHistory().size > 0}
			<div class="mb-6">
				<LatencyChart pingHistory={filteredPingHistory()} />
			</div>
		{/if}

		{#if !error && peers.length > 0}
			<div class="mb-6">
				<div class="relative">
					<input
						type="text"
						bind:value={searchQuery}
						placeholder="Search by name, IP, user, OS, tag, status (online, offline), or connection type (direct, derp, peer-relay)..."
						class="w-full rounded-md border border-input bg-background px-4 py-2.5 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
					/>
					{#if searchQuery}
						<button
							onclick={() => searchQuery = ''}
							aria-label="Clear search"
							class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
						>
							<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
								<line x1="18" y1="6" x2="6" y2="18"></line>
								<line x1="6" y1="6" x2="18" y2="18"></line>
							</svg>
						</button>
					{/if}
				</div>
				{#if searchQuery && filteredPeers().length > 0}
					<p class="mt-2 text-sm text-muted-foreground">
						Found {filteredPeers().length} peer{filteredPeers().length === 1 ? '' : 's'}
					</p>
				{/if}
			</div>
		{/if}

		{#if loading}
			<div class="text-center py-8">
				<p class="text-muted-foreground">Loading peers...</p>
			</div>
		{:else if peers.length === 0}
			<div class="text-center py-8">
				<p class="text-muted-foreground">No peers found</p>
			</div>
		{:else if filteredPeers().length === 0}
			<div class="text-center py-8">
				<p class="text-muted-foreground">No peers match your search</p>
			</div>
		{:else}
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
				{#each filteredPeers() as peer (peer.ip)}
					<CanaryPeerCard {peer} pingResult={pingResults.get(peer.ip)} />
				{/each}
			</div>

			<div class="mt-6 text-sm text-muted-foreground">
				<p class="font-semibold mb-2">Connection Types:</p>
				<ul class="space-y-1 mb-4">
					<li><span class="text-green-600 dark:text-green-400 font-semibold">● Direct</span> - Direct peer-to-peer UDP connection (fastest, &lt;10ms typical)</li>
					<li><span class="text-yellow-600 dark:text-yellow-400 font-semibold">⚡ DERP Relay</span> - Relayed through Tailscale server (may upgrade to direct in ~5-30s)</li>
					<li><span class="text-blue-600 dark:text-blue-400 font-semibold">🔄 Peer Relay</span> - Relayed through another peer node</li>
				</ul>
				<p class="text-xs italic">
					Note: Connections initially use DERP relay during discovery, then automatically upgrade to direct paths when possible. This is normal behavior.
				</p>
			</div>
		{/if}
	</div>
</div>
