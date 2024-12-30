<script lang="ts">
	import '../app.css';
	let { children } = $props();
</script>

<nav class="grid w-fit grid-cols-5 gap-4 p-4">
	<a class="underline" href="/">home</a>
	<a class="underline" href="/account">account</a>
	<a class="underline" href="/signin">signin</a>
	<a class="underline" href="/signup">signup</a>
	<button class="underline" onclick={() => fetch('/api/signout')}>signout</button>
</nav>
<div class="p-4">
	{@render children()}
</div>
