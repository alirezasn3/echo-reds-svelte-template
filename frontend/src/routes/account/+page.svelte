<script lang="ts">
	import { onMount } from 'svelte';

	let user: any;
	let status: number;

	onMount(async () => {
		const res = await fetch('/api/account');
		status = res.status;
		if (res.status === 200) {
			user = await res.json();
		} else {
			console.log(res.status);
		}
	});
</script>

<div>{status}</div>
<div>{user?.id}</div>
<div>{user?.phoneNumber}</div>
