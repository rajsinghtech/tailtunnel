<script lang="ts">
	import type { PeerInfo, PingResult } from '$lib/types/canary';
	import { cn } from '$lib/utils/style';

	let { peer, pingResult }: { peer: PeerInfo; pingResult?: PingResult } = $props();

	function getConnectionIcon(type?: string): string {
		switch (type) {
			case 'direct':
				return '●';
			case 'derp':
				return '⚡';
			case 'peer-relay':
				return '🔄';
			case 'offline':
				return '○';
			default:
				return '?';
		}
	}

	function getConnectionColor(type?: string): string {
		switch (type) {
			case 'direct':
				return 'text-green-600 dark:text-green-400';
			case 'derp':
				return 'text-yellow-600 dark:text-yellow-400';
			case 'peer-relay':
				return 'text-blue-600 dark:text-blue-400';
			case 'offline':
				return 'text-gray-400';
			default:
				return 'text-gray-600';
		}
	}

	function getConnectionLabel(type?: string): string {
		switch (type) {
			case 'direct':
				return 'Direct';
			case 'derp':
				return 'DERP Relay';
			case 'peer-relay':
				return 'Peer Relay';
			case 'offline':
				return 'Offline';
			default:
				return 'Unknown';
		}
	}

	function getConnectionDescription(type?: string): string {
		switch (type) {
			case 'direct':
				return 'Peer-to-peer UDP connection';
			case 'derp':
				return 'May upgrade to direct (~5-30s)';
			case 'peer-relay':
				return 'Relayed through another peer';
			case 'offline':
				return 'Not reachable';
			default:
				return '';
		}
	}

	function formatPingError(error?: string): string {
		if (!error) return 'Unknown error';
		if (error.includes('context deadline exceeded') || error.includes('timeout')) {
			return 'Ping timeout (peer may be slow to respond)';
		}
		if (error.includes('no such host')) {
			return 'Host not found';
		}
		if (error.includes('connection refused')) {
			return 'Connection refused';
		}
		// Keep original error but make it more user-friendly
		return error.replace('error Post "http://local-tailscaled.sock/localapi/v0/ping', 'Ping failed')
			.replace('": context deadline exceeded', ' (timeout)')
			.substring(0, 100); // Truncate long errors
	}

</script>

<div
	class={cn(
		'rounded-lg border bg-card shadow-sm transition-all hover:shadow-md hover:border-primary/50',
		!peer.online && 'opacity-60'
	)}
>
	<div class="p-4">
		<div class="flex items-start justify-between mb-3">
			<div class="flex-1 min-w-0">
				<h3 class="text-base font-semibold truncate mb-1">{peer.hostName}</h3>
				<p class="text-xs text-muted-foreground truncate">{peer.dnsName}</p>
				<p class="text-xs text-muted-foreground truncate">{peer.ip}</p>
			</div>
			<span
				class={cn(
					'ml-2 inline-flex items-center rounded-full px-2 py-0.5 text-xs font-semibold whitespace-nowrap',
					peer.online
						? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-100'
						: 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-100'
				)}
			>
				{peer.online ? 'Online' : 'Offline'}
			</span>
		</div>

		<div class="space-y-2">
			<div class="flex items-center gap-2 text-sm">
				<span class="inline-flex items-center rounded-md border px-2 py-0.5 text-xs font-medium">
					{peer.os}
				</span>
				{#if peer.active}
					<span
						class="inline-flex items-center rounded-md bg-blue-50 dark:bg-blue-950 px-2 py-0.5 text-xs font-medium text-blue-700 dark:text-blue-300"
					>
						Active
					</span>
				{/if}
			</div>

			{#if peer.userDisplay || peer.userLogin}
				<div class="flex items-center gap-1 text-xs text-muted-foreground">
					<svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="flex-shrink-0">
						<path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"></path>
						<circle cx="12" cy="7" r="4"></circle>
					</svg>
					<span class="truncate">{peer.userDisplay || peer.userLogin}</span>
				</div>
			{/if}

			{#if peer.tags && peer.tags.length > 0}
				<div class="flex flex-wrap gap-1">
					{#each peer.tags.slice(0, 3) as tag}
						<span class="inline-flex items-center rounded-md bg-blue-50 dark:bg-blue-950 px-1.5 py-0.5 text-xs font-medium text-blue-700 dark:text-blue-300 border border-blue-200 dark:border-blue-800">
							{tag.replace('tag:', '')}
						</span>
					{/each}
					{#if peer.tags.length > 3}
						<span class="inline-flex items-center rounded-md bg-muted px-1.5 py-0.5 text-xs font-medium text-muted-foreground">
							+{peer.tags.length - 3}
						</span>
					{/if}
				</div>
			{/if}

			{#if pingResult && pingResult.success}
				<div class="border-t pt-2">
					<div class="flex items-center justify-between mb-1">
						<div class="flex items-center gap-2">
							<span class={cn('text-lg', getConnectionColor(pingResult.connectionType))}>
								{getConnectionIcon(pingResult.connectionType)}
							</span>
							<span class="text-sm font-medium">
								{getConnectionLabel(pingResult.connectionType)}
							</span>
							{#if pingResult.derpRegion}
								<span class="text-xs text-muted-foreground">({pingResult.derpRegion})</span>
							{/if}
						</div>
						<span class="text-sm font-semibold text-primary">
							{pingResult.latencyMs.toFixed(1)}ms
						</span>
					</div>
					<p class="text-xs text-muted-foreground">
						{getConnectionDescription(pingResult.connectionType)}
					</p>
					{#if pingResult.endpoint}
						<p class="text-xs text-muted-foreground mt-1 font-mono">{pingResult.endpoint}</p>
					{:else if pingResult.peerRelay}
						<p class="text-xs text-muted-foreground mt-1 font-mono">{pingResult.peerRelay}</p>
					{/if}
				</div>
			{:else if pingResult && !pingResult.success}
				<div class="border-t pt-2">
					<div class="text-sm text-yellow-600 dark:text-yellow-400">
						{formatPingError(pingResult.error)}
					</div>
				</div>
			{/if}
		</div>
	</div>
</div>
