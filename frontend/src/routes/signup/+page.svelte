<script lang="ts">
	let phoneNumber = '';
	let password = '';

	async function signup() {
		const res = await fetch('/api/signup', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ phoneNumber, password })
		});
		console.log(res.status);
	}
</script>

<input type="text" placeholder="Phone Number" bind:value={phoneNumber} />
<input type="password" placeholder="Password" bind:value={password} />
<button on:click={signup}>signup</button>
